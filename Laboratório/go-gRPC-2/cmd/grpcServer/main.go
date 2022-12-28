package main

import (
	"database/sql"
	"net"

	"example.com/m/internal/database"
	"example.com/m/internal/pb"
	"example.com/m/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// fazendo conexao com o banco
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	// dizendo para fechar a conexao quando tudo finalizar
	defer db.Close()

	// instanciando o repositorio
	categoryDb := database.NewCategory(db)
	// instanciando o servico
	categoryService := service.NewCategoryService(*categoryDb)

	// criando o grpc server
	grpcServer := grpc.NewServer()
	// registrando nosso servico no grpc server
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	// necessario para usar no cliente do server
	reflection.Register(grpcServer)

	// abrindo uma conexao tcp
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
