package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/muktar-gif/Project-2-611/fileproto"
)

type fileTransferServer struct {
	pb.UnimplementedFileServiceServer
}

// Function to check for errors in file operations
func checkFileOper(e error) {
	if e != nil {
		panic(e)
	}
}

func (s *fileTransferServer) GetFileChunk(fileRequest *pb.FileSegmentRequest, stream pb.FileService_GetFileChunkServer) error {

	log.Println(os.Getwd())
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

			// Converts unsigned 64bit in little endian order to decimal
			// checkNum := binary.LittleEndian.Uint64(numBytes[:8])

			if err := stream.Send(&pb.FileData{DataChunk: numBytes[:8]}); err != nil {
				return err
			}

			// for getByte := range numBytes {

			// }

			// Increments total read bytes
			totalBytesRead += 8
		}

	}

	return nil
}

func main() {

	lis, err := net.Listen("tcp", ":5003")

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, &fileTransferServer{})
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
