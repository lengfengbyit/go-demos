package main

import (
	"context"
	"lengfengbyit/go-demos/grpc/hello"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50001"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &hello.HelloReply{
		Message: "Hello " + in.GetName(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
