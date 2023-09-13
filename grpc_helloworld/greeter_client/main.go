package main

import (
	"context"
	"fmt"

	hello "github.com/iziyang/grpc_helloworld/helloworld"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := hello.NewHelloServiceClient(conn)
	response, err := client.SayHello(context.Background(), &hello.HelloRequest{Name: "World"})
	if err != nil {
		panic(err)
	}

	fmt.Println(response.Message)
}
