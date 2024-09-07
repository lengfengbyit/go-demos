package main

import (
	"context"
	"lengfengbyit/go-demos/grpc/greet"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const addr = ":9091"

// 发起一个最简单的 RPC 调用
func main() {
	// 1. 创建一个 grpc 客户端，使用非安全方式(不使用 ssl 加密)
	// insecure.NewCredentials() 返回一个非安全的证书，否则这里需要传入一个 ssl 证书
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("new client failed, err: %v\n", err)
	}
	defer conn.Close()

	// 2. 创建一个 rpc 客户端
	client := greet.NewGreetServiceClient(conn)

	// 3. 构造请求参数
	requestParams := &greet.GreetRequest{
		Name:    "ZhangSan",
		Age:     20,
		Hobbies: []string{"篮球", "乒乓球", "足球"},
	}

	// 4. 发送 rpc 请求
	response, err := client.Greet(context.Background(), requestParams)
	if err != nil {
		log.Fatalf("call failed, err: %v\n", err)
	}

	log.Printf("response: %v\n", response.GetMessage())
}
