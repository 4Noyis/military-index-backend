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

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Country{})
	require.NoError(t, err, "Failed to migrate test database")

	return db
}

// seedCountries adds sample countries to the database
func seedCountries(t *testing.T, db *gorm.DB) []models.Country {
	countries := []models.Country{
		{Name: "Turkey", Code: "TUR", FlagURL: "/flags/turkey.svg"},
		{Name: "United States", Code: "USA", FlagURL: "/flags/usa.svg"},
		{Name: "Russia", Code: "RUS", FlagURL: "/flags/russia.svg"},
		{Name: "China", Code: "CHN", FlagURL: "/flags/china.svg"},
		{Name: "Germany", Code: "DEU", FlagURL: "/flags/germany.svg"},
	}

	for i := range countries {
		err := db.Create(&countries[i]).Error
		require.NoError(t, err, "Failed to seed country")
	}

	return countries
}

func TestCountryRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		country *models.Country
		wantErr bool
	}{
		{
			name: "Valid country",
			country: &models.Country{
				Name:    "Japan",
				Code:    "JPN",
				FlagURL: "/flags/japan.svg",
			},
			wantErr: false,
		},
		{
			name: "Country with minimal data",
			country: &models.Country{
				Name: "France",
				Code: "FRA",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.country)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.country.ID, "Country ID should be set after creation")
			}
		})
	}
}

func TestCountryRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	countries := seedCountries(t, db)

	tests := []struct {
		name    string
		id      uint
		want    string // expected country name
		wantErr bool
	}{
		{
			name:    "Existing country",
			id:      countries[0].ID,
			want:    "Turkey",
			wantErr: false,
		},
		{
			name:    "Non-existent country",
			id:      9999,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, err := repo.GetByID(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, country)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, country)
				assert.Equal(t, tt.want, country.Name)
			}
		})
	}
}

func TestCountryRepository_GetByCode(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	seedCountries(t, db)

	tests := []struct {
		name    string
		code    string
		want    string // expected country name
		wantErr bool
	}{
		{
			name:    "Existing country code",
			code:    "TUR",
			want:    "Turkey",
			wantErr: false,
		},
		{
			name:    "Another existing country",
			code:    "USA",
			want:    "United States",
			wantErr: false,
		},
		{
			name:    "Non-existent country code",
			code:    "XXX",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, err := repo.GetByCode(ctx, tt.code)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, country)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, country)
				assert.Equal(t, tt.want, country.Name)
				assert.Equal(t, tt.code, country.Code)
			}
		})
	}
}

func TestCountryRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	seedCountries(t, db) // 5 countries

	tests := []struct {
		name      string
		page      int
		limit     int
		wantCount int
		wantTotal int64
	}{
		{
			name:      "First page with limit 3",
			page:      1,
			limit:     3,
			wantCount: 3,
			wantTotal: 5,
		},
		{
			name:      "Second page with limit 3",
			page:      2,
			limit:     3,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "All countries",
			page:      1,
			limit:     10,
			wantCount: 5,
			wantTotal: 5,
		},
		{
			name:      "Empty page",
			page:      10,
			limit:     10,
			wantCount: 0,
			wantTotal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countries, total, err := repo.GetAll(ctx, tt.page, tt.limit)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(countries))
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}

func TestCountryRepository_Search(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	seedCountries(t, db)

	tests := []struct {
		name      string
		query     string
		page      int
		limit     int
		wantCount int
		wantTotal int64
	}{
		{
			name:      "Search for 'United'",
			query:     "United",
			page:      1,
			limit:     10,
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "Search for 'a' (matches multiple)",
			query:     "a",
			page:      1,
			limit:     10,
			wantCount: 4, // Turkey, United States, Russia, China, Germany (4 have 'a')
			wantTotal: 4,
		},
		{
			name:      "Search with no results",
			query:     "xyz",
			page:      1,
			limit:     10,
			wantCount: 0,
			wantTotal: 0,
		},
		{
			name:      "Case insensitive search",
			query:     "turkey",
			page:      1,
			limit:     10,
			wantCount: 1,
			wantTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countries, total, err := repo.Search(ctx, tt.query, tt.page, tt.limit)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(countries))
			assert.Equal(t, tt.wantTotal, total)
		})
	}
}

func TestCountryRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	countries := seedCountries(t, db)

	t.Run("Update existing country", func(t *testing.T) {
		// Get the first country
		country := countries[0]
		country.Name = "Republic of Turkey"
		country.FlagURL = "/flags/turkey-new.svg"

		err := repo.Update(ctx, &country)
		assert.NoError(t, err)

		// Verify the update
		updated, err := repo.GetByID(ctx, country.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Republic of Turkey", updated.Name)
		assert.Equal(t, "/flags/turkey-new.svg", updated.FlagURL)
	})
}

func TestCountryRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	countries := seedCountries(t, db)

	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "Delete existing country",
			id:      countries[0].ID,
			wantErr: false,
		},
		{
			name:    "Delete non-existent country",
			id:      9999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify deletion
				_, err := repo.GetByID(ctx, tt.id)
				assert.Error(t, err, "Country should not exist after deletion")
			}
		})
	}
}

func TestCountryRepository_PaginationConsistency(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCountryRepository(db)
	ctx := context.Background()

	seedCountries(t, db) // 5 countries

	// Get all countries in one call
	allCountries, total1, err := repo.GetAll(ctx, 1, 100)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total1)
	assert.Len(t, allCountries, 5)

	// Get countries in pages and verify they match
	page1, total2, err := repo.GetAll(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, total1, total2)
	assert.Len(t, page1, 2)

	page2, total3, err := repo.GetAll(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Equal(t, total1, total3)
	assert.Len(t, page2, 2)

	page3, total4, err := repo.GetAll(ctx, 3, 2)
	assert.NoError(t, err)
	assert.Equal(t, total1, total4)
	assert.Len(t, page3, 1)

	// Verify no duplicates across pages
	combined := append(page1, page2...)
	combined = append(combined, page3...)

	assert.Len(t, combined, 5)

	// Check all IDs are unique
	seen := make(map[uint]bool)
	for _, c := range combined {
		assert.False(t, seen[c.ID], "Duplicate ID found in pagination")
		seen[c.ID] = true
	}
}
