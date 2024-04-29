package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	filePb "github.com/muktar-gif/Project-2-611/fileproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Function to check for errors in file operations
func checkFileOper(e error) {
	if e != nil {
		panic(e)
	}
}

// Function too pull jobs from the job queue and count prime numbers in C or less bytes
// Counted results are inserted into a result queue
func worker(C *int, wg *sync.WaitGroup) {

	defer wg.Done()

	// numOfJobsDone := 0

	// Loop through job queue channel
	// for job := range jobQueue {

	// 	// Sleep worker between 400 and 600 ms
	// 	randTime := rand.IntN(600-400) + 400
	// 	time.Sleep(time.Duration(randTime) * time.Millisecond)

	// 	// Open job datafile
	// 	f, err := os.Open(job.datafile)
	// 	checkFileOper(err)

	// 	// Move file pointer to job start
	// 	f.Seek(int64(job.start), 0)

	// 	numOfPrimes := 0
	// 	totalJobLenBytes := 0

	// 	// While the total read is less than the job's length
	// 	for totalJobLenBytes < job.length {

	// 		// Buffer for reading C bytes
	// 		jobData := make([]byte, *C)

	// 		readJob, err := f.Read(jobData)
	// 		checkFileOper(err)

	// 		// Tracks total bytes read
	// 		totalJobLenBytes += readJob

	// 		// Corrects if buffer reads more than the job length
	// 		if totalJobLenBytes > job.length {
	// 			readJob -= (totalJobLenBytes - job.length)
	// 		}

	// 		totalBytesRead := 0

	// 		// While the number of 8 bytes left is less than the single job
	// 		for totalBytesRead < readJob {

	// 			var numBytes []byte

	// 			// Inserts 0s if less than 8 bytes are left else take 8 bytes
	// 			if (readJob - totalBytesRead) < 8 {

	// 				bytesLeft := readJob - totalBytesRead
	// 				zeroBytes := make([]byte, 8-bytesLeft)
	// 				numBytes = append(jobData[totalBytesRead:totalBytesRead+bytesLeft], zeroBytes...)

	// 			} else {

	// 				numBytes = jobData[totalBytesRead : totalBytesRead+8]

	// 			}

	// 			// Converts unsigned 64bit in little endian order to decimal
	// 			checkNum := binary.LittleEndian.Uint64(numBytes[:8])

	// 			// Checks and adds the number of primes within the whole job
	// 			if big.NewInt(int64(checkNum)).ProbablyPrime(0) {
	// 				numOfPrimes++
	// 			}

	// 			// Increments total read bytes
	// 			totalBytesRead += 8
	// 		}

	// 	}

	// 	// Inserts job results into result channel
	// 	makeResult := result{job, numOfPrimes}
	// 	resultQueue <- makeResult

	// 	slog.Info(fmt.Sprintf("Job: %#v Primes in Job: %d", job, numOfPrimes))

	// 	numOfJobsDone++
	// }

	// completedJobs <- numOfJobsDone
}

func main() {

	//var opts []grpc.DialOption
	conn, err := grpc.Dial("localhost:5003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := filePb.NewFileServiceClient(conn)

	fileSeg := &filePb.FileSegmentRequest{Datafile: "", Start: 2, Length: 2}

	stream, err := client.GetFileChunk(context.Background(), fileSeg)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Hello")
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		log.Println(feature)
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}

	}
}
