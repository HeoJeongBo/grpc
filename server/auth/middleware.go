package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// Middleware validates JWT tokens and adds user information to the request context
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Allow unauthenticated requests to pass through
			// Individual handlers can check if user is authenticated
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(string)
	return userID, ok
}

// RequireAuth is a helper to check if a user is authenticated
func RequireAuth(ctx context.Context) (string, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok || userID == "" {
		return "", ErrUnauthorized
	}
	return userID, nil
}
