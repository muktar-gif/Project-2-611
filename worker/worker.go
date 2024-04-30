package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"math/rand/v2"
	"time"

	//"io"
	"log"
	"sync"

	filePb "github.com/muktar-gif/Project-2-611/fileproto"
	jobPb "github.com/muktar-gif/Project-2-611/jobproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
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

	fileConn, err := grpc.Dial("localhost:5003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer fileConn.Close()

	fileClient := filePb.NewFileServiceClient(fileConn)

	jobConn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer jobConn.Close()

	jobClient := jobPb.NewJobServiceClient(jobConn)

	for {

		// Sleep worker between 400 and 600 ms
		randTime := rand.IntN(600-400) + 400
		time.Sleep(time.Duration(randTime) * time.Millisecond)

		// Call to dispatcher server to request job
		getJob, err := jobClient.RequestJob(context.Background(), &emptypb.Empty{})

		if err != nil {
			log.Fatalf("client.RequestJob failed: %v", err)
		}

		fmt.Println(getJob)

		fileSeg := &filePb.FileSegmentRequest{Datafile: getJob.Datafile, Start: getJob.Start, Length: getJob.Length}

		// Call to file server to get data
		stream, err := fileClient.GetFileChunk(context.Background(), fileSeg)

		if err != nil {
			panic(err)
		}

		numOfPrimes := 0

		for {

			fileData, err := stream.Recv()

			fmt.Println("here")
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("client.FileJob failed: %v", err)
			}

			// Converts unsigned 64bit in little endian order to decimal
			checkNum := binary.LittleEndian.Uint64(fileData.DataChunk)
			// Checks and adds the number of primes within the whole job
			if big.NewInt(int64(checkNum)).ProbablyPrime(0) {
				numOfPrimes++
			}
		}

		fmt.Println(getJob, " Primes are ", numOfPrimes)
	}

}
