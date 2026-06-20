package middleware

import (
	"context"
	"net/http"
	"strings"

	"kafe-omborxona/internal/service"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
	RoleKey     contextKey = "role"
)

func Auth(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			var tokenStr string
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenStr = parts[1]
				}
			}
			if tokenStr == "" {
				tokenStr = r.URL.Query().Get("token")
			}

			if tokenStr == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			claims, err := authSvc.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			userID := int(claims["user_id"].(float64))
			username := claims["username"].(string)
			role := claims["role"].(string)

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UsernameKey, username)
			ctx = context.WithValue(ctx, RoleKey, role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			for _, allowed := range roles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		})
	}
}

func CORS(frontendURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", frontendURL)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) int {
	if v, ok := r.Context().Value(UserIDKey).(int); ok {
		return v
	}
	return 0
}

func GetRole(r *http.Request) string {
	if v, ok := r.Context().Value(RoleKey).(string); ok {
		return v
	}
	return ""
}
