package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	list, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Could not connect to %v", err)
	}

	grpcServer := grpc.NewServer()

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("Could not Serve: %v", err)
	}
}
