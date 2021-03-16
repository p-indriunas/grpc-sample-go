package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	gw "github.com/p-indriunas/grpc-sample-go/gen/go"
	"github.com/p-indriunas/grpc-sample-go/internal/services"
)

//
// Based on tutorial article:
// https://medium.com/swlh/rest-over-grpc-with-grpc-gateway-for-go-9584bfcbb835
//

func registerGrpcHandlers(grpcServer *grpc.Server) {
	gw.RegisterEchoServiceServer(grpcServer, services.NewEchoServiceServer())

	// Register GRPC services here
	// ...
}

func startGrpcServer(grpcServer *grpc.Server, address string, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		lis, err := net.Listen("tcp4", address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return
		}

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
			return
		}
	}()
}

func startHttpServer(address string, handler http.Handler, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := http.ListenAndServe(address, handler); err != nil {
			log.Fatalf("failed to serve http: %v", err)
			return
		}
	}()
}

func registerHttpGatewayHandlers(ctx context.Context, grpcAddress string, mux *runtime.ServeMux) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if err := gw.RegisterEchoServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return err
	}

	// Register additional service handlers here:
	// ...

	return nil
}

func httpGrpcRouter(grpcHandler http.Handler, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcHandler.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}

	grpcAddress := "localhost:5000"
	grpcServer := grpc.NewServer()
	startGrpcServer(grpcServer, grpcAddress, wg)

	registerGrpcHandlers(grpcServer)

	httpAddress := "localhost:5001"
	httpServer := runtime.NewServeMux()
	startHttpServer(httpAddress, httpGrpcRouter(grpcServer, httpServer), wg)

	registerHttpGatewayHandlers(ctx, grpcAddress, httpServer)

	wg.Wait()
	grpcServer.GracefulStop()
}
