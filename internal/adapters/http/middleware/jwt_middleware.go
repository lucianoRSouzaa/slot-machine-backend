package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	handler_error "slot-machine/internal/adapters/http/handler/error"
	"slot-machine/internal/domain/contextkeys"
	"slot-machine/internal/domain/ports"
	"strings"

	"github.com/gorilla/mux"
)

type ContextKey string


func JWTMiddleware(jwtManager ports.JWTManager) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(handler_error.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "Authorization header missing",
				})

				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(handler_error.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "Invalid Authorization header format. Expected 'Bearer <token>'",
				})

				return
			}

			tokenString := parts[1]

			claims, err := jwtManager.VerifyAccessToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(handler_error.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: "Invalid token",
				})

				return
			}

			ctx := context.WithValue(r.Context(), contextkeys.ContextKeyUserID, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(contextkeys.ContextKeyUserID).(string)
	if !ok {
		return "", errors.New("userID not found in context")
	}
	return userID, nil
}
