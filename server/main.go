package main

import (
	"fmt"
	"log"
	"net"

	pb "gobook/pb"

	"google.golang.org/grpc"
)

const (
	defaultPort = 50051
)

func main() {
	port := defaultPort

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	gobookServer := NewGoBookServer()
	pb.RegisterGoBookRPCServer(grpcServer, gobookServer)

	fmt.Println("=========================================")
	fmt.Println("       GoBook Server - gRPC Lab")
	fmt.Println("=========================================")
	fmt.Printf("  Server listening on port %d...\n", port)
	fmt.Println("  Waiting for clients to connect...")
	fmt.Println("=========================================")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
