package main

import (
	"context"
	"fmt"
	pb "lengfengbyit/go-demos/grpc/greet"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	addr = ":9091"
)

type server struct {
	pb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.GreetResponse{
		Message: fmt.Sprintf("Hello %s, age: %d, hobbits: %v", in.GetName(), in.GetAge(), in.GetHobbies()),
	}, nil
}

// 实现一个 SSL 加密的 rpc 服务端
func main() {
	// 1. 创建一个 TCP 服务
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("tcp listen failed: %v", err)
	}

	// 2.0 创建一个 insecure 的 TLS 证书
	creds, err := credentials.NewServerTLSFromFile("./ssl/test.pem", "./ssl/test.key")
	if err != nil {
		log.Fatalf("insecure tls failed: %v", err)
	}

	// 2.1 创建一个使用 SSL 加密 RPC 服务
	rpcService := grpc.NewServer(grpc.Creds(creds))

	// 3. 注册服务
	pb.RegisterGreetServiceServer(rpcService, &server{})

	// 4. 启动 RPC 服务
	log.Printf("rpc server running: %s", addr)
	err = rpcService.Serve(listen)

	// 5. 服务结束
	log.Fatal(err)
}
