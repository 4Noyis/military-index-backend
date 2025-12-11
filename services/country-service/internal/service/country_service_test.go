package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/4Noyis/military-index-backend/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCountryRepository is a mock implementation of CountryRepository
type MockCountryRepository struct {
	mock.Mock
}

func (m *MockCountryRepository) GetAll(ctx context.Context, page, limit int) ([]models.Country, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]models.Country), args.Get(1).(int64), args.Error(2)
}

func (m *MockCountryRepository) GetByID(ctx context.Context, id uint) (*models.Country, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryRepository) GetByCode(ctx context.Context, code string) (*models.Country, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryRepository) Search(ctx context.Context, query string, page, limit int) ([]models.Country, int64, error) {
	args := m.Called(ctx, query, page, limit)
	return args.Get(0).([]models.Country), args.Get(1).(int64), args.Error(2)
}

func (m *MockCountryRepository) Create(ctx context.Context, country *models.Country) error {
	args := m.Called(ctx, country)
	return args.Error(0)
}

func (m *MockCountryRepository) Update(ctx context.Context, country *models.Country) error {
	args := m.Called(ctx, country)
	return args.Error(0)
}

func (m *MockCountryRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCountryService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		mockCountries []models.Country
		mockTotal     int64
		mockError     error
		expectPage    int
		expectLimit   int
		wantErr       bool
	}{
		{
			name:  "Valid pagination",
			page:  1,
			limit: 10,
			mockCountries: []models.Country{
				{ID: 1, Name: "Turkey", Code: "TUR"},
				{ID: 2, Name: "USA", Code: "USA"},
			},
			mockTotal:   2,
			mockError:   nil,
			expectPage:  1,
			expectLimit: 10,
			wantErr:     false,
		},
		{
			name:  "Invalid page number (< 1)",
			page:  0,
			limit: 10,
			mockCountries: []models.Country{
				{ID: 1, Name: "Turkey", Code: "TUR"},
			},
			mockTotal:   1,
			mockError:   nil,
			expectPage:  1, // Should default to 1
			expectLimit: 10,
			wantErr:     false,
		},
		{
			name:  "Invalid limit (> 100)",
			page:  1,
			limit: 200,
			mockCountries: []models.Country{
				{ID: 1, Name: "Turkey", Code: "TUR"},
			},
			mockTotal:   1,
			mockError:   nil,
			expectPage:  1,
			expectLimit: 10, // Should default to 10
			wantErr:     false,
		},
		{
			name:          "Repository error",
			page:          1,
			limit:         10,
			mockCountries: []models.Country{},
			mockTotal:     0,
			mockError:     errors.New("database error"),
			expectPage:    1,
			expectLimit:   10,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCountryRepository)
			mockRepo.On("GetAll", mock.Anything, tt.expectPage, tt.expectLimit).
				Return(tt.mockCountries, tt.mockTotal, tt.mockError)

			service := NewCountryService(mockRepo)
			ctx := context.Background()

			countries, total, err := service.GetAll(ctx, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockCountries, countries)
				assert.Equal(t, tt.mockTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCountryService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		id          uint
		mockCountry *models.Country
		mockError   error
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Valid ID",
			id:          1,
			mockCountry: &models.Country{ID: 1, Name: "Turkey", Code: "TUR"},
			mockError:   nil,
			wantErr:     false,
		},
		{
			name:        "Invalid ID (0)",
			id:          0,
			mockCountry: nil,
			mockError:   nil,
			wantErr:     true,
			errMsg:      "invalid country ID",
		},
		{
			name:        "Country not found",
			id:          999,
			mockCountry: nil,
			mockError:   errors.New("country not found"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCountryRepository)

			if tt.id != 0 {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockCountry, tt.mockError)
			}

			service := NewCountryService(mockRepo)
			ctx := context.Background()

			country, err := service.GetByID(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockCountry, country)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCountryService_GetByCode(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectCode  string
		mockCountry *models.Country
		mockError   error
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "Valid code",
			code:        "TUR",
			expectCode:  "TUR",
			mockCountry: &models.Country{ID: 1, Name: "Turkey", Code: "TUR"},
			mockError:   nil,
			wantErr:     false,
		},
		{
			name:        "Lowercase code (should normalize)",
			code:        "tur",
			expectCode:  "TUR",
			mockCountry: &models.Country{ID: 1, Name: "Turkey", Code: "TUR"},
			mockError:   nil,
			wantErr:     false,
		},
		{
			name:        "Code with whitespace (should trim)",
			code:        " TUR ",
			expectCode:  "TUR",
			mockCountry: &models.Country{ID: 1, Name: "Turkey", Code: "TUR"},
			mockError:   nil,
			wantErr:     false,
		},
		{
			name:        "Empty code",
			code:        "",
			expectCode:  "",
			mockCountry: nil,
			mockError:   nil,
			wantErr:     true,
			errMsg:      "country code is required",
		},
		{
			name:        "Country not found",
			code:        "XXX",
			expectCode:  "XXX",
			mockCountry: nil,
			mockError:   errors.New("country not found"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCountryRepository)

			if tt.expectCode != "" {
				mockRepo.On("GetByCode", mock.Anything, tt.expectCode).
					Return(tt.mockCountry, tt.mockError)
			}

			service := NewCountryService(mockRepo)
			ctx := context.Background()

			country, err := service.GetByCode(ctx, tt.code)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockCountry, country)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCountryService_Search(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		page          int
		limit         int
		mockCountries []models.Country
		mockTotal     int64
		mockError     error
		expectPage    int
		expectLimit   int
		wantErr       bool
		errMsg        string
	}{
		{
			name:  "Valid search",
			query: "Turkey",
			page:  1,
			limit: 10,
			mockCountries: []models.Country{
				{ID: 1, Name: "Turkey", Code: "TUR"},
			},
			mockTotal:   1,
			mockError:   nil,
			expectPage:  1,
			expectLimit: 10,
			wantErr:     false,
		},
		{
			name:          "Empty query",
			query:         "",
			page:          1,
			limit:         10,
			mockCountries: []models.Country{},
			mockTotal:     0,
			mockError:     nil,
			expectPage:    1,
			expectLimit:   10,
			wantErr:       true,
			errMsg:        "search query is required",
		},
		{
			name:  "Invalid pagination (corrected)",
			query: "test",
			page:  -1,
			limit: 200,
			mockCountries: []models.Country{
				{ID: 1, Name: "Test", Code: "TST"},
			},
			mockTotal:   1,
			mockError:   nil,
			expectPage:  1,  // Should be corrected
			expectLimit: 10, // Should be corrected
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCountryRepository)

			if tt.query != "" {
				mockRepo.On("Search", mock.Anything, tt.query, tt.expectPage, tt.expectLimit).
					Return(tt.mockCountries, tt.mockTotal, tt.mockError)
			}

			service := NewCountryService(mockRepo)
			ctx := context.Background()

			countries, total, err := service.Search(ctx, tt.query, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockCountries, countries)
				assert.Equal(t, tt.mockTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCountryService_Create(t *testing.T) {
	tests := []struct {
		name      string
		country   *models.Country
		mockError error
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Valid country",
			country: &models.Country{
				Name:    "Japan",
				Code:    "JPN",
				FlagURL: "/flags/japan.svg",
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name: "Country with lowercase code (should normalize)",
			country: &models.Country{
				Name: "France",
				Code: "fra",
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "Nil country",
			country:   nil,
			mockError: nil,
			wantErr:   true,
			errMsg:    "country data is required",
		},
		{
			name: "Missing name",
			country: &models.Country{
				Code: "TST",
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "country name is required",
		},
		{
			name: "Missing code",
			country: &models.Country{
				Name: "Test Country",
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "country code is required",
		},
		{
			name: "Invalid code length (too short)",
			country: &models.Country{
				Name: "Test",
				Code: "T",
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "country code must be 2-3 characters",
		},
		{
			name: "Invalid code length (too long)",
			country: &models.Country{
				Name: "Test",
				Code: "TEST",
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "country code must be 2-3 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCountryRepository)

			if !tt.wantErr && tt.country != nil {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Country")).
					Return(tt.mockError)
			}

			service := NewCountryService(mockRepo)
			ctx := context.Background()

			err := service.Create(ctx, tt.country)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				// Verify normalization
				if tt.country != nil {
					assert.Equal(t, tt.country.Code, strings.ToUpper(tt.country.Code))
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCountryService_Update(t *testing.T) {
	existingCountry := &models.Country{
		ID:      1,
		Name:    "Turkey",
		Code:    "TUR",
		FlagURL: "/flags/turkey.svg",
	}

	tests := []struct {
		name           string
		id             uint
		country        *models.Country
		mockGetCountry *models.Country
		mockGetError   error
		mockUpdError   error
		wantErr        bool
		errMsg         string
	}{
		{
			name: "Valid update",
			id:   1,
			country: &models.Country{
				Name:    "Republic of Turkey",
				Code:    "TUR",
				FlagURL: "/flags/turkey-new.svg",
			},
			mockGetCountry: existingCountry,
			mockGetError:   nil,
			mockUpdError:   nil,
			wantErr:        false,
		},
		{
			name: "Invalid ID (0)",
			id:   0,
			country: &models.Country{
				Name: "Test",
				Code: "TST",
			},
			mockGetCountry: nil,
			mockGetError:   nil,
			mockUpdError:   nil,
			wantErr:        true,
			errMsg:         "invalid country ID",
		},
		{
			name: "Country not found",
			id:   999,
			country: &models.Country{
				Name: "Test",
				Code: "TST",
			},
			mockGetCountry: nil,
			mockGetError:   errors.New("country not found"),
			mockUpdError:   nil,
			wantErr:        true,
		},
		{
			name: "Invalid update data",
			id:   1,
			country: &models.Country{
				Name: "", // Empty name
				Code: "TUR",
			},
			mockGetCountry: existingCountry,
			mockGetError:   nil,
			mockUpdError:   nil,
			wantErr:        true,
			errMsg:         "country name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCountryRepository)

			if tt.id != 0 {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockGetCountry, tt.mockGetError)

				if tt.mockGetError == nil && tt.mockGetCountry != nil && tt.country.Name != "" {
					mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Country")).
						Return(tt.mockUpdError)
				}
			}

			service := NewCountryService(mockRepo)
			ctx := context.Background()

			err := service.Update(ctx, tt.id, tt.country)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCountryService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockError error
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Valid delete",
			id:        1,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "Invalid ID (0)",
			id:        0,
			mockError: nil,
			wantErr:   true,
			errMsg:    "invalid country ID",
		},
		{
			name:      "Country not found",
			id:        999,
			mockError: errors.New("country not found"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCountryRepository)

			if tt.id != 0 {
				mockRepo.On("Delete", mock.Anything, tt.id).
					Return(tt.mockError)
			}

			service := NewCountryService(mockRepo)
			ctx := context.Background()

			err := service.Delete(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
