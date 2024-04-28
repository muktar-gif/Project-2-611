package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/muktar-gif/Project-2-611/fileproto"
)

type fileTransferServer struct{}

func (s *fileTransferServer) GetFileChunk(fileRequesst *pb.FileSegmentRequest, stream *pb.FileData) error {
	return nil
}

func main() {

	lis, err := net.Listen("tcp", ":5000")

	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	err = grpcServer.Serve(lis)

	pb.RegisterFileServiceServer(grpcServer, &fileTransferServer{})

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
