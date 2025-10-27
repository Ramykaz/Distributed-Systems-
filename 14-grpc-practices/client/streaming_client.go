package main

import (
	"context"
	"io"
	"log"
	"time"

  pb "distributed-systems/14-grpc-practices/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50062", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamerClient(conn)

	stream, err := client.ListItems(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error calling ListItems: %v", err)
	}

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving item: %v", err)
		}
		log.Printf("Received: %s = %s", item.Key, item.Value)
		time.Sleep(time.Millisecond * 200)
	}
}
