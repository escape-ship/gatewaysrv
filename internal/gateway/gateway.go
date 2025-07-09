package gateway

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/escape-ship/gatewaysrv/config"
)

type Gateway struct {
	config *config.Config
	logger *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (*Gateway, error) {
	return &Gateway{
		config: cfg,
		logger: logger,
	}, nil
}

func (g *Gateway) RegisterHealthChecks(mux *http.ServeMux) {
	mux.HandleFunc("/health", g.healthHandler)
	mux.HandleFunc("/ready", g.readinessHandler)
}

func (g *Gateway) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"gatewaysrv"}`))
}

func (g *Gateway) readinessHandler(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you would check if all downstream services are ready
	// For now, we'll just return OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready","service":"gatewaysrv"}`))
}

func (g *Gateway) Shutdown(ctx context.Context) error {
	g.logger.Info("Gateway shutting down")
	// Add any cleanup logic here
	return nil
}
