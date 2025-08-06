package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/escape-ship/gatewaysrv/config"
	"github.com/escape-ship/gatewaysrv/internal/gateway"
	"github.com/escape-ship/gatewaysrv/internal/middleware"
	"github.com/escape-ship/gatewaysrv/pkg/logger"
	"github.com/escape-ship/protos/gen"
)

type App struct {
	config     *config.Config
	logger     *slog.Logger
	httpServer *http.Server
	gateway    *gateway.Gateway
}

func New(cfg *config.Config) (*App, error) {
	// Initialize logger
	logger := logger.New(cfg.App.LogLevel)

	// Create gateway
	gw, err := gateway.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway: %w", err)
	}

	app := &App{
		config:  cfg,
		logger:  logger,
		gateway: gw,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	// Setup HTTP server
	if err := a.setupHTTPServer(ctx); err != nil {
		return fmt.Errorf("failed to setup HTTP server: %w", err)
	}

	// Start server
	go func() {
		addr := fmt.Sprintf("%s:%d", a.config.App.Host, a.config.App.Port)
		a.logger.Info("Starting HTTP server", "address", addr)

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("HTTP server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	return a.gracefulShutdown()
}

func (a *App) setupHTTPServer(ctx context.Context) error {
	// Create gRPC gateway mux
	mux := runtime.NewServeMux()

	// Register service handlers
	if err := a.registerServiceHandlers(ctx, mux); err != nil {
		return fmt.Errorf("failed to register service handlers: %w", err)
	}

	// Setup middleware chain
	handler := a.setupMiddleware(mux)

	// Create HTTP server
	a.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", a.config.App.Host, a.config.App.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return nil
}

func (a *App) registerServiceHandlers(ctx context.Context, mux *runtime.ServeMux) error {
	// gRPC connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Register account service
	accountAddr := a.config.GetServiceAddress("account")
	if err := gen.RegisterAccountServiceHandlerFromEndpoint(ctx, mux, accountAddr, opts); err != nil {
		return fmt.Errorf("failed to register account service: %w", err)
	}

	// Register product service
	productAddr := a.config.GetServiceAddress("product")
	if err := gen.RegisterProductServiceHandlerFromEndpoint(ctx, mux, productAddr, opts); err != nil {
		return fmt.Errorf("failed to register product service: %w", err)
	}

	// Register payment service
	paymentAddr := a.config.GetServiceAddress("payment")
	if err := gen.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, paymentAddr, opts); err != nil {
		return fmt.Errorf("failed to register payment service: %w", err)
	}

	// Register order service
	orderAddr := a.config.GetServiceAddress("order")
	if err := gen.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, orderAddr, opts); err != nil {
		return fmt.Errorf("failed to register order service: %w", err)
	}

	a.logger.Info("All service handlers registered successfully")
	return nil
}

func (a *App) setupMiddleware(handler http.Handler) http.Handler {
	// Create middleware chain
	middlewares := []middleware.Middleware{
		middleware.NewLogging(a.logger),
		middleware.NewCORS(a.config.CORS),
		middleware.NewRecovery(a.logger),
		middleware.NewAuth(a.config.Auth.JWTSecret), // Auth는 마지막에 적용
	}

	// Apply middleware in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].Handle(handler)
	}

	return handler
}

func (a *App) gracefulShutdown() error {
	// Create channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	sig := <-quit
	a.logger.Info("Received shutdown signal", "signal", sig.String())

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Error("HTTP server shutdown error", "error", err)
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	// Shutdown gateway
	if err := a.gateway.Shutdown(ctx); err != nil {
		a.logger.Error("Gateway shutdown error", "error", err)
		return fmt.Errorf("gateway shutdown error: %w", err)
	}

	a.logger.Info("Server shutdown completed")
	return nil
}
