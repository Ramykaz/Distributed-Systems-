package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "distributed-systems/grpc-streaming/proto"
	"google.golang.org/grpc"
)

type weatherServer struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *weatherServer) GetWeatherStream(req *pb.CityRequest, stream pb.WeatherService_GetWeatherStreamServer) error {
	city := req.City
	conditions := []string{"Sunny", "Cloudy", "Rainy", "Windy", "Stormy"}

	for i, condition := range conditions {
		update := &pb.WeatherUpdate{
			City:        city,
			Condition:   condition,
			Temperature: 20 + float32(i)*1.5,
		}

		if err := stream.Send(update); err != nil {
			return err
		}

		fmt.Printf("Sent update %d for %s\n", i+1, city)
		time.Sleep(2 * time.Second)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterWeatherServiceServer(grpcServer, &weatherServer{})

	fmt.Println("🌦️  Weather gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
