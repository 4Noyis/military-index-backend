package router

import (
	"net/http"

	"github.com/4Noyis/military-index-backend/services/country-service/internal/handler"
	"github.com/4Noyis/military-index-backend/shared/middleware"
	"github.com/gorilla/mux"
)

// SetupRouter configures all routes for the country service
func SetupRouter(countryHandler *handler.CountryHandler) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.Logger)
	router.Use(middleware.ErrorHandler)
	router.Use(middleware.CORS)

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Country routes
	api.HandleFunc("/countries", countryHandler.GetAll).Methods("GET", "OPTIONS")
	api.HandleFunc("/countries/{id:[0-9]+}", countryHandler.GetByID).Methods("GET", "OPTIONS")
	api.HandleFunc("/countries/code/{code:[A-Z]{2,3}}", countryHandler.GetByCode).Methods("GET", "OPTIONS")
	api.HandleFunc("/countries", countryHandler.Create).Methods("POST", "OPTIONS")
	api.HandleFunc("/countries/{id:[0-9]+}", countryHandler.Update).Methods("PUT", "OPTIONS")
	api.HandleFunc("/countries/{id:[0-9]+}", countryHandler.Delete).Methods("DELETE", "OPTIONS")

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"country-service"}`))
	}).Methods("GET")

	return router
}
