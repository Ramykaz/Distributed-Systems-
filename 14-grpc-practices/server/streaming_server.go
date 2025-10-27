package main

import (
	"log"
	"net"
	"time"

  pb "distributed-systems/14-grpc-practices/proto"

	"google.golang.org/grpc"
)

type streamerServer struct {
	pb.UnimplementedStreamerServer
	items []pb.Item
}

func (s *streamerServer) ListItems(_ *pb.Empty, stream pb.Streamer_ListItemsServer) error {
	for _, item := range s.items {
		if err := stream.Send(&item); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 500) // simulate delay
	}
	return nil
}

func main() {
	items := []pb.Item{
		{Key: "1", Value: "Apple"},
		{Key: "2", Value: "Banana"},
		{Key: "3", Value: "Cherry"},
	}

	lis, err := net.Listen("tcp", ":50062")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStreamerServer(grpcServer, &streamerServer{items: items})

	log.Println("Streaming gRPC server running on port 50062")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
