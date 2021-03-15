package echo

import (
	"context"
	grpc "github.com/p-indriunas/grpc-sample-go/gen/go"
)


type EchoServiceImpl struct {
}

func (r *EchoServiceImpl) dd() {

}

func (r *EchoServiceImpl) Echo(context.Context, *grpc.EchoRequest) (*grpc.EchoResponse, error) {
	return nil, nil
}

func (r *EchoServiceImpl) EchoStream(*grpc.EchoRequest, grpc.EchoService_EchoStreamServer) error {
	return nil
}