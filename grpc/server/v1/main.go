package main

import (
	"context"
	"fmt"
	pb "lengfengbyit/go-demos/grpc/greet"
	"log"
	"net"

	"google.golang.org/grpc"
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

// 实现一个最简单的 rpc 服务端，不使用认证方式
func main() {
	// 1. 创建一个 TCP 服务
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("tcp listen failed: %v", err)
	}

	// 2. 创建一个 RPC 服务
	rpcService := grpc.NewServer()

	// 3. 注册服务
	pb.RegisterGreetServiceServer(rpcService, &server{})

	// 4. 启动 RPC 服务
	log.Printf("rpc server running: %s", addr)
	err = rpcService.Serve(listen)

	// 5. 服务结束
	log.Fatal(err)
}
