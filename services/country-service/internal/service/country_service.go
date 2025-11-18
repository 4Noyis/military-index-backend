package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/4Noyis/military-index-backend/services/country-service/internal/repository"
	"github.com/4Noyis/military-index-backend/shared/models"
)

// CountryService defines the interface for country business logic
type CountryService interface {
	GetAll(ctx context.Context, page, limit int) ([]models.Country, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Country, error)
	GetByCode(ctx context.Context, code string) (*models.Country, error)
	Search(ctx context.Context, query string, page, limit int) ([]models.Country, int64, error)
	Create(ctx context.Context, country *models.Country) error
	Update(ctx context.Context, id uint, country *models.Country) error
	Delete(ctx context.Context, id uint) error
}

// countryService implements CountryService
type countryService struct {
	repo repository.CountryRepository
}

// NewCountryService creates a new country service
func NewCountryService(repo repository.CountryRepository) CountryService {
	return &countryService{repo: repo}
}

// GetAll retrieves all countries with pagination
func (s *countryService) GetAll(ctx context.Context, page, limit int) ([]models.Country, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	countries, total, err := s.repo.GetAll(ctx, page, limit)
	if err != nil {
		log.Printf("Error fetching countries: %v", err)
		return nil, 0, err
	}

	return countries, total, nil
}

// GetByID retrieves a country by ID
func (s *countryService) GetByID(ctx context.Context, id uint) (*models.Country, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid country ID")
	}

	country, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching country ID %d: %v", id, err)
		return nil, err
	}

	return country, nil
}

// GetByCode retrieves a country by code
func (s *countryService) GetByCode(ctx context.Context, code string) (*models.Country, error) {
	if code == "" {
		return nil, fmt.Errorf("country code is required")
	}

	// Normalize code to uppercase
	code = strings.ToUpper(strings.TrimSpace(code))

	country, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		log.Printf("Error fetching country code %s: %v", code, err)
		return nil, err
	}

	return country, nil
}

// Search searches countries by name
func (s *countryService) Search(ctx context.Context, query string, page, limit int) ([]models.Country, int64, error) {
	if query == "" {
		return nil, 0, fmt.Errorf("search query is required")
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	countries, total, err := s.repo.Search(ctx, query, page, limit)
	if err != nil {
		log.Printf("Error searching countries with query '%s': %v", query, err)
		return nil, 0, err
	}

	return countries, total, nil
}

// Create creates a new country
func (s *countryService) Create(ctx context.Context, country *models.Country) error {
	// Validate required fields
	if err := s.validateCountry(country); err != nil {
		return err
	}

	// Normalize code to uppercase
	country.Code = strings.ToUpper(strings.TrimSpace(country.Code))
	country.Name = strings.TrimSpace(country.Name)

	if err := s.repo.Create(ctx, country); err != nil {
		log.Printf("Error creating country: %v", err)
		return err
	}

	log.Printf("Country created: %s (%s)", country.Name, country.Code)
	return nil
}

// Update updates an existing country
func (s *countryService) Update(ctx context.Context, id uint, country *models.Country) error {
	if id == 0 {
		return fmt.Errorf("invalid country ID")
	}

	// Check if country exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Validate updated fields
	if err := s.validateCountry(country); err != nil {
		return err
	}

	// Update fields
	existing.Name = strings.TrimSpace(country.Name)
	existing.Code = strings.ToUpper(strings.TrimSpace(country.Code))
	existing.FlagURL = strings.TrimSpace(country.FlagURL)

	if err := s.repo.Update(ctx, existing); err != nil {
		log.Printf("Error updating country ID %d: %v", id, err)
		return err
	}

	log.Printf("Country updated: ID %d", id)
	return nil
}

// Delete deletes a country
func (s *countryService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid country ID")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting country ID %d: %v", id, err)
		return err
	}

	log.Printf("Country deleted: ID %d", id)
	return nil
}

// validateCountry validates country data
func (s *countryService) validateCountry(country *models.Country) error {
	if country == nil {
		return fmt.Errorf("country data is required")
	}

	if strings.TrimSpace(country.Name) == "" {
		return fmt.Errorf("country name is required")
	}

	if strings.TrimSpace(country.Code) == "" {
		return fmt.Errorf("country code is required")
	}

	// Validate code length (should be 2-3 characters)
	code := strings.TrimSpace(country.Code)
	if len(code) < 2 || len(code) > 3 {
		return fmt.Errorf("country code must be 2-3 characters")
	}

	return nil
}
