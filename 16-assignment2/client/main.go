package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "distributed-systems/16-assignment2/proto"
	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run main.go <port>")
	}
	port := os.Args[1]
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Commands: put <key> <value> | get <key> | list | exit")

	for {
		fmt.Print("> ")
		cmdLine, _ := reader.ReadString('\n')
		cmdLine = strings.TrimSpace(cmdLine)
		if cmdLine == "" {
			continue
		}

		parts := strings.Split(cmdLine, " ")
		switch parts[0] {
		case "put":
			if len(parts) != 3 {
				fmt.Println("Usage: put <key> <value>")
				continue
			}
			_, err := client.Put(context.Background(), &pb.KeyValue{Key: parts[1], Value: parts[2]})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Stored and replicated")
			}

		case "get":
			if len(parts) != 2 {
				fmt.Println("Usage: get <key>")
				continue
			}
			res, err := client.Get(context.Background(), &pb.Key{Key: parts[1]})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("%s = %s\n", parts[1], res.Value)
			}

		case "list":
			stream, err := client.List(context.Background(), &pb.Empty{})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			for {
				kv, err := stream.Recv()
				if err != nil {
					break
				}
				fmt.Printf("%s = %s\n", kv.Key, kv.Value)
			}

		case "exit":
			fmt.Println("Bye!")
			return

		default:
			fmt.Println("Commands: put/get/list/exit")
		}
	}
}
