package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
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
	consolidatorGrpcServer *grpc.Server
	dispatcherGrpcServer   *grpc.Server

	// Channels to store jobs
	jobQueue = make(chan job)

	// Channels to store results
	resultQueue = make(chan result)

	// Channel to communites total jobs to consolidator
	totalJobs = make(chan int, 1)

	// Channel to track total primes
	totalPrime = make(chan int)

	done = make(chan bool)
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

func (s *consolidatorServer) PushResult(ctx context.Context, pushedResults *pb.JobResult) (*emptypb.Empty, error) {

	// Signals to close server, workers will terminate
	if s.expectedJobs == s.jobsReceived {
		done <- true
		return &emptypb.Empty{}, nil
	} else {
		makeResult := result{job{pushedResults.JobFound.Datafile, int(pushedResults.JobFound.Start), int(pushedResults.JobFound.Length), int(pushedResults.JobFound.CValue)}, int(pushedResults.NumOfPrimes)}
		resultQueue <- makeResult
		s.jobsReceived++

		fmt.Println("Consolidator -- Printing pushed result from worker, Job: ", makeResult.jobFound, " # of Primes: ", makeResult.numOfPrimes)
		fmt.Println("Expected total jobs:", s.expectedJobs, "Jobs received:", s.jobsReceived)
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
		<-done
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
		<-done
		grpcServer.GracefulStop()
		fmt.Println("Stopping consolidator server...")
	}()

	fmt.Println("Starting consolidator server...")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	// 	totalPrime := 0

	// 	// Loops through resule queue and totals the number of primes
	// 	for result := range resultQueue {
	// 		totalPrime += result.numOfPrimes
	// 	}

	// 	returnTotal <- totalPrime
	// 	done <- true
	// }
}

// func oldMain() {

// 	trackStart := time.Now()

// 	// Prints out min, max, average, and median for the list of jobs a worker finished and elapsed time of main
// 	slices.Sort(completedJobsArray)

// 	comLen := len(completedJobsArray)
// 	if comLen != 0 {
// 		fmt.Println("Min # job a worker completed:", completedJobsArray[0])
// 		fmt.Println("Max # job a worker completed:", completedJobsArray[len(completedJobsArray)-1])
// 		fmt.Println("Average # job a worker completed:", (float64(totalSum) / float64(len(completedJobsArray))))

// 		var medianJob float64
// 		if comLen%2 == 0 {
// 			medianJob = float64(completedJobsArray[comLen/2]+completedJobsArray[comLen/2-1]) / float64(2)
// 		} else {
// 			medianJob = float64(completedJobsArray[comLen/2])
// 		}

// 		fmt.Println("Median # job a worker completed:", medianJob)
// 	}

// 	fmt.Println("Total primes numbers are", prime)

// 	trackEnd := time.Now()
// 	fmt.Println("Elapsed Time:", trackEnd.Sub(trackStart))

//}

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
