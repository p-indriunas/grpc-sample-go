package services

import (
	"context"
	grpc "github.com/p-indriunas/grpc-sample-go/gen/go"
)


type EchoServiceServer struct {
	grpc.UnimplementedEchoServiceServer
}

func NewEchoServiceServer() grpc.EchoServiceServer {
	return &EchoServiceServer{}
}

func (r *EchoServiceServer) Echo(ctx context.Context, request *grpc.EchoRequest) (*grpc.EchoResponse, error) {
	return &grpc.EchoResponse{Echo: request.Echo}, nil
}

func (r *EchoServiceServer) EchoStream(request *grpc.EchoRequest, responseStream grpc.EchoService_EchoStreamServer) error {
	return nil
}