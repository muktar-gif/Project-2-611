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
	expectedJobs int
	jobsReceived int
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

	// Channel to track total primes
	totalPrime = make(chan int)

	doneConsolidator = make(chan bool)
	doneDispatcher   = make(chan bool)
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

func (s *consolidatorServer) PushResult(ctx context.Context, pushedResults *pb.JobResult) (*pb.TerminateRequest, error) {

	// Signals to close server, workers will terminate
	if s.expectedJobs == s.jobsReceived {
		doneConsolidator <- true
		return &emptypb.Empty{}, nil
	} else {

		// Pushes result into the queue
		makeResult := result{job{pushedResults.JobFound.Datafile, int(pushedResults.JobFound.Start), int(pushedResults.JobFound.Length), int(pushedResults.JobFound.CValue)}, int(pushedResults.NumOfPrimes)}
		resultQueue <- makeResult
		s.jobsReceived++

		slog.Info(fmt.Sprintf("Consolidator-- Job: datafile: %s, start: %d, length: %d, # Primes: %d",
			pushedResults.JobFound.Datafile, pushedResults.JobFound.Start, pushedResults.JobFound.Length, pushedResults.NumOfPrimes))

	}

	// Signal to close result queue
	if s.expectedJobs == s.jobsReceived {
		close(resultQueue)
	}

	return &emptypb.Empty{}, nil
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
		grpcServer.GracefulStop()
		fmt.Println("Stopping dispatcher server...")
	}()

	pb.RegisterJobServiceServer(grpcServer, &dispatcherServer{})
	fmt.Println("Starting dispatcher server...")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Function to count the total amount of primes in the result queue
func consolidator(wg *sync.WaitGroup) {

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
		//<-doneConsolidator

		for range doneConsolidator {
		}

		grpcServer.GracefulStop()
		fmt.Println("Stopping consolidator server...")

		// Signal to dispatcher to close
		doneDispatcher <- true
	}()

	fmt.Println("Starting consolidator server...")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {

	pathname := flag.String("pathname", "", "File path to binary file")
	N := flag.Int("N", 64*1024, "Number of bytes to segment input file")
	C := flag.Int("C", 1024, "Number of bytes to read from data file")
	flag.Parse()

	var wg sync.WaitGroup

	wg.Add(1)
	go dispatcher(pathname, N, C, &wg)

	wg.Add(1)
	go consolidator(&wg)

	wg.Wait()

}
