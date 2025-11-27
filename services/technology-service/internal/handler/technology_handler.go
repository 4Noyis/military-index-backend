package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/service"
	"github.com/4Noyis/military-index-backend/shared/models"
	"github.com/4Noyis/military-index-backend/shared/utils"
)

// TechnologyHandler handles HTTP reqyest for technologies
type TechnologyHandler struct {
	service service.TechnologyService
}

// NewTechnologyHandler creates a new technology handler
func NewTechnologyHandler(service service.TechnologyService) *TechnologyHandler {
	return &TechnologyHandler{service: service}
}

// GetAll handles GET /api/v1/technologies
func (h *TechnologyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	sortBy := r.URL.Query().Get("sort")
	setOrder := r.URL.Query().Get("order")

	// Get search query if provided
	query := r.URL.Query().Get("q")

	var technologies []models.Technology
	var total int64
	var err error

	if query != "" {
		// Search technologies
		technologies, total, err = h.service.Search(ctx, query, page, limit)
	} else {
		// GetAll technologies
		technologies, total, err = h.service.GetAll(ctx, page, limit, sortBy, setOrder)
	}

	if err != nil {
		utils.SendInternalError(w, "failed to fetch technologies")
		return
	}
	// send paginated response
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	utils.SendPaginated(w, technologies, page, limit, total)
}

// GetByID handles GET /api/v1/technologies/:id
func (h *TechnologyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse ID from URL
	idStr := r.URL.Path[len("/api/v1/technologies/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "invalid technology ID")
		return
	}

	technology, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		utils.SendNotFound(w, "technology not found")
		return
	}

	utils.SendSuccess(w, http.StatusOK, technology)
}

// GetByCountry handles GET /api/v1/technologies/country/:code
func (h *TechnologyHandler) GetByCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse country code from URL
	code := r.URL.Path[len("/api/v1/technologies/country/"):]
	if code == "" {
		utils.SendBadRequest(w, "country code is required")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	technologies, total, err := h.service.GetByCountry(ctx, code, page, limit)
	if err != nil {
		utils.SendInternalError(w, "failed to fetch technologies")
		return
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	utils.SendPaginated(w, technologies, page, limit, total)
}

// GetByCategory handles GET /api/v1/technologies/category/:name
func (h *TechnologyHandler) GetByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse category name from URL
	categoryName := r.URL.Path[len("/api/v1/technologies/category/"):]
	if categoryName == "" {
		utils.SendBadRequest(w, "category name is required")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	technologies, total, err := h.service.GetByCategory(ctx, categoryName, page, limit)
	if err != nil {
		utils.SendInternalError(w, "failed to fetch technologies")
		return
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	utils.SendPaginated(w, technologies, page, limit, total)
}

// GetByStatus handles GET /api/v1/technologies/status/:status
func (h *TechnologyHandler) GetByStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse status from URL
	status := r.URL.Path[len("/api/v1/technologies/status/"):]
	if status == "" {
		utils.SendBadRequest(w, "status is required")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	technologies, total, err := h.service.GetByStatus(ctx, status, page, limit)
	if err != nil {
		utils.SendBadRequest(w, err.Error())
		return
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	utils.SendPaginated(w, technologies, page, limit, total)
}

// GetByYearRange handles GET /api/v1/technologies/year-range
func (h *TechnologyHandler) GetByYearRange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	startYearStr := r.URL.Query().Get("start")
	endYearStr := r.URL.Query().Get("end")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if startYearStr == "" || endYearStr == "" {
		utils.SendBadRequest(w, "start and end year are required")
		return
	}

	startYear, err := strconv.Atoi(startYearStr)
	if err != nil {
		utils.SendBadRequest(w, "invalid start year")
		return
	}

	endYear, err := strconv.Atoi(endYearStr)
	if err != nil {
		utils.SendBadRequest(w, "invalid end year")
		return
	}

	technologies, total, err := h.service.GetByYearRange(ctx, startYear, endYear, page, limit)
	if err != nil {
		utils.SendBadRequest(w, err.Error())
		return
	}

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	utils.SendPaginated(w, technologies, page, limit, total)
}

// Create handles POST /api/v1/technologies
func (h *TechnologyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var technology models.Technology
	if err := json.NewDecoder(r.Body).Decode(&technology); err != nil {
		utils.SendBadRequest(w, "invalid request body", err.Error())
		return
	}

	err := h.service.Create(ctx, &technology)
	if err != nil {
		utils.SendBadRequest(w, err.Error())
		return
	}

	utils.SendCreated(w, technology, "Technology created successfully")
}

// Update handles PUT /api/v1/technologies/:id
func (h *TechnologyHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse ID from URL
	idStr := r.URL.Path[len("/api/v1/technologies/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "invalid technology ID")
		return
	}

	var technology models.Technology
	if err := json.NewDecoder(r.Body).Decode(&technology); err != nil {
		utils.SendBadRequest(w, "invalid request body", err.Error())
		return
	}

	err = h.service.Update(ctx, uint(id), &technology)
	if err != nil {
		if err.Error() == "technology not found" {
			utils.SendNotFound(w, err.Error())
			return
		}
		utils.SendBadRequest(w, err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, technology)
}

// Delete handles DELETE /api/v1/technologies/:id
func (h *TechnologyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse ID from URL
	idStr := r.URL.Path[len("/api/v1/technologies/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "invalid technology ID")
		return
	}

	err = h.service.Delete(ctx, uint(id))
	if err != nil {
		if err.Error() == "technology not found" {
			utils.SendNotFound(w, err.Error())
			return
		}
		utils.SendInternalError(w, "failed to delete technology")
		return
	}

	utils.SendSuccess(w, http.StatusOK, map[string]string{"message": "technology deleted successfully"})
}
