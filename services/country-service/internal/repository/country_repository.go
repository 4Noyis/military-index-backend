package repository

import (
	"context"
	"fmt"

	"github.com/4Noyis/military-index-backend/shared/models"
	"gorm.io/gorm"
)

// CountryRepository defines the interface for country data operations
type CountryRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]models.Country, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Country, error)
	GetByCode(ctx context.Context, code string) (*models.Country, error)
	Search(ctx context.Context, query string, page, limit int) ([]models.Country, int64, error)
	Create(ctx context.Context, country *models.Country) error
	Update(ctx context.Context, country *models.Country) error
	Delete(ctx context.Context, id uint) error
}

// countryRepository implements CountryRepository
type countryRepository struct {
	db *gorm.DB
}

// NewCountryRepository creates a new country repository
func NewCountryRepository(db *gorm.DB) CountryRepository {
	return &countryRepository{db: db}
}

// GetAll retrieves all countries with pagination
func (r *countryRepository) GetAll(ctx context.Context, page, limit int) ([]models.Country, int64, error) {
	var countries []models.Country
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Country{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count countries: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results
	if err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&countries).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch countries: %w", err)
	}

	return countries, total, nil
}

// GetByID retrieves a country by ID
func (r *countryRepository) GetByID(ctx context.Context, id uint) (*models.Country, error) {
	var country models.Country

	if err := r.db.WithContext(ctx).First(&country, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to fetch country: %w", err)
	}

	return &country, nil
}

// GetByCode retrieves a country by code
func (r *countryRepository) GetByCode(ctx context.Context, code string) (*models.Country, error) {
	var country models.Country

	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&country).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("country not found")
		}
		return nil, fmt.Errorf("failed to fetch country: %w", err)
	}

	return &country, nil
}

// Search searches countries by name
func (r *countryRepository) Search(ctx context.Context, query string, page, limit int) ([]models.Country, int64, error) {
	var countries []models.Country
	var total int64

	searchPattern := "%" + query + "%"

	// Count matching records
	if err := r.db.WithContext(ctx).
		Model(&models.Country{}).
		Where("name ILIKE ?", searchPattern).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count countries: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Fetch matching results
	if err := r.db.WithContext(ctx).
		Where("name ILIKE ?", searchPattern).
		Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&countries).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search countries: %w", err)
	}

	return countries, total, nil
}

// Create creates a new country
func (r *countryRepository) Create(ctx context.Context, country *models.Country) error {
	if err := r.db.WithContext(ctx).Create(country).Error; err != nil {
		return fmt.Errorf("failed to create country: %w", err)
	}
	return nil
}

// Update updates an existing country
func (r *countryRepository) Update(ctx context.Context, country *models.Country) error {
	if err := r.db.WithContext(ctx).Save(country).Error; err != nil {
		return fmt.Errorf("failed to update country: %w", err)
	}
	return nil
}

// Delete deletes a country by ID
func (r *countryRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Country{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete country: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}
