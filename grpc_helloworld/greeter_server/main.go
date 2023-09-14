package main

import (
	"context"
	"fmt"
	"net"

	hello "github.com/iziyang/grpc_helloworld/helloworld"
	"google.golang.org/grpc"
)

type server struct {
	hello.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	hello.RegisterHelloServiceServer(grpcServer, &server{})

	fmt.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
