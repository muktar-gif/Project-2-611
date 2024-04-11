package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"math/rand/v2"
	"os"
	"slices"
	"sync"
	"time"
)

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

// Function to check for errors in file operations
func checkFileOper(e error) {
	if e != nil {
		panic(e)
	}
}

// Global queue to store number of completed jobs for each worker (for stats)
var completedJobs = make(chan int)

// Function to read file and creates N or less sized jobs for a job queue
func dispatcher(pathname *string, N *int, jobQueue chan<- job, wg *sync.WaitGroup) {

	// Opens file and checks for error
	f, err := os.Open(*pathname)
	checkFileOper(err)

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
			makeJob := job{*pathname, start, readJob}
			jobQueue <- makeJob
			start += readJob

		}
	}

	// Closes queue once there are no more jobs to insert
	close(jobQueue)
}

// Function too pull jobs from the job queue and count prime numbers in C or less bytes
// Counted results are inserted into a result queue
func worker(C *int, jobQueue <-chan job, resultQueue chan<- result, wg *sync.WaitGroup) {

	defer wg.Done()

	numOfJobsDone := 0

	// Loop through job queue channel
	for job := range jobQueue {

		// Sleep worker between 400 and 600 ms
		randTime := rand.IntN(600-400) + 400
		time.Sleep(time.Duration(randTime) * time.Millisecond)

		// Open job datafile
		f, err := os.Open(job.datafile)
		checkFileOper(err)

		// Move file pointer to job start
		f.Seek(int64(job.start), 0)

		numOfPrimes := 0
		totalJobLenBytes := 0

		// While the total read is less than the job's length
		for totalJobLenBytes < job.length {

			// Buffer for reading C bytes
			jobData := make([]byte, *C)

			readJob, err := f.Read(jobData)
			checkFileOper(err)

			// Tracks total bytes read
			totalJobLenBytes += readJob

			// Corrects if buffer reads more than the job length
			if totalJobLenBytes > job.length {
				readJob -= (totalJobLenBytes - job.length)
			}

			totalBytesRead := 0

			// While the number of 8 bytes left is less than the single job
			for totalBytesRead < readJob {

				var numBytes []byte

				// Inserts 0s if less than 8 bytes are left else take 8 bytes
				if (readJob - totalBytesRead) < 8 {

					bytesLeft := readJob - totalBytesRead
					zeroBytes := make([]byte, 8-bytesLeft)
					numBytes = append(jobData[totalBytesRead:totalBytesRead+bytesLeft], zeroBytes...)

				} else {

					numBytes = jobData[totalBytesRead : totalBytesRead+8]

				}

				// Converts unsigned 64bit in little endian order to decimal
				checkNum := binary.LittleEndian.Uint64(numBytes[:8])

				// Checks and adds the number of primes within the whole job
				if big.NewInt(int64(checkNum)).ProbablyPrime(0) {
					numOfPrimes++
				}

				// Increments total read bytes
				totalBytesRead += 8
			}

		}

		// Inserts job results into result channel
		makeResult := result{job, numOfPrimes}
		resultQueue <- makeResult

		slog.Info(fmt.Sprintf("Job: %#v Primes in Job: %d", job, numOfPrimes))

		numOfJobsDone++
	}

	completedJobs <- numOfJobsDone
}

// Function to count the total amount of primes in the result queue
func consolidator(returnTotal chan<- int, resultQueue <-chan result, done chan bool) {

	totalPrime := 0

	// Loops through resule queue and totals the number of primes
	for result := range resultQueue {
		totalPrime += result.numOfPrimes
	}

	returnTotal <- totalPrime
	done <- true
}

func main() {

	trackStart := time.Now()

	// Default command line arguments
	pathname := flag.String("pathname", "", "File path to binary file")
	M := flag.Int("M", 1, "Number of worker threads")
	N := flag.Int("N", 64*1024, "Number of bytes to segment input file")
	C := flag.Int("C", 1024, "Number of bytes to read from data file")
	flag.Parse()

	// Channels to store jobs and results
	jobQueue := make(chan job)
	resultQueue := make(chan result)
	getTotalPrime := make(chan int)

	// Used for channel sync between jobQueue and resultQueue
	var wg sync.WaitGroup
	done := make(chan bool, 1)

	// Stores worker's number of completed jobs
	completedJobsArray := make([]int, *M)
	totalSum := 0

	// Calling all routines
	go dispatcher(pathname, N, jobQueue, &wg)

	for i := 0; i < *M; i++ {
		wg.Add(1)
		go worker(C, jobQueue, resultQueue, &wg)
	}

	go consolidator(getTotalPrime, resultQueue, done)

	// Takes all number of completed jobs per worker for statistics
	for i := 0; i < *M; i++ {
		completedJobsArray[i] = <-completedJobs
		totalSum += completedJobsArray[i]
	}

	// Waits until all workers are done to close resultQueue
	wg.Wait()
	close(resultQueue)

	// Channel to wait for consolidator to finish (Both not needed, either can be used to sync)
	prime := <-getTotalPrime
	<-done

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

	fmt.Println("Total primes numbers are", prime)

	trackEnd := time.Now()
	fmt.Println("Elapsed Time:", trackEnd.Sub(trackStart))

}
