package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	jwtSecret string
}

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuth(jwtSecret string) *Auth {
	return &Auth{
		jwtSecret: jwtSecret,
	}
}

func (a *Auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 인증이 필요 없는 경로 예외 처리
		nonAuthPaths := map[string]bool{
			"/products": true,
			"/health":   true,
			"/ready":    true,
			"/register": true,
			"/login":    true,
		}
		if r.Method == http.MethodOptions || nonAuthPaths[r.URL.Path] || strings.HasPrefix(r.URL.Path, "/oauth/") {
			next.ServeHTTP(w, r)
			return
		}

		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if token has Bearer prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(a.jwtSecret), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
	user, ok := ctx.Value(UserContextKey).(*UserClaims)
	return user, ok
}
