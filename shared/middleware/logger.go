package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger middleware logs HTTP requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		startTime := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status
		}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Calculate request duration
		duration := time.Since(startTime)

		// Log request details
		log.Printf(
			"[%s] %s %s | Status: %d | Duration: %v | IP: %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			wrapped.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}
