package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/github-carlos/Desafios-Full-Cycle/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	AddUserVerbose(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Carlos Eduardo",
		Email: "carlos@email.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make request to Server: %v", err)
	}

	fmt.Print(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Carlos Eduardo",
		Email: "carlos@email.com",
	}

	res, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make request to gRPC server: %v", err)
	}

	for {
		stream, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive message %v", err)
		}

		fmt.Println("Status:", stream.Status, "User:", stream.GetUser())
	}

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{Id: "1", Name: "Carlos 1", Email: "carloseduardo1@email.com"},
		&pb.User{Id: "2", Name: "Carlos 2", Email: "carloseduardo2@email.com"},
		&pb.User{Id: "3", Name: "Carlos 3", Email: "carloseduardo3@email.com"},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Could not create stream %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Could not close stream: %v", err)
	}

	fmt.Println(res)
}

func AddUsersStreamBidirectional(client pb.UserServiceClient) {
	stream, err := client.AddUserStereamBidirect(context.Background())

	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{Id: "1", Name: "Carlos 1", Email: "carloseduardo1@email.com"},
		&pb.User{Id: "2", Name: "Carlos 2", Email: "carloseduardo2@email.com"},
		&pb.User{Id: "3", Name: "Carlos 3", Email: "carloseduardo3@email.com"},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending User: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving response: %v", err)
				break
			}
			fmt.Printf("Receiving User %v with Status %v", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()
	<-wait
}
