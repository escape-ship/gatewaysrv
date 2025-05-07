package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	accountEndpoint := "localhost:9090"
	err := gw.RegisterAccountHandlerFromEndpoint(ctx, mux, accountEndpoint, opts)
	if err != nil {
		return err
	}
	orderEndpoint := "localhost:9091"
	err = gw.RegisterProductServiceHandlerFromEndpoint(ctx, mux, orderEndpoint, opts)
	if err != nil {
		return err
	}
	paymentEndpoint := "localhost:9092"
	err = gw.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, paymentEndpoint, opts)
	if err != nil {
		return err
	}

	fmt.Println("Serving gRPC-Gateway on http://0.0.0.0:8081")
	s := &http.Server{
		Addr:    ":8081",
		Handler: allowCORS(mux),
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	// return http.ListenAndServe(":8081", authMiddleware(corsHandler))
	return s.ListenAndServe()
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-credentials", "true")
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	grpclog.Infof("Preflight request for %s", r.URL.Path)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
