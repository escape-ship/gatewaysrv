package main

import (
	"context"
	"os"

	"github.com/escape-ship/gatewaysrv/config"
	"github.com/escape-ship/gatewaysrv/internal/app"
	"github.com/escape-ship/gatewaysrv/pkg/errors"
	"github.com/escape-ship/gatewaysrv/pkg/logger"
)

func main() {
	// Initialize logger for startup
	startupLogger := logger.New("info")

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		startupLogger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize application logger with configured level
	appLogger := logger.New(cfg.App.LogLevel)

	// Create application
	application, err := app.New(cfg)
	if err != nil {
		appLogger.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	// Create context for graceful shutdown
	ctx := context.Background()

	// Run application
	if err := application.Run(ctx); err != nil {
		if customErr, ok := err.(*errors.Error); ok {
			appLogger.Error("Application error", 
				"error", customErr.Error(),
				"code", customErr.Code,
				"stack", customErr.Stack,
			)
		} else {
			appLogger.Error("Application error", "error", err)
		}
		os.Exit(1)
	}

	appLogger.Info("Application stopped successfully")
}
