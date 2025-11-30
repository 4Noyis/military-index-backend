package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/4Noyis/military-index-backend/services/api-gateway/internal/router"
)

func main() {
	log.Println("Starting API Gateway...")

	// Create gateway with proxies
	gateway, err := router.NewGateway()
	if err != nil {
		log.Fatal("Failed to create gateway:", err)
	}

	// Setup routes
	handler := gateway.SetupRoutes()

	// Get port from environment or use default
	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("API Gateway listening on port %s", port)
		log.Printf("Health check: http://localhost:%s/health", port)
		log.Printf("Technologies API: http://localhost:%s/api/v1/technologies", port)
		log.Printf("Countries API: http://localhost:%s/api/v1/countries", port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("API Gateway stopped gracefully")
}
