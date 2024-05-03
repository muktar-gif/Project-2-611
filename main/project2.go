package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/muktar-gif/Project-2-611/jobproto"
)

type dispatcherServer struct {
	pb.UnimplementedJobServiceServer
}

type consolidatorServer struct {
	pb.UnimplementedJobServiceServer
	expectedJobs    int
	jobsReceived    int
	countConnection int
	mu              sync.Mutex
}

// Job Description
type job struct {
	datafile string
	start    int
	length   int
	cValue   int
}

// Result Description
type result struct {
	jobFound    job
	numOfPrimes int
}

var (

	// Channels to store jobs
	jobQueue = make(chan job)

	// Channels to store results
	resultQueue = make(chan result)

	// Channel to communites total jobs to consolidator
	totalJobs = make(chan int, 1)

	doneConsolidator = make(chan bool)
	doneDispatcher   = make(chan bool)

	completedJobs = make(chan int)
)

// Function to check for errors in file operations
func checkFileOper(e error) {
	if e != nil {
		panic(e)
	}
}

// Global queue to store number of completed jobs for each worker (for stats)
// var completedJobs = make(chan int)

func (s *dispatcherServer) RequestJob(ctx context.Context, empty *emptypb.Empty) (*pb.Job, error) {

	select {
	// Job Queue is not empty pulls job from queue and returns the job
	case job := <-jobQueue:
		protoJob := &pb.Job{Datafile: job.datafile, Start: int32(job.start), Length: int32(job.length), CValue: int32(job.cValue)}
		return protoJob, nil
	// Job queue is empty
	default:
		return &pb.Job{Datafile: "", Start: 0, Length: 0, CValue: 0}, nil
	}

}

func (s *consolidatorServer) EstablishConnection(ctx context.Context, connection *pb.Connected) (*emptypb.Empty, error) {

	if connection.Connection {
		s.mu.Lock()
		s.countConnection++
		s.mu.Unlock()
	} else {
		s.mu.Lock()
		s.countConnection--
		s.mu.Unlock()
	}

	return &emptypb.Empty{}, nil
}

func (s *consolidatorServer) PushResult(ctx context.Context, info *pb.PushInfo) (*pb.TerminateRequest, error) {

	// Makes sure all workers have terminate so server can close
	if info.RequestConfirmation.Confirmed {

		completedJobs <- int(info.RequestConfirmation.NumOfJobCompleted)
		doneConsolidator <- true
		// Closes check for all workers connection
		if s.countConnection == 0 {
			close(resultQueue)
			close(doneConsolidator)
		}

		return nil, nil
		// Signals to close server, workers will terminate
	} else if s.expectedJobs == s.jobsReceived {

		return &pb.TerminateRequest{Request: true}, nil
		// Prevents null pushes
	} else if info.Result.JobFound.Datafile != "" {

		// Pushes result into the queue
		makeResult := result{job{info.Result.JobFound.Datafile, int(info.Result.JobFound.Start), int(info.Result.JobFound.Length), int(info.Result.JobFound.CValue)}, int(info.Result.NumOfPrimes)}
		resultQueue <- makeResult

		s.mu.Lock()
		s.jobsReceived++
		s.mu.Unlock()

		slog.Info(fmt.Sprintf("Consolidator-- Job: datafile: %s, start: %d, length: %d, # Primes: %d",
			info.Result.JobFound.Datafile, info.Result.JobFound.Start, info.Result.JobFound.Length, info.Result.NumOfPrimes))

	}

	// Signal to close result queue
	if s.expectedJobs == s.jobsReceived {
		return &pb.TerminateRequest{Request: true}, nil
	}

	return &pb.TerminateRequest{Request: false}, nil
}

// Function to read file and creates N or less sized jobs for a job queue
func dispatcher(pathname *string, N *int, C *int, wg *sync.WaitGroup) {

	defer wg.Done()

	// Opens file and checks for error
	f, err := os.Open(*pathname)
	checkFileOper(err)

	// Calculates the number of jobs and updates jobqueue buffer
	fileSize, err := f.Stat()
	jobTotal := int64(math.Ceil(float64(fileSize.Size()) / float64(*N)))
	jobQueue = make(chan job, jobTotal)

	totalJobs <- int(jobTotal)

	// Prepares job byte to store data from file with N bytes
	jobByte := make([]byte, *N)
	var readJob int
	var start int = 0

	for err != io.EOF {

		// Reads from data file, returns length read (of bytes)
		readJob, err = f.Read(jobByte)

		// Checks error, except end of file
		if err != nil && err != io.EOF {
			panic(err)
		}

		if err != io.EOF {

			// Creates job and adds it to job queue
			makeJob := job{*pathname, start, readJob, *C}
			jobQueue <- makeJob
			start += readJob

		}
	}

	// Closes queue once there are no more jobs to insert
	close(jobQueue)

	// Start listening for dispatcher server
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	go func() {

		// Waiting for signal to close
		<-doneDispatcher
		fmt.Println("Stopping dispatcher server...")
		grpcServer.GracefulStop()

	}()

	pb.RegisterJobServiceServer(grpcServer, &dispatcherServer{})
	fmt.Println("Starting dispatcher server...")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Function to count the total amount of primes in the result queue
func consolidator(returnTotal chan<- int, wg *sync.WaitGroup) {

	defer wg.Done()

	// Start listening for consolidator server
	lis, err := net.Listen("tcp", ":5002")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	// Retrieves the number of jobs the dispatcher has created
	getExpectedJobs := <-totalJobs

	// Buffers the expecting numbers of results
	resultQueue = make(chan result, getExpectedJobs)

	pb.RegisterJobServiceServer(grpcServer, &consolidatorServer{expectedJobs: getExpectedJobs, jobsReceived: 0})

	go func() {

		// Waiting for signal to close
		for range doneConsolidator {
		}
		fmt.Println("Stopping consolidator server...")
		grpcServer.GracefulStop()

		// Signal to dispatcher to close
		doneDispatcher <- true
	}()

	fmt.Println("Starting consolidator server...")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	totalPrime := 0

	// Loops through resule queue and totals the number of primes
	for result := range resultQueue {
		totalPrime += result.numOfPrimes
	}

	returnTotal <- totalPrime
}

func main() {

	pathname := flag.String("pathname", "", "File path to binary file")
	N := flag.Int("N", 64*1024, "Number of bytes to segment input file")
	C := flag.Int("C", 1024, "Number of bytes to read from data file")
	flag.Parse()

	getTotalPrime := make(chan int)

	var wg sync.WaitGroup

	wg.Add(1)
	go dispatcher(pathname, N, C, &wg)

	wg.Add(1)
	go consolidator(getTotalPrime, &wg)

	wg.Wait()

}
