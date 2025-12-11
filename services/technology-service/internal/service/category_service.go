package service

import (
	"context"
	"fmt"
	"log"

	"github.com/4Noyis/military-index-backend/services/technology-service/internal/repository"
	"github.com/4Noyis/military-index-backend/shared/models"
)

// CategoryService defines the interface for category business logic
type CategoryService interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id uint) (*models.Category, error)
	GetByName(ctx context.Context, name string) (*models.Category, error)
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
