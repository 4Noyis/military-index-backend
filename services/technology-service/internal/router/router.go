package router

import (
	"net/http"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/handler"
	"github.com/4Noyis/military-index-backend/shared/middleware"
)

// SetupRoutes configures all routes for the technology service
func SetupRoutes(techHandler *handler.TechnologyHandler, categoryHandler *handler.CategoryHandler) http.Handler {
	mux := http.NewServeMux()

	// Apply middleware
	var handler http.Handler = mux

	// Category routes
	// GET /api/v1/categories - List all categories
	mux.HandleFunc("GET /api/v1/categories", categoryHandler.GetAll)

	// GET /api/v1/categories/{id} - Get single category
	mux.HandleFunc("GET /api/v1/categories/", categoryHandler.GetByID)

	// Technology routes
	// GET /api/v1/technologies - List all technologies (with optional search ?q=)
	mux.HandleFunc("GET /api/v1/technologies", techHandler.GetAll)

	// GET /api/v1/technologies/{id} - Get single technology
	mux.HandleFunc("GET /api/v1/technologies/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a specific route pattern
		path := r.URL.Path
		if len(path) > len("/api/v1/technologies/") {
			segment := path[len("/api/v1/technologies/"):]

			// Check for special routes
			if len(segment) > 8 && segment[:8] == "country/" {
				techHandler.GetByCountry(w, r)
				return
			}
			if len(segment) > 9 && segment[:9] == "category/" {
				techHandler.GetByCategory(w, r)
				return
			}
			if len(segment) > 7 && segment[:7] == "status/" {
				techHandler.GetByStatus(w, r)
				return
			}

			// Otherwise, treat as ID
			techHandler.GetByID(w, r)
		}
	})

	// GET /api/v1/technologies/country/{code} - Get by country
	mux.HandleFunc("GET /api/v1/technologies/country/", techHandler.GetByCountry)

	// GET /api/v1/technologies/category/{name} - Get by category
	mux.HandleFunc("GET /api/v1/technologies/category/", techHandler.GetByCategory)

	// GET /api/v1/technologies/status/{status} - Get by status
	mux.HandleFunc("GET /api/v1/technologies/status/", techHandler.GetByStatus)

	// GET /api/v1/technologies/year-range?start=2000&end=2020 - Get by year range
	mux.HandleFunc("GET /api/v1/technologies/year-range", techHandler.GetByYearRange)

	// POST /api/v1/technologies - Create technology
	mux.HandleFunc("POST /api/v1/technologies", techHandler.Create)

	// PUT /api/v1/technologies/{id} - Update technology
	mux.HandleFunc("PUT /api/v1/technologies/", techHandler.Update)

	// DELETE /api/v1/technologies/{id} - Delete technology
	mux.HandleFunc("DELETE /api/v1/technologies/", techHandler.Delete)

	// Apply middleware (CORS, Logger, Error Handler)
	handler = middleware.CORS(handler)
	handler = middleware.Logger(handler)
	handler = middleware.ErrorHandler(handler)

	return handler
}
