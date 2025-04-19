package middleware

import (
	"context"
	"net/http"

	"richisntreal-backend/internal/api/auth"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

// AuthMiddleware injects an Authenticator and sets userID in context.
func AuthMiddleware(a auth.Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uid, err := a.Authenticate(r)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// FromContext retrieves the authenticated userID.
func FromContext(ctx context.Context) int64 {
	if v := ctx.Value(UserIDKey); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}
