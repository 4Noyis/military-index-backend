package service

import (
	"context"
	"fmt"
	"log"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/repository"
	"github.com/4Noyis/military-index-backend/shared/models"
)

// TechnologyService defines the interface for technology business logic
type TechnologyService interface {
	GetAll(ctx context.Context, page, limit int, sortBy, sortOrder string) ([]models.Technology, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Technology, error)
	GetByCountry(ctx context.Context, code string, page, limit int) ([]models.Technology, int64, error)
	GetByCategory(ctx context.Context, categoryName string, page, limit int) ([]models.Technology, int64, error)
	GetByStatus(ctx context.Context, status string, page, limit int) ([]models.Technology, int64, error)
	GetByYearRange(ctx context.Context, startYear, endYear int, page, limit int) ([]models.Technology, int64, error)
	Search(ctx context.Context, query string, page, limit int) ([]models.Technology, int64, error)
	Create(ctx context.Context, technology *models.Technology) error
	Update(ctx context.Context, id uint, technology *models.Technology) error
	Delete(ctx context.Context, id uint) error
}

// technologyService implements TechnologyService
type technologyService struct {
	repo repository.TechnologyRepository
}

// NewTechnologyService creates a new technology service
func NewTechnologyService(repo repository.TechnologyRepository) TechnologyService {
	return &technologyService{repo: repo}
}

// GetAll retrieves all technologies with pagination
func (s *technologyService) GetAll(ctx context.Context, page, limit int, sortBy, sortOrder string) ([]models.Technology, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	technologies, total, err := s.repo.GetAll(ctx, page, limit, sortBy, sortOrder)
	if err != nil {
		log.Printf("Error fetching technologies: %v", err)
		return nil, 0, err
	}

	return technologies, total, nil
}

// GetByID retrieves a technology by ID
func (s *technologyService) GetByID(ctx context.Context, id uint) (*models.Technology, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid technology ID")
	}

	technology, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching technology ID %d: %v", id, err)
		return nil, err
	}

	return technology, nil
}

