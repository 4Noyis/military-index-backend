package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/4Noyis/military-index-backend/services/country-service/internal/service"
	"github.com/4Noyis/military-index-backend/shared/models"
	"github.com/4Noyis/military-index-backend/shared/utils"
	"github.com/gorilla/mux"
)

// CountryHandler handles HTTP requests for countries
type CountryHandler struct {
	service service.CountryService
}

// NewCountryHandler creates a new country handler
func NewCountryHandler(service service.CountryService) *CountryHandler {
	return &CountryHandler{service: service}
}

// GetAll handles GET /api/v1/countries
func (h *CountryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Get search query if provided
	query := r.URL.Query().Get("q")

	var countries []models.Country
	var total int64
	var err error

	if query != "" {
		// Search countries
		countries, total, err = h.service.Search(ctx, query, page, limit)
	} else {
		// Get all countries
		countries, total, err = h.service.GetAll(ctx, page, limit)
	}

	if err != nil {
		utils.SendInternalError(w, "Failed to fetch countries", err.Error())
		return
	}

	// Send paginated response
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	utils.SendPaginated(w, countries, page, limit, total)
}

// GetByID handles GET /api/v1/countries/{id}
func (h *CountryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "Invalid country ID")
		return
	}

	country, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		utils.SendNotFound(w, "Country not found")
		return
	}

	utils.SendSuccess(w, http.StatusOK, country)
}

// GetByCode handles GET /api/v1/countries/code/{code}
func (h *CountryHandler) GetByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	code := vars["code"]
	if code == "" {
		utils.SendBadRequest(w, "Country code is required")
		return
	}

	country, err := h.service.GetByCode(ctx, code)
	if err != nil {
		utils.SendNotFound(w, "Country not found")
		return
	}

	utils.SendSuccess(w, http.StatusOK, country)
}

// Create handles POST /api/v1/countries
func (h *CountryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var country models.Country
	if err := json.NewDecoder(r.Body).Decode(&country); err != nil {
		utils.SendBadRequest(w, "Invalid request body", err.Error())
		return
	}

	if err := h.service.Create(ctx, &country); err != nil {
		utils.SendBadRequest(w, "Failed to create country", err.Error())
		return
	}

	utils.SendCreated(w, country, "Country created successfully")
}

// Update handles PUT /api/v1/countries/{id}
func (h *CountryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "Invalid country ID")
		return
	}

	var country models.Country
	if err := json.NewDecoder(r.Body).Decode(&country); err != nil {
		utils.SendBadRequest(w, "Invalid request body", err.Error())
		return
	}

	if err := h.service.Update(ctx, uint(id), &country); err != nil {
		utils.SendBadRequest(w, "Failed to update country", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, country, "Country updated successfully")
}

// Delete handles DELETE /api/v1/countries/{id}
func (h *CountryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "Invalid country ID")
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		utils.SendNotFound(w, "Country not found")
		return
	}

	utils.SendNoContent(w)
}
