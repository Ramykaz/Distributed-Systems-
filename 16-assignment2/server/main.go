package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	pb "distributed-systems/16-assignment2/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedKeyValueStoreServer
	mu    sync.Mutex
	store map[string]string
	peers []string
}

func (s *server) Put(ctx context.Context, kv *pb.KeyValue) (*pb.Ack, error) {
	s.mu.Lock()
	s.store[kv.Key] = kv.Value
	s.mu.Unlock()

	// Replicate to peers
	for _, peer := range s.peers {
		go func(peer string) {
			conn, err := grpc.Dial(peer, grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to connect to peer %s: %v", peer, err)
				return
			}
			defer conn.Close()

			client := pb.NewKeyValueStoreClient(conn)
			_, err = client.Replicate(context.Background(), kv)
			if err != nil {
				log.Printf("Replication to %s failed: %v", peer, err)
			}
		}(peer)
	}

	return &pb.Ack{Success: true}, nil
}

func (s *server) Get(ctx context.Context, key *pb.Key) (*pb.Value, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.store[key.Key]
	if !ok {
		return &pb.Value{Value: ""}, nil
	}
	return &pb.Value{Value: value}, nil
}

func (s *server) List(_ *pb.Empty, stream pb.KeyValueStore_ListServer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range s.store {
		if err := stream.Send(&pb.KeyValue{Key: k, Value: v}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Replicate(ctx context.Context, kv *pb.KeyValue) (*pb.Ack, error) {
	s.mu.Lock()
	s.store[kv.Key] = kv.Value
	s.mu.Unlock()
	log.Printf("Replicated: %s = %s", kv.Key, kv.Value)
	return &pb.Ack{Success: true}, nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: go run main.go <port> <peer1,peer2>")
	}

	port := os.Args[1]
	peerStr := os.Args[2]
	peers := []string{}
	if peerStr != "" {
		for _, p := range strings.Split(peerStr, ",") {
			peers = append(peers, fmt.Sprintf("localhost:%s", p))
		}
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterKeyValueStoreServer(s, &server{
		store: make(map[string]string),
		peers: peers,
	})

	log.Printf("Server running on port %s with peers %v", port, peers)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
