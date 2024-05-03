package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math/big"
	"math/rand/v2"
	"os"
	"time"

	//"io"
	"log"

	filePb "github.com/muktar-gif/Project-2-611/tree/direction-version/fileproto"
	jobPb "github.com/muktar-gif/Project-2-611/tree/direction-version/jobproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Function too pull jobs from the job queue and count prime numbers in C or less bytes
// Counted results are inserted into a result queue

func main() {

	defaultC := 1024
	C := flag.Int("C", defaultC, "Number of bytes to read from data file")
	flag.Parse()

	// Create dispatcher client
	dispatcherConn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer dispatcherConn.Close()

	dispatcherClient := jobPb.NewJobServiceClient(dispatcherConn)

	// Create consolidator client
	consolidatorConn, err := grpc.Dial("localhost:5002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer consolidatorConn.Close()

	consolidatorClient := jobPb.NewJobServiceClient(consolidatorConn)

	// Create file server client
	fileConn, err := grpc.Dial("localhost:5003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer fileConn.Close()

	fileClient := filePb.NewFileServiceClient(fileConn)

	for {

		// Sleep worker between 400 and 600 ms
		randTime := rand.IntN(600-400) + 400
		time.Sleep(time.Duration(randTime) * time.Millisecond)

		// Call to dispatcher server to request job, worker will continuously request a job
		getJob, err := dispatcherClient.RequestJob(context.Background(), &emptypb.Empty{})
		numOfPrimes := 0

		// Terminates after failing to submitted to dispatcher
		if err != nil {
			fmt.Println("Terminating worker...")
			os.Exit(0)
		}

		if getJob.Datafile != "" {

			var fileSeg *filePb.FileSegmentRequest

			// Priortizes C value from worker, else it will use the one expected from the dispatcher
			if *C == defaultC {
				fileSeg = &filePb.FileSegmentRequest{Datafile: getJob.Datafile, Start: getJob.Start, Length: getJob.Length, CValue: getJob.CValue}
			} else {
				fileSeg = &filePb.FileSegmentRequest{Datafile: getJob.Datafile, Start: getJob.Start, Length: getJob.Length, CValue: int32(*C)}
			}

			// Call to file server to get data
			stream, err := fileClient.GetFileChunk(context.Background(), fileSeg)

			if err != nil {
				panic(err)
			}

			for {

				// Steam file data from job
				fileData, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatalf("fileClient.FileJob failed: %v", err)
				}

				// Converts unsigned 64bit in little endian order to decimal
				checkNum := binary.LittleEndian.Uint64(fileData.DataChunk[:8])

				// Checks and adds the number of primes within the whole job
				if big.NewInt(int64(checkNum)).ProbablyPrime(0) {
					numOfPrimes++
				}
			}
		}

		// Call to consolidator server to push job, worker will continuously push a job
		// Will receive termination from consolidator if no more jobs were found
		pushResults := &jobPb.JobResult{JobFound: getJob, NumOfPrimes: int32(numOfPrimes)}
		_, err = consolidatorClient.PushResult(context.Background(), pushResults)

		// Terminates after failing to submitted to consolidator
		if err != nil {
			fmt.Println("Terminating worker...")
			os.Exit(0)
		}

		if err != nil {
			log.Fatalf("consolidatorClient.RequestJob failed: %v", err)
		}
	}
}
