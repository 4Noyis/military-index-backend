package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorHandler middleware catches panics and handles errors
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				// Send error response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				response := map[string]interface{}{
					"success": false,
					"error": map[string]string{
						"code":    "INTERNAL_ERROR",
						"message": "An internal server error occurred",
					},
				}

				json.NewEncoder(w).Encode(response)
			}
		}()

		// Call next handler
		next.ServeHTTP(w, r)
	})
}
