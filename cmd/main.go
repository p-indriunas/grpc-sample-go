package main

import (
	"context"
	"flag"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw "github.com/p-indriunas/grpc-sample-go/gen/go"
	echo "github.com/p-indriunas/grpc-sample-go/internal/echo"
)

//
// Based on tutorial article:
// https://medium.com/swlh/rest-over-grpc-with-grpc-gateway-for-go-9584bfcbb835
//

var (
	// command-line options:
	// gRPC server endpoint
	//grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8080", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	grpcServer := grpc.NewServer()
	gw.RegisterEchoServiceServer(grpcServer, new(echo.EchoServiceImpl))

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	router := runtime.NewServeMux()
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8081",
		grpc.WithInsecure(),
	)

	if err = gw.RegisterEchoServiceHandler(ctx, router, conn); err != nil {
		return err
	}

	//err = gw.RegisterEchoServiceHandlerServer(ctx, router, grpcServer)
	//if err != nil {
	//	return err
	//}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", httpGrpcRouter(grpcServer, router))
}

func httpGrpcRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
