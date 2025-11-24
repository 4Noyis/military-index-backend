package repository

import (
	"context"
	"fmt"

	"github.com/4Noyis/military-index-backend/shared/models"
	"gorm.io/gorm"
)

// TechnologyRepository defines the interface for tech
type TechnologyRepository interface {
	GetAll(ctx context.Context, page, limit int) ([]models.Technology, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Technology, error)
	GetByCountry(ctx context.Context, code string) ([]models.Technology, int64, error)
	GetByCategory(ctx context.Context, categoryName string) ([]models.Technology, int64, error)
	GetByStatus(ctx context.Context, status string) ([]models.Technology, int64, error)
	GetByYearRange(ctx context.Context, startYear, endYear int) ([]models.Technology, int64, error)
	Create(ctx context.Context, technology *models.Technology) error
	Update(ctx context.Context, technology *models.Technology) error
	Delete(ctx context.Context, id uint) error
}

// technologyRepository implements TechnologyRepository
type technologyRepository struct {
	db *gorm.DB
}

// NewTechnologyRepository creates a new technology repository
func NewTechnologyRepository(db *gorm.DB) TechnologyRepository {
	return &technologyRepository{db: db}
}

func (r *technologyRepository) GetAll(ctx context.Context, page, limit int) ([]models.Technology, int64, error) {
	var technologies []models.Technology
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&models.Technology{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count technologies: %w", err)
	}

	// calculate offset
	offset := (page - 1) * limit

	// Fetch paginated results with relationships
	if err := r.db.WithContext(ctx).
		Preload("Country").
		Preload("Category").
		Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&technologies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch technologies: %w", err)
	}

	return technologies, total, nil
}

// GetByID retrieves a technology by ID
func (r *technologyRepository) GetByID(ctx context.Context, id uint) (*models.Technology, error) {
	var technology models.Technology

	if err := r.db.WithContext(ctx).
		Preload("Country").
		Preload("Category").
		First(&technology, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("technology not found")
		}
		return nil, fmt.Errorf("failed to fetch technology: %w", err)
	}

	return &technology, nil
}

// GetByCountry retrieves technologies by country code
func (r *technologyRepository) GetByCountry(ctx context.Context, code string) ([]models.Technology, int64, error) {
	var technologies []models.Technology
	var total int64

	// Join with countries table to filter by code
	query := r.db.WithContext(ctx).
		Joins("JOIN countries ON countries.id = technologies.country_id").
		Where("countries.code = ?", code)

	// Count matching records
	if err := query.Model(&models.Technology{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count technologes: %w", err)
	}

	// Fetch technologies with related Country and Category
	if err := query.
		Preload("Country").
		Preload("Category").
		Order("technologies.name ASC").
		Find(&technologies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch technologies: %w", err)
	}

	return technologies, total, nil
}

// GetByCategory retrieves technologies by category
func (r *technologyRepository) GetByCategory(ctx context.Context, categoryName string) ([]models.Technology, int64, error) {
	var technologies []models.Technology
	var total int64

	query := r.db.WithContext(ctx).
		Joins("JOIN tech_categories ON tech_categories.id = technologies.category_id").
		Where("tech_categories.name = ?", categoryName)

	// Count matching records
	if err := query.Model(&models.Technology{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch technologies: %w", err)
	}

	// Fetch technologies with related category
	if err := query.Preload("Country").Preload("Category").Order("technologies.name ASC").Find(&technologies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch technologies: %w", err)
	}

	return technologies, total, nil
}

// GetByStatus retrieves technologies by status. status: 'historical', 'current', 'future', 'concept', 'prototype'
func (r *technologyRepository) GetByStatus(ctx context.Context, status string) ([]models.Technology, int64, error) {
	var technologies []models.Technology
	var total int64

	// Count matching records
	if err := r.db.WithContext(ctx).
		Model(&models.Technology{}).
		Where("status = ?", status).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count technologies: %w", err)
	}

	// Fetch results with relationships
	if err := r.db.WithContext(ctx).
		Preload("Country").
		Preload("Category").
		Where("status = ?", status).
		Order("name ASC").
		Find(&technologies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch technologies: %w", err)
	}

	return technologies, total, nil
}

// GetByYearRange retrieves technologies by year range. exp: all developed in range of 2002-2025 technologies
func (r *technologyRepository) GetByYearRange(ctx context.Context, startYear, endYear int) ([]models.Technology, int64, error) {
	var technologies []models.Technology
	var total int64

	// Count matching records
	if err := r.db.WithContext(ctx).
		Model(&models.Technology{}).
		Where("year_developed >= ? AND year_developed <= ?", startYear, endYear).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count technologies: %w", err)
	}

	// Fetch results with relationships
	if err := r.db.WithContext(ctx).
		Preload("Country").
		Preload("Category").
		Where("year_developed >= ? AND year_developed <= ?", startYear, endYear).
		Order("year_developed ASC").
		Find(&technologies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch technologies: %w", err)
	}

	return technologies, total, nil
}

// Create creates a new technology
func (r *technologyRepository) Create(ctx context.Context, technology *models.Technology) error {
	if err := r.db.WithContext(ctx).Create(technology).Error; err != nil {
		return fmt.Errorf("Failed to create country: %w", err)
	}
	return nil
}

// Update updates an existing technology
func (r *technologyRepository) Update(ctx context.Context, technology *models.Technology) error {
	if err := r.db.WithContext(ctx).Save(technology).Error; err != nil {
		return fmt.Errorf("failed to update technology: %w", err)
	}
	return nil
}

// Delete deletes a technology by ID
func (r *technologyRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Technology{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete technology: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("technology not found")
	}

	return nil
}
