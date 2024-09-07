package main

import (
	"context"
	"errors"
	"fmt"
	pb "lengfengbyit/go-demos/grpc/greet"
	"log"
	"net"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const (
	port = ":9091"
)

type server struct {
	pb.UnimplementedGreetServiceServer
}

// checkToken 检查请求 token
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("token not found")
	}

	if v, ok := md["token"]; ok && v[0] == "123456789" {
		return nil
	}

	return errors.New("token invalid")
}

func (s *server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Received: %v", in.GetName())
	return &pb.GreetResponse{
		Message: fmt.Sprintf("Hello %s, age: %d, hobbits: %v", in.GetName(), in.GetAge(), in.GetHobbies()),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 添加证书
	pemPath, _ := filepath.Abs("key/test.pem")
	keyPath, _ := filepath.Abs("key/test.key")

	creds, err := credentials.NewServerTLSFromFile(pemPath, keyPath)
	if err != nil {
		log.Fatalf("creds err:", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterGreetServiceServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
