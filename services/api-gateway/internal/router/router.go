package router

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/4Noyis/military-index-backend/services/api-gateway/internal/proxy"
)

// Gateway represents the API Gateway
type Gateway struct {
	technologyProxy *proxy.ServiceProxy
	countryProxy    *proxy.ServiceProxy
}

// NewGateway creates a new API Gateway with configured proxies
func NewGateway() (*Gateway, error) {
	// Get service URLs from environment variables
	technologyServiceURL := getEnv("TECHNOLOGY_SERVICE_URL", "http://localhost:8081")
	countryServiceURL := getEnv("COUNTRY_SERVICE_URL", "http://localhost:8082")

	// Create proxies
	technologyProxy, err := proxy.NewServiceProxy("technology-service", technologyServiceURL)
	if err != nil {
		return nil, err
	}

	countryProxy, err := proxy.NewServiceProxy("country-service", countryServiceURL)
	if err != nil {
		return nil, err
	}

	log.Printf("Technology Service proxy: %s", technologyServiceURL)
	log.Printf("Country Service proxy: %s", countryServiceURL)

	return &Gateway{
		technologyProxy: technologyProxy,
		countryProxy:    countryProxy,
	}, nil
}

// SetupRoutes configures all routes for the API Gateway
func (g *Gateway) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Route to Technology Service
	// All /api/v1/technologies/* routes go to technology-service
	mux.HandleFunc("/api/v1/technologies/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Routing %s %s to technology-service", r.Method, r.URL.Path)
		g.technologyProxy.ServeHTTP(w, r)
	})

	// Route to Country Service
	// All /api/v1/countries/* routes go to country-service
	mux.HandleFunc("/api/v1/countries/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Routing %s %s to country-service", r.Method, r.URL.Path)
		g.countryProxy.ServeHTTP(w, r)
	})

	// Health check endpoint
	mux.HandleFunc("/health", g.healthCheck)

	// Root handler - must be last
	mux.HandleFunc("/", g.rootHandler)

	// Apply CORS middleware
	handler := corsMiddleware(mux)
	handler = loggingMiddleware(handler)

	return handler
}

// healthCheck returns the health status of the gateway
func (g *Gateway) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"service": "api-gateway",
		"version": "1.0.0",
	})
}

// rootHandler handles requests to the root path
func (g *Gateway) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Military Index API Gateway",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"health":       "/health",
			"technologies": "/api/v1/technologies",
			"countries":    "/api/v1/countries",
		},
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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

// loggingMiddleware logs all requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for health checks to reduce noise
		if !strings.HasPrefix(r.URL.Path, "/health") {
			log.Printf("Gateway: %s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
