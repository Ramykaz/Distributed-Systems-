package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "distributed-systems/grpc-streaming/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewWeatherServiceClient(conn)

	stream, err := client.GetWeatherStream(context.Background(), &pb.CityRequest{City: "Ankara"})
	if err != nil {
		log.Fatalf("Error starting stream: %v", err)
	}

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("✅ Stream ended.")
			break
		}
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}
		fmt.Printf("🌤️ %s: %s, %.1f°C\n", update.City, update.Condition, update.Temperature)
		time.Sleep(1 * time.Second)
	}
}
