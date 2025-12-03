package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/handler"
	"github.com/4Noyis/military-index-backend/services/technology-service/internal/repository"
	"github.com/4Noyis/military-index-backend/services/technology-service/internal/router"
	"github.com/4Noyis/military-index-backend/services/technology-service/internal/service"
	"github.com/4Noyis/military-index-backend/shared/config"
)

func main() {
	log.Println("Starting Technology Service...")

	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected")

	// Initialize layers
	technologyRepo := repository.NewTechnologyRepository(db)
	technologyService := service.NewTechnologyService(technologyRepo)
	technologyHandler := handler.NewTechnologyHandler(technologyService)

	// Setup router
	r := router.SetupRoutes(technologyHandler)

	// Get port from environment or use default
	port := os.Getenv("TECHNOLOGY_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Technology Service listening on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}
