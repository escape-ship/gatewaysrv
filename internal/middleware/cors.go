package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/escape-ship/gatewaysrv/config"
)

type CORS struct {
	config config.CORS
}

func NewCORS(config config.CORS) *CORS {
	return &CORS{
		config: config,
	}
}

func (c *CORS) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		if c.isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Set other CORS headers
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.config.AllowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.config.AllowedHeaders, ", "))
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(c.config.MaxAge))

		if c.config.AllowedCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *CORS) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowedOrigin := range c.config.AllowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}

	return false
}
