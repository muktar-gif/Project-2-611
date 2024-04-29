package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
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

type fileTransferServer struct {
	pb.UnimplementedJobServiceServer
}

// Job Description
type job struct {
	datafile string
	start    int
	length   int
}

// Result Description
type result struct {
	jobFound    job
	numOfPrimes int
}

var (

	// Channels to store jobs
	jobQueue = make(chan job)
	path     = ""
	globalN  = 0
)

// Function to check for errors in file operations
func checkFileOper(e error) {
	if e != nil {
		panic(e)
	}
}

// Global queue to store number of completed jobs for each worker (for stats)
var completedJobs = make(chan int)

func (s *dispatcherServer) RequestJob(ctx context.Context, empty *emptypb.Empty) (*pb.Job, error) {

	//return &pb.Job{Datafile: "HELLO", Start: int32(2), Length: int32(4)}, nil

	//fmt.Println("ruuning from server")

	select {
	case job := <-jobQueue:
		protoJob := &pb.Job{Datafile: job.datafile, Start: int32(job.start), Length: int32(job.length)}
		return protoJob, nil
	default:
		//job queue is empty
		return nil, nil
	}

}

func fillJobQueue() {

	// Opens file and checks for error
	f, err := os.Open(path)
	checkFileOper(err)

	// Prepares job byte to store data from file with N bytes
	jobByte := make([]byte, globalN)
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

			fmt.Println("Added job")
			// Creates job and adds it to job queue
			makeJob := job{path, start, readJob}
			jobQueue <- makeJob
			start += readJob

		}
	}

}

// Function to read file and creates N or less sized jobs for a job queue
func dispatcher(pathname *string, N *int, wg *sync.WaitGroup) {

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterJobServiceServer(grpcServer, &dispatcherServer{})
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Function to count the total amount of primes in the result queue
func consolidator() {

	// 	totalPrime := 0

	// 	// Loops through resule queue and totals the number of primes
	// 	for result := range resultQueue {
	// 		totalPrime += result.numOfPrimes
	// 	}

	// 	returnTotal <- totalPrime
	// 	done <- true
	// }

	// func oldMain() {

	// 	trackStart := time.Now()

	// 	// Default command line arguments
	// 	pathname := flag.String("pathname", "", "File path to binary file")
	// 	M := flag.Int("M", 1, "Number of worker threads")
	// 	N := flag.Int("N", 64*1024, "Number of bytes to segment input file")
	// 	//C := flag.Int("C", 1024, "Number of bytes to read from data file")
	// 	flag.Parse()

	// 	// Channels to store jobs and results
	// 	// jobQueue := make(chan job)
	// 	resultQueue := make(chan result)
	// 	getTotalPrime := make(chan int)

	// 	// Used for channel sync between jobQueue and resultQueue
	// 	var wg sync.WaitGroup
	// 	done := make(chan bool, 1)

	// 	// Stores worker's number of completed jobs
	// 	completedJobsArray := make([]int, *M)
	// 	totalSum := 0

	// 	// Calling all routines
	// 	go dispatcher(pathname, N, jobQueue, &wg)

	// 	// for i := 0; i < *M; i++ {
	// 	// 	wg.Add(1)
	// 	// 	go worker(C, jobQueue, resultQueue, &wg)
	// 	// }

	// 	go consolidator(getTotalPrime, resultQueue, done)

	// 	// Takes all number of completed jobs per worker for statistics
	// 	for i := 0; i < *M; i++ {
	// 		completedJobsArray[i] = <-completedJobs
	// 		totalSum += completedJobsArray[i]
	// 	}

	// 	// Waits until all workers are done to close resultQueue
	// 	wg.Wait()
	// 	close(resultQueue)

	// 	// Channel to wait for consolidator to finish (Both not needed, either can be used to sync)
	// 	prime := <-getTotalPrime
	// 	<-done

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

}

func main() {

	pathname := flag.String("pathname", "", "File path to binary file")
	N := flag.Int("N", 64*1024, "Number of bytes to segment input file")
	//C := flag.Int("C", 1024, "Number of bytes to read from data file")
	flag.Parse()
	fmt.Println("run go routine")

	var wg sync.WaitGroup
	wg.Add(1)

	path = *pathname
	globalN = *N
	go dispatcher(pathname, N, &wg)

	for job := range jobQueue {
		// Process the job here
		fmt.Println("Processing job:", job)
	}

	wg.Wait()
	// Closes queue once there are no more jobs to insert
	close(jobQueue)

	//go consolidator()

}
