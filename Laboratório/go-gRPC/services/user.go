package services

import (
	"context"
	"fmt"

	"github.com/github-carlos/Desafios-Full-Cycle/pb"
)

// type UserServiceServer interface {
// 	AddUser(context.Context, *User) (*User, error)
// 	mustEmbedUnimplementedUserServiceServer()
// }

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (u *UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	fmt.Println("Saving Name", req.Name)

	return &pb.User{
		Id:    "123",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}
