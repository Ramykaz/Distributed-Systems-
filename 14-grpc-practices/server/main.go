package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "distributed-systems/grpc-practices/proto"
	"google.golang.org/grpc"
)

// greeterServer implements the Greeter service
type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello handles the unary RPC call
func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received request for name: %s", req.Name)
	return &pb.HelloReply{Message: "Hello, " + req.Name + "!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50061")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &greeterServer{})

	fmt.Println("Unary gRPC server running on port 50061")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
