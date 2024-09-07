package main

import (
	"context"
	"log"
	"net"

	pb "github.com/akshaychikhalkar/GoTaskQueue_v2/tasks" // Correct import path

	"google.golang.org/grpc"
)

const (
	port = ":50051" // Port for gRPC server
)

// server is used to implement ConsumerService
type server struct {
	pb.UnimplementedConsumerServiceServer // Embed the correct unimplemented server struct
}

// Implement ConsumeTask method for ConsumerService
func (s *server) ConsumeTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	log.Printf("Consumed task: Type: %v, Value: %v", req.TaskType, req.TaskValue)
	return &pb.TaskResponse{Success: true}, nil
}

func main() {
	// Start the gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterConsumerServiceServer(grpcServer, &server{}) // Register ConsumerService

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
