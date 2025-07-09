package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type Logging struct {
	logger *slog.Logger
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewLogging(logger *slog.Logger) *Logging {
	return &Logging{
		logger: logger,
	}
}

func (l *Logging) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(rw, r)

		// Log request
		duration := time.Since(start)
		l.logger.Info("HTTP request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", duration.String(),
			"user_agent", r.UserAgent(),
			"remote_addr", r.RemoteAddr,
		)
	})
}