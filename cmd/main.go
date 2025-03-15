package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	gw "github.com/escape-ship/gatewaysrv/proto/gen" // Update
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	// Create CORS handler to allow cross-origin requests
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Or you can restrict to specific origins (e.g., ["http://localhost:3000"])
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	orderEndpoint := "localhost:9091"
	err := gw.RegisterOrderHandlerFromEndpoint(ctx, mux, orderEndpoint, opts)
	if err != nil {
		return err
	}
	accountEndpoint := "localhost:9090"
	err = gw.RegisterAccountHandlerFromEndpoint(ctx, mux, accountEndpoint, opts)
	if err != nil {
		return err
	}

	fmt.Println("Serving gRPC-Gateway on http://0.0.0.0:8081")

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", corsHandler)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
