package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

type Recovery struct {
	logger *slog.Logger
}

func NewRecovery(logger *slog.Logger) *Recovery {
	return &Recovery{
		logger: logger,
	}
}

func (r *Recovery) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				r.logger.Error("Panic recovered",
					"error", err,
					"method", req.Method,
					"path", req.URL.Path,
					"stack", string(debug.Stack()),
				)

				// Send error response
				http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, req)
	})
}