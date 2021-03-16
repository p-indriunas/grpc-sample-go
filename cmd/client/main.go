package main

import (
	"context"
	"fmt"
	gw "github.com/p-indriunas/grpc-sample-go/gen/go"
	"google.golang.org/grpc"
	"log"
	"os"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	grpcAddress := "localhost:5000"
	conn, err := grpc.Dial(grpcAddress, opts...)
	if err != nil {
		log.Fatalf("failed to dial grpc: %v", err)
	}
	defer conn.Close()

	client := gw.NewEchoServiceClient(conn)

	message := "hello grpc"
	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	request := &gw.EchoRequest{
		Echo: message,
	}

	response, err := client.Echo(context.TODO(), request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Echo)
}
