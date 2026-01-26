package main

import (
	"CRM/internal/contact_service"
	"CRM/rpc/proto/contacts"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Define the port for the gRPC server
	grpcPort := ":50051"

	// Create a TCP listener on the specified port
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create an instance of our ContactService implementation
	contactServer := contact_service.NewServer()

	// Register the ContactService with the gRPC server
	rpc.RegisterContactServiceServer(s, contactServer)

	fmt.Printf("gRPC server listening on %s\n", grpcPort)

	// Start serving requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
