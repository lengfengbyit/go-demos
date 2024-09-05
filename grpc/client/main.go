package main

import (
	"context"
	"lengfengbyit/go-demos/grpc/hello"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50001"
	defaultName = "world"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}

	defer conn.Close()

	c := hello.NewGreeterClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &hello.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v\n", err)
	}
	log.Printf("Greeting: %s\n", r.GetMessage())
}
