package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/muktar-gif/Project-2-611/fileproto"
)

type fileTransferServer struct {
	pb.UnimplementedFileServiceServer
}

func (s *fileTransferServer) GetFileChunk(fileRequest *pb.FileSegmentRequest, stream pb.FileService_GetFileChunkServer) error {
	return nil
}

func main() {

	lis, err := net.Listen("tcp", ":5003")

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
