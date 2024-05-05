package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	pb "github.com/muktar-gif/Project-2-611/fileproto"
)

type fileTransferServer struct {
	pb.UnimplementedFileServiceServer
}

type serverAddress struct {
	port      int
	ipAddress string
}

// Function to check for errors in file operations
func checkFileOper(e error) {
	if e != nil {
		panic(e)
	}
}

func (s *fileTransferServer) GetFileChunk(fileRequest *pb.FileSegmentRequest, stream pb.FileService_GetFileChunkServer) error {

	// Open job datafile
	f, err := os.Open(fileRequest.Datafile)
	checkFileOper(err)

	// Move file pointer to job start
	f.Seek(int64(fileRequest.Start), 0)

	totalJobLenBytes := 0

	// While the total read is less than the job's length
	for totalJobLenBytes < int(fileRequest.Length) {

		// Buffer for reading C bytes
		jobData := make([]byte, fileRequest.CValue)

		readJob, err := f.Read(jobData)
		checkFileOper(err)

		// Tracks total bytes read
		totalJobLenBytes += readJob

		// Corrects if buffer reads more than the job length
		if totalJobLenBytes > int(fileRequest.Length) {
			readJob -= (totalJobLenBytes - int(fileRequest.Length))
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

			// Sends bytes to client
			if err := stream.Send(&pb.FileData{DataChunk: numBytes[:8]}); err != nil {
				return err
			}

			totalBytesRead += 8
		}

	}

	return nil
}

func main() {

	serverList := make(map[string]serverAddress)
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

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(serverList["fileserver"].port))

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, &fileTransferServer{})

	fmt.Println("Starting fileserver...")
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
