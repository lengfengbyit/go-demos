package main

import (
	"context"
	pb "lengfengbyit/go-demos/grpc/greet"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address     = ":9091"
	defaultName = "world"
)

// token 认证 实现 credentials.PerRPCCredentials 接口
type ClientTokenAuth struct {
}

func (a ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": "123456789", // 这里携带需要验证的信息，如： token, appkey 等
	}, nil
}

// 是否传输安全， 返回 true, 则需要配合证书使用，否则不用加证书
func (a ClientTokenAuth) RequireTransportSecurity() bool {
	return true
}

func main() {

	creds, _ := credentials.NewClientTLSFromFile("key/test.pem", "www.example.com")

	var opts = make([]grpc.DialOption, 0, 2)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}

	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Greet(ctx, &pb.GreetRequest{Name: name, Age: 20, Hobbies: []string{"go", "php"}})
	if err != nil {
		log.Fatalf("could not greet: %v\n", err)
	}
	log.Printf("Greeting: %s\n", r.GetMessage())
}
