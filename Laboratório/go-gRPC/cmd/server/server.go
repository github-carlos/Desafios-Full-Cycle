package main

import (
	"log"
	"net"

	"github.com/github-carlos/Desafios-Full-Cycle/pb"
	"github.com/github-carlos/Desafios-Full-Cycle/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	list, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Could not connect to %v", err)
	}

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	pb.RegisterUserServiceServer(grpcServer, &services.UserService{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("Could not Serve: %v", err)
	}
}
