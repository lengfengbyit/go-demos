package main

import (
	"context"
	"lengfengbyit/go-demos/grpc/greet"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const addr = ":9091"

type ClientTokenAuth struct{}

// GetRequestMetadata 获取请求元数据
func (a ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"token": "123123"}, nil
}

// RequireTransportSecurity 是否传输安全，返回 true, 则需要配合证书使用，否则不用加证书
func (a ClientTokenAuth) RequireTransportSecurity() bool { return true }

// 发起一个使用 SSL 加密的自定义认证的 RPC 调用
func main() {

	// 1.0 创建一个 insecure 的 TLS 证书
	// example.com 生成证书时在 openssl.cnf 文件中定义的域名，防止其他域名调用
	creds, err := credentials.NewClientTLSFromFile("./ssl/test.pem", "example.com")
	if err != nil {
		log.Fatalf("failed to create tls credentials %v", err)
	}

	// 1.1. 创建一个 grpc 客户端，使用 ssl 加密
	var opts = make([]grpc.DialOption, 0, 2)
	// ssl 加密
	opts = append(opts, grpc.WithTransportCredentials(creds))
	// 自定义认证
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))
	conn, err := grpc.NewClient(addr, opts...)
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
