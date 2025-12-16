package middleware

import (
	"net/http"
	"strings"

	"github.com/4Noyis/military-index-backend/shared/auth"
	"github.com/4Noyis/military-index-backend/shared/utils"
)

// JWTAuth middleware validates JWT tokens for protected routes
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendUnauthorized(w, "Missing authorization header")
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendUnauthorized(w, "Invalid authorization format. Use: Bearer <token>")
			return
		}

		// Extract token
		tokenString := parts[1]

		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			utils.SendUnauthorized(w, "Invalid or expired token")
			return
		}

		// Token is valid, you could add claims to request context here if needed
		// ctx := context.WithValue(r.Context(), "username", claims.Username)
		// r = r.WithContext(ctx)

		// Log successful authentication
		_ = claims // Use claims if needed

		// Proceed to next handler
		next.ServeHTTP(w, r)
	})
}
