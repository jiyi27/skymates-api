package middleware

import (
	"context"
	"net/http"
	"skymates-api/pkg/auth"
	"strings"
)

// Auth 身份验证中间件
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderValue := r.Header.Get("Authorization")
		if authHeaderValue == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeaderValue, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateJwtToken(bearerToken[1])
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "signature"):
				http.Error(w, err.Error(), http.StatusUnauthorized)
			case strings.Contains(err.Error(), "expired"):
				http.Error(w, err.Error(), http.StatusUnauthorized)
			default:
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
