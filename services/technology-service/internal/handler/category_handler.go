package handler

import (
	"net/http"
	"strconv"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/service"
	"github.com/4Noyis/military-index-backend/shared/utils"
)

// CategoryHandler handles HTTP requests for categories
type CategoryHandler struct {
	service service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// GetAll handles GET /api/v1/categories
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.service.GetAll(ctx)
	if err != nil {
		utils.SendInternalError(w, "failed to fetch categories")
		return
	}

	utils.SendSuccess(w, http.StatusOK, categories)
}

// GetByID handles GET /api/v1/categories/:id
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse ID from URL
	idStr := r.URL.Path[len("/api/v1/categories/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "invalid category ID")
		return
	}

	category, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		utils.SendNotFound(w, "category not found")
		return
	}

	utils.SendSuccess(w, http.StatusOK, category)
}
