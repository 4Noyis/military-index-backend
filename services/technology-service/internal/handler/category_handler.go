package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/service"
	"github.com/4Noyis/military-index-backend/shared/models"
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

// GetByID handles GET /api/v1/categories/{id}
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from URL path
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/v1/categories/")

	// Parse ID from URL
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

// Create handles POST /api/v1/categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.SendBadRequest(w, "Invalid request body", err.Error())
		return
	}

	if err := h.service.Create(ctx, &category); err != nil {
		utils.SendBadRequest(w, "Failed to create category", err.Error())
		return
	}

	utils.SendCreated(w, category, "Category created successfully")
}

// Update handles PUT /api/v1/categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from URL path
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/v1/categories/")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "Invalid category ID")
		return
	}

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.SendBadRequest(w, "Invalid request body", err.Error())
		return
	}

	if err := h.service.Update(ctx, uint(id), &category); err != nil {
		utils.SendBadRequest(w, "Failed to update category", err.Error())
		return
	}

	utils.SendSuccess(w, http.StatusOK, category, "Category updated successfully")
}

// Delete handles DELETE /api/v1/categories/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract ID from URL path
	path := r.URL.Path
	idStr := strings.TrimPrefix(path, "/api/v1/categories/")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendBadRequest(w, "Invalid category ID")
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		utils.SendNotFound(w, "Category not found")
		return
	}

	utils.SendNoContent(w)
}
