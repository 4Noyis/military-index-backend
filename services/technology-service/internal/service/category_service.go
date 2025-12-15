package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/repository"
	"github.com/4Noyis/military-index-backend/shared/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id uint) (*models.Category, error)
	GetByName(ctx context.Context, name string) (*models.Category, error)
	Create(ctx context.Context, category *models.Category) error
	Update(ctx context.Context, id uint, category *models.Category) error
	Delete(ctx context.Context, id uint) error
}

// categoryService implements CategoryService
type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// GetAll retrieves all categories
func (s *categoryService) GetAll(ctx context.Context) ([]models.Category, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		return nil, err
	}

	return categories, nil
}

// GetByID retrieves a category by ID
func (s *categoryService) GetByID(ctx context.Context, id uint) (*models.Category, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid category ID")
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching category by ID %d: %v", id, err)
		return nil, err
	}

	return category, nil
}

// GetByName retrieves a category by name
func (s *categoryService) GetByName(ctx context.Context, name string) (*models.Category, error) {
	if name == "" {
		return nil, fmt.Errorf("category name cannot be empty")
	}

	category, err := s.repo.GetByName(ctx, name)
	if err != nil {
		log.Printf("Error fetching category by name %s: %v", name, err)
		return nil, err
	}

	return category, nil
}

// Create creates a new category
func (s *categoryService) Create(ctx context.Context, category *models.Category) error {
	// validate required fields
	if err := s.validatesCategory(category); err != nil {
		return err
	}

	// Normalize category name all words first letter uppercase
	category.Name = cases.Title(language.English).String(strings.ToLower(strings.TrimSpace(category.Name)))

	if err := s.repo.Create(ctx, category); err != nil {
		log.Printf("Error creating category: %v", err)
		return err
	}

	log.Printf("Category created: %s", category.Name)
	return nil
}

// Update updates an existing category
func (s *categoryService) Update(ctx context.Context, id uint, category *models.Category) error {
	if id == 0 {
		return fmt.Errorf("invalid category ID")
	}

	// check if category exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Validate updated fields
	if err := s.validatesCategory(category); err != nil {
		return err
	}

	// Update fields - only update if provided (non-empty)
	existing.Name = cases.Title(language.English).String(strings.ToLower(strings.TrimSpace(category.Name)))

	// Only update optional fields if they are provided
	if category.Description != "" {
		existing.Description = strings.TrimSpace(category.Description)
	}
	if category.IconURL != "" {
		existing.IconURL = strings.TrimSpace(category.IconURL)
	}

	// Validate updated fields
	if err := s.repo.Update(ctx, existing); err != nil {
		log.Printf("Error updating category ID %d: %v", id, err)
		return err
	}

	log.Printf("Category updated ID %d", id)
	return nil
}

// Delete deletes a category
func (s *categoryService) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid category ID")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Printf("Error deleting category ID %d: %v", id, err)
		return err
	}

	log.Printf("Category deleted: ID %d", id)
	return nil
}

// validateCategory validates category data
func (s *categoryService) validatesCategory(category *models.Category) error {
	if category == nil {
		return fmt.Errorf("category data is required")
	}

	if strings.TrimSpace(category.Name) == "" {
		return fmt.Errorf("category name is required")
	}

	// Description and IconURL are optional fields
	return nil
}
