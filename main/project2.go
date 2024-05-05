package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

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

type serverAddress struct {
	port      int
	ipAddress string
}

var (
	// Channels to store jobs
	jobQueue = make(chan job)

	// Channels to store results
	resultQueue = make(chan result)

	// Channel to communites total jobs to consolidator
	totalJobs = make(chan int, 1)

	// Channels to signal to dispatcher the closure of the consolidator
	doneDispatcher = make(chan bool)

	// Global array to store number of completed jobs for each worker (for stats)
	completedJobsArray = []int{}

	// List of server addresses and ports
	serverList = make(map[string]serverAddress)

	// Store ports from server list
	dispatcherPort   string
	consolidatorPort string
)

// Function to check for errors in file operations
func checkFileOper(e error) {
	if e != nil {
		panic(e)
	}
}

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

		// Tracks workers completed jobs for stats
		completedJobsArray = append(completedJobsArray, int(info.RequestConfirmation.NumOfJobCompleted))

		// Check for all workers connection
		if s.countConnection == 0 {
			close(resultQueue)
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

	// Start listening for dispatcher server
	lis, err := net.Listen("tcp", ":"+dispatcherPort)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	wg.Add(1)
	// Routine to start adding jobs to queue
	go func() {

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
	lis, err := net.Listen("tcp", ":"+consolidatorPort)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	// Retrieves the number of jobs the dispatcher has created
	getExpectedJobs := <-totalJobs

	// Buffers the expecting numbers of results
	resultQueue = make(chan result, getExpectedJobs)

	pb.RegisterJobServiceServer(grpcServer, &consolidatorServer{expectedJobs: getExpectedJobs, jobsReceived: 0})

	wg.Add(1)
	// Routine to start collecting jobs results to result queue
	go func() {

		defer wg.Done()

		totalPrime := 0

		// Loops through resule queue and totals the number of primes
		for result := range resultQueue {
			totalPrime += result.numOfPrimes
		}

		returnTotal <- totalPrime

		// Signal to dispatcher to close
		doneDispatcher <- true

		fmt.Println("Stopping consolidator server...")
		grpcServer.GracefulStop()

	}()

	fmt.Println("Starting consolidator server...")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {

	trackStart := time.Now()

	configFile, err := os.Open("primes_config.txt")

	if err != nil {
		panic(err)
	}

	readFile := bufio.NewScanner(configFile)

	// Read line
	for readFile.Scan() {
		serverInfo := readFile.Text()
		serverData := strings.Split(serverInfo, " ")
		getPort, _ := strconv.Atoi(serverData[2])
		serverList[serverData[0]] = serverAddress{getPort, serverData[1]}
	}

	dispatcherPort = strconv.Itoa(serverList["dispatcher"].port)
	consolidatorPort = strconv.Itoa(serverList["consolidator"].port)

	pathname := flag.String("pathname", "", "File path to binary file")
	N := flag.Int("N", 64*1024, "Number of bytes to segment input file")
	C := flag.Int("C", 1024, "Number of bytes to read from data file")
	flag.Parse()

	getTotalPrime := make(chan int, 2)

	var wg sync.WaitGroup

	wg.Add(1)
	go dispatcher(pathname, N, C, &wg)

	wg.Add(1)
	go consolidator(getTotalPrime, &wg)

	wg.Wait()

	// Stores worker's number of completed jobs
	totalSum := 0

	for getTotal := range completedJobsArray {
		totalSum += getTotal
	}

	// Prints out min, max, average, and median for the list of jobs a worker finished and elapsed time of main
	slices.Sort(completedJobsArray)

	comLen := len(completedJobsArray)
	if comLen != 0 {
		fmt.Println("Min # job a worker completed:", completedJobsArray[0])
		fmt.Println("Max # job a worker completed:", completedJobsArray[len(completedJobsArray)-1])
		fmt.Println("Average # job a worker completed:", (float64(totalSum) / float64(len(completedJobsArray))))

		var medianJob float64
		if comLen%2 == 0 {
			medianJob = float64(completedJobsArray[comLen/2]+completedJobsArray[comLen/2-1]) / float64(2)
		} else {
			medianJob = float64(completedJobsArray[comLen/2])
		}

		fmt.Println("Median # job a worker completed:", medianJob)
	}

	fmt.Println("Total primes numbers are", <-getTotalPrime)

	trackEnd := time.Now()
	fmt.Println("Elapsed Time:", trackEnd.Sub(trackStart))
}
