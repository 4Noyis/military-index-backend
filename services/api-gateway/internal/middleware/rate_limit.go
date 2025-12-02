package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TokenBucket represents a token bucket for rate limiting
type TokenBucket struct {
	tokens         float64
	capacity       float64
	refillRate     float64 // tokens per second
	lastRefillTime time.Time
	mutex          sync.Mutex
}

// RateLimiter manages rate limiting for multiple clients
type RateLimiter struct {
	buckets       map[string]*TokenBucket
	mutex         sync.RWMutex
	capacity      float64
	refillRate    float64
	cleanupTicker *time.Ticker
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	// Get configuration from environment variables
	requestsPerMinute := getEnvInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60)
	
	capacity := float64(requestsPerMinute)
	refillRate := capacity / 60.0 // tokens per second

	log.Printf("Rate Limiter initialized: %d requests/minute (%.2f req/sec)", requestsPerMinute, refillRate)

	rl := &RateLimiter{
		buckets:    make(map[string]*TokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}

	// Start cleanup goroutine to remove old entries
	rl.startCleanup()

	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	bucket, exists := rl.buckets[ip]
	if !exists {
		bucket = &TokenBucket{
			tokens:         rl.capacity,
			capacity:       rl.capacity,
			refillRate:     rl.refillRate,
			lastRefillTime: time.Now(),
		}
		rl.buckets[ip] = bucket
	}
	rl.mutex.Unlock()

	return bucket.consume()
}

// consume attempts to consume one token from the bucket
func (tb *TokenBucket) consume() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Seconds()
	tb.tokens += elapsed * tb.refillRate

	// Cap tokens at capacity
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	tb.lastRefillTime = now

	// Try to consume a token
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// startCleanup starts a background goroutine to clean up old entries
func (rl *RateLimiter) startCleanup() {
	rl.cleanupTicker = time.NewTicker(5 * time.Minute)
	
	go func() {
		for range rl.cleanupTicker.C {
			rl.cleanup()
		}
	}()
}

// cleanup removes buckets that haven't been used recently
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	count := 0

	for ip, bucket := range rl.buckets {
		bucket.mutex.Lock()
		if bucket.lastRefillTime.Before(cutoff) {
			delete(rl.buckets, ip)
			count++
		}
		bucket.mutex.Unlock()
	}

	if count > 0 {
		log.Printf("Rate limiter cleanup: removed %d inactive buckets", count)
	}
}

// Stop stops the cleanup ticker
func (rl *RateLimiter) Stop() {
	if rl.cleanupTicker != nil {
		rl.cleanupTicker.Stop()
	}
}

// Middleware returns a rate limiting middleware
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		ip := getClientIP(r)

		// Check if request is allowed
		if !rl.Allow(ip) {
			// Rate limit exceeded
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-RateLimit-Limit", strconv.FormatFloat(rl.capacity, 'f', 0, 64))
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error": map[string]string{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
				},
			})
			return
		}

		// Request allowed, proceed
		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr, strip port
	ip := r.RemoteAddr
	// Remove port if present (IPv4: "IP:port", IPv6: "[IP]:port")
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		// Check if it's IPv6 with brackets
		if strings.HasPrefix(ip, "[") && strings.Contains(ip[idx:], "]") {
			// IPv6 format [::1]:port - keep everything before the last colon
			ip = ip[:idx]
			// Remove brackets
			ip = strings.Trim(ip, "[]")
		} else {
			// IPv4 format 127.0.0.1:port - keep everything before the last colon
			ip = ip[:idx]
		}
	}
	return ip
}

// getEnvInt gets an integer environment variable or returns default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