// GetByCountry retrieves technologies by country code
func (s *technologyService) GetByCountry(ctx context.Context, code string, page, limit int) ([]models.Technology, int64, error) {
	// Validate country code
	if code == "" {
		return nil, 0, fmt.Errorf("country code cannot be empty")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	technologies, total, err := s.repo.GetByCountry(ctx, code, page, limit)
	if err != nil {
		log.Printf("Error fetching technologies for country %s: %v", code, err)
		return nil, 0, err
	}

	return technologies, total, nil
}

// GetByCategory retrieves technologies by category name
func (s *technologyService) GetByCategory(ctx context.Context, categoryName string, page, limit int) ([]models.Technology, int64, error) {
	// Validate category name
	if categoryName == "" {
		return nil, 0, fmt.Errorf("category name cannot be empty")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	technologies, total, err := s.repo.GetByCategory(ctx, categoryName, page, limit)
	if err != nil {
		log.Printf("Error fetching technologies for category %s: %v", categoryName, err)
		return nil, 0, err
	}

	return technologies, total, nil
}

// GetByStatus retrieves technologies by status
func (s *technologyService) GetByStatus(ctx context.Context, status string, page, limit int) ([]models.Technology, int64, error) {
	// Validate status
	validStatuses := map[string]bool{
		"historical": true,
		"current":    true,
		"future":     true,
		"concept":    true,
		"prototype":  true,
	}
	if !validStatuses[status] {
		return nil, 0, fmt.Errorf("invalid status: must be one of historical, current, future, concept, prototype")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	technologies, total, err := s.repo.GetByStatus(ctx, status, page, limit)
	if err != nil {
		log.Printf("Error fetching technologies with status %s: %v", status, err)
		return nil, 0, err
	}

	return technologies, total, nil
}

// GetByYearRange retrieves technologies by year range
func (s *technologyService) GetByYearRange(ctx context.Context, startYear, endYear int, page, limit int) ([]models.Technology, int64, error) {
	// Validate year range
	if startYear > endYear {
		return nil, 0, fmt.Errorf("start year cannot be greater than end year")
	}
	if startYear < 1900 || endYear > 2100 {
		return nil, 0, fmt.Errorf("year range must be between 1900 and 2100")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	technologies, total, err := s.repo.GetByYearRange(ctx, startYear, endYear, page, limit)
	if err != nil {
		log.Printf("Error fetching technologies for year range %d-%d: %v", startYear, endYear, err)
		return nil, 0, err
	}

	return technologies, total, nil
}

// Search searches technologies by name or description
func (s *technologyService) Search(ctx context.Context, query string, page, limit int) ([]models.Technology, int64, error) {
	// Validate query
	if query == "" {
		return nil, 0, fmt.Errorf("search query cannot be empty")
	}
	if len(query) < 2 {
		return nil, 0, fmt.Errorf("search query must be at least 2 characters")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	technologies, total, err := s.repo.Search(ctx, query, page, limit)
	if err != nil {
		log.Printf("Error searching technologies with query '%s': %v", query, err)
		return nil, 0, err
	}

	return technologies, total, nil
}

// Create creates a new technology
func (s *technologyService) Create(ctx context.Context, technology *models.Technology) error {
	// Validate required fields
	if technology.Name == "" {
		return fmt.Errorf("technology name is required")
	}
	if technology.CountryID == 0 {
		return fmt.Errorf("country ID is required")
	}
	if technology.CategoryID == 0 {
		return fmt.Errorf("category ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"historical": true,
		"current":    true,
		"future":     true,
		"concept":    true,
		"prototype":  true,
	}
	if technology.Status != "" && !validStatuses[technology.Status] {
		return fmt.Errorf("invalid status: must be one of historical, current, future, concept, prototype")
	}

	// Validate year ranges
	if technology.YearDeveloped != nil && (*technology.YearDeveloped < 1900 || *technology.YearDeveloped > 2100) {
		return fmt.Errorf("year developed must be between 1900 and 2100")
	}
	if technology.YearDeployed != nil && (*technology.YearDeployed < 1900 || *technology.YearDeployed > 2100) {
		return fmt.Errorf("year deployed must be between 1900 and 2100")
	}
	if technology.YearDeveloped != nil && technology.YearDeployed != nil && *technology.YearDeveloped > *technology.YearDeployed {
		return fmt.Errorf("year developed cannot be greater than year deployed")
	}

	err := s.repo.Create(ctx, technology)
	if err != nil {
		log.Printf("Error creating technology: %v", err)
		return err
	}

	log.Printf("Technology created successfully: %s (ID: %d)", technology.Name, technology.ID)
	return nil
}

// Update updates an existing technology
func (s *technologyService) Update(ctx context.Context, id uint, technology *models.Technology) error {
	// Validate ID
	if id == 0 {
		return fmt.Errorf("invalid technology ID")
	}

	// Check if technology exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("technology not found")
	}

	// Validate required fields
	if technology.Name == "" {
		return fmt.Errorf("technology name is required")
	}
	if technology.CountryID == 0 {
		return fmt.Errorf("country ID is required")
	}
	if technology.CategoryID == 0 {
		return fmt.Errorf("category ID is required")
	}

	// Validate status
	validStatuses := map[string]bool{
		"historical": true,
		"current":    true,
		"future":     true,
		"concept":    true,
		"prototype":  true,
	}
	if technology.Status != "" && !validStatuses[technology.Status] {
		return fmt.Errorf("invalid status: must be one of historical, current, future, concept, prototype")
	}

	// Validate year ranges
	if technology.YearDeveloped != nil && (*technology.YearDeveloped < 1900 || *technology.YearDeveloped > 2100) {
		return fmt.Errorf("year developed must be between 1900 and 2100")
	}
	if technology.YearDeployed != nil && (*technology.YearDeployed < 1900 || *technology.YearDeployed > 2100) {
		return fmt.Errorf("year deployed must be between 1900 and 2100")
	}
	if technology.YearDeveloped != nil && technology.YearDeployed != nil && *technology.YearDeveloped > *technology.YearDeployed {
		return fmt.Errorf("year developed cannot be greater than year deployed")
	}

	// Set the ID to ensure we're updating the right record
	technology.ID = existing.ID
	technology.CreatedAt = existing.CreatedAt

	err = s.repo.Update(ctx, technology)
	if err != nil {
		log.Printf("Error updating technology ID %d: %v", id, err)
		return err
	}

	log.Printf("Technology updated successfully: %s (ID: %d)", technology.Name, technology.ID)
	return nil
}

// Delete deletes a technology by ID
func (s *technologyService) Delete(ctx context.Context, id uint) error {
	// Validate ID
	if id == 0 {
		return fmt.Errorf("invalid technology ID")
	}

	// Check if technology exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("technology not found")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		log.Printf("Error deleting technology ID %d: %v", id, err)
		return err
	}

	log.Printf("Technology deleted successfully (ID: %d)", id)
	return nil
}

