package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	handler_error "slot-machine/internal/adapters/http/handler/error"
	"slot-machine/internal/domain/contextkeys"
	"slot-machine/internal/domain/ports"
	"slot-machine/internal/infrastructure/config"
)

func AdminMiddleware(jwtManager ports.JWTManager) func(http.Handler) http.Handler {
	adminSecret := config.GetRequiredEnv("ADMIN_SECRET")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader != "" && strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				JWTMiddleware(jwtManager)(next).ServeHTTP(w, r)
				return
			}

			requestSecret := r.Header.Get("X-Admin-Secret")
			if requestSecret != "" && requestSecret == adminSecret {
				ctx := context.WithValue(r.Context(), contextkeys.ContextKeyUserID, "admin")
				ctx = context.WithValue(ctx, contextkeys.ContextKeyIsAdmin, true)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(handler_error.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
			})
		})
	}
}
