package repository

import (
	"context"
	"testing"

	"github.com/4Noyis/military-index-backend/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Simple test setup without relationship loading (avoids SQLite datetime issues)
func setupSimpleTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	err = db.AutoMigrate(&models.Country{}, &models.Category{}, &models.Technology{})
	require.NoError(t, err, "Failed to migrate test database")

	return db
}

func seedSimpleTestData(t *testing.T, db *gorm.DB) ([]models.Country, []models.Category, []models.Technology) {
	// Seed countries
	countries := []models.Country{
		{Name: "Turkey", Code: "TUR"},
		{Name: "United States", Code: "USA"},
		{Name: "Russia", Code: "RUS"},
	}
	for i := range countries {
		require.NoError(t, db.Create(&countries[i]).Error)
	}

	// Seed categories
	categories := []models.Category{
		{Name: "Aircraft"},
		{Name: "Naval"},
		{Name: "Ground Vehicles"},
		{Name: "Drones"},
	}
	for i := range categories {
		require.NoError(t, db.Create(&categories[i]).Error)
	}

	// Seed technologies
	year2023 := 2023
	year2024 := 2024
	year2010 := 2010
	year2015 := 2015

	technologies := []models.Technology{
		{Name: "KAAN Fighter", CountryID: countries[0].ID, CategoryID: categories[0].ID, YearDeveloped: &year2023, Status: "prototype"},
		{Name: "Bayraktar TB2", CountryID: countries[0].ID, CategoryID: categories[3].ID, YearDeveloped: &year2015, YearDeployed: &year2015, Status: "current"},
		{Name: "F-35 Lightning II", CountryID: countries[1].ID, CategoryID: categories[0].ID, YearDeveloped: &year2010, YearDeployed: &year2015, Status: "current"},
		{Name: "Su-57", CountryID: countries[2].ID, CategoryID: categories[0].ID, YearDeveloped: &year2010, Status: "current"},
		{Name: "Altay MBT", CountryID: countries[0].ID, CategoryID: categories[2].ID, YearDeveloped: &year2024, Status: "prototype"},
	}

	for i := range technologies {
		require.NoError(t, db.Create(&technologies[i]).Error)
	}

	return countries, categories, technologies
}

func TestSimpleTechnologyRepository_Create(t *testing.T) {
	db := setupSimpleTestDB(t)
	repo := NewTechnologyRepository(db)
	ctx := context.Background()

	countries, categories, _ := seedSimpleTestData(t, db)

	year := 2023
	tech := &models.Technology{
		Name:          "Test Aircraft",
		CountryID:     countries[0].ID,
		CategoryID:    categories[0].ID,
		Description:   "Test description",
		YearDeveloped: &year,
		Status:        "concept",
	}

	err := repo.Create(ctx, tech)
	assert.NoError(t, err)
	assert.NotZero(t, tech.ID)
}

func TestSimpleTechnologyRepository_Update(t *testing.T) {
	db := setupSimpleTestDB(t)
	repo := NewTechnologyRepository(db)
	ctx := context.Background()

	_, _, technologies := seedSimpleTestData(t, db)

	// Update the first technology (without loading relationships)
	var technology models.Technology
	err := db.First(&technology, technologies[0].ID).Error
	require.NoError(t, err)

	technology.Name = "Updated KAAN"
	technology.Description = "Updated description"

	err = repo.Update(ctx, &technology)
	assert.NoError(t, err)

	// Verify the update (without Preload)
	var updated models.Technology
	err = db.First(&updated, technology.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Updated KAAN", updated.Name)
	assert.Equal(t, "Updated description", updated.Description)
}

func TestSimpleTechnologyRepository_Delete(t *testing.T) {
	db := setupSimpleTestDB(t)
	repo := NewTechnologyRepository(db)
	ctx := context.Background()

	_, _, technologies := seedSimpleTestData(t, db)

	// Delete existing technology
	err := repo.Delete(ctx, technologies[0].ID)
	assert.NoError(t, err)

	// Verify deletion
	var deleted models.Technology
	err = db.First(&deleted, technologies[0].ID).Error
	assert.Error(t, err)

	// Delete non-existent technology
	err = repo.Delete(ctx, 9999)
	assert.Error(t, err)
}

func TestSimpleTechnologyRepository_GetByCountryCount(t *testing.T) {
	db := setupSimpleTestDB(t)
	repo := NewTechnologyRepository(db)
	ctx := context.Background()

	seedSimpleTestData(t, db)

	// Count Turkish technologies
	techs, total, err := repo.GetByCountry(ctx, "TUR", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total) // KAAN, TB2, Altay

	// Just verify we got results, don't check relationships
	assert.Len(t, techs, 3)
}

func TestSimpleTechnologyRepository_GetByStatus(t *testing.T) {
	db := setupSimpleTestDB(t)
	repo := NewTechnologyRepository(db)
	ctx := context.Background()

	seedSimpleTestData(t, db)

	// Get current status technologies
	techs, total, err := repo.GetByStatus(ctx, "current", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total) // TB2, F-35, Su-57
	assert.Len(t, techs, 3)

	// Get prototype status technologies
	techs, total, err = repo.GetByStatus(ctx, "prototype", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total) // KAAN, Altay
	assert.Len(t, techs, 2)
}

func TestSimpleTechnologyRepository_GetByYearRange(t *testing.T) {
	db := setupSimpleTestDB(t)
	repo := NewTechnologyRepository(db)
	ctx := context.Background()

	seedSimpleTestData(t, db)

	// Technologies from 2020-2025
	techs, total, err := repo.GetByYearRange(ctx, 2020, 2025, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total) // KAAN (2023), Altay (2024)
	assert.Len(t, techs, 2)

	// Technologies from 2010-2015
	techs, total, err = repo.GetByYearRange(ctx, 2010, 2015, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total) // TB2, F-35, Su-57
	assert.Len(t, techs, 3)
}
