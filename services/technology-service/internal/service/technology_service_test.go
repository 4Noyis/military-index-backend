package service

import (
	"context"
	"errors"
	"testing"

	"github.com/4Noyis/military-index-backend/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTechnologyRepository is a mock implementation of TechnologyRepository
type MockTechnologyRepository struct {
	mock.Mock
}

func (m *MockTechnologyRepository) GetAll(ctx context.Context, page, limit int, sortBy, sortOrder string) ([]models.Technology, int64, error) {
	args := m.Called(ctx, page, limit, sortBy, sortOrder)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyRepository) GetByID(ctx context.Context, id uint) (*models.Technology, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Technology), args.Error(1)
}

func (m *MockTechnologyRepository) GetByCountry(ctx context.Context, code string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, code, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyRepository) GetByCategory(ctx context.Context, categoryName string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, categoryName, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyRepository) GetByStatus(ctx context.Context, status string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, status, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyRepository) GetByYearRange(ctx context.Context, startYear, endYear int, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, startYear, endYear, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyRepository) Search(ctx context.Context, query string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, query, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyRepository) Create(ctx context.Context, technology *models.Technology) error {
	args := m.Called(ctx, technology)
	return args.Error(0)
}

func (m *MockTechnologyRepository) Update(ctx context.Context, technology *models.Technology) error {
	args := m.Called(ctx, technology)
	return args.Error(0)
}

func (m *MockTechnologyRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTechnologyService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		sortBy        string
		sortOrder     string
		mockReturn    []models.Technology
		mockTotal     int64
		mockError     error
		expectedPage  int
		expectedLimit int
		wantErr       bool
	}{
		{
			name:      "Valid request with defaults",
			page:      0,
			limit:     0,
			sortBy:    "name",
			sortOrder: "asc",
			mockReturn: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
			},
			mockTotal:     1,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10,
			wantErr:       false,
		},
		{
			name:      "Valid request with custom pagination",
			page:      2,
			limit:     20,
			sortBy:    "year_developed",
			sortOrder: "desc",
			mockReturn: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
			},
			mockTotal:     50,
			mockError:     nil,
			expectedPage:  2,
			expectedLimit: 20,
			wantErr:       false,
		},
		{
			name:          "Invalid limit (too high)",
			page:          1,
			limit:         200,
			sortBy:        "name",
			sortOrder:     "asc",
			mockReturn:    []models.Technology{},
			mockTotal:     0,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10, // Should be normalized
			wantErr:       false,
		},
		{
			name:          "Repository error",
			page:          1,
			limit:         10,
			sortBy:        "name",
			sortOrder:     "asc",
			mockReturn:    []models.Technology{},
			mockTotal:     0,
			mockError:     errors.New("database error"),
			expectedPage:  1,
			expectedLimit: 10,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)
			mockRepo.On("GetAll", mock.Anything, tt.expectedPage, tt.expectedLimit, tt.sortBy, tt.sortOrder).
				Return(tt.mockReturn, tt.mockTotal, tt.mockError)

			service := NewTechnologyService(mockRepo)
			technologies, total, err := service.GetAll(context.Background(), tt.page, tt.limit, tt.sortBy, tt.sortOrder)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTotal, total)
				assert.Equal(t, len(tt.mockReturn), len(technologies))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTechnologyService_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		id         uint
		mockReturn *models.Technology
		mockError  error
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Valid ID",
			id:         1,
			mockReturn: &models.Technology{ID: 1, Name: "KAAN Fighter"},
			mockError:  nil,
			wantErr:    false,
		},
		{
			name:       "Zero ID",
			id:         0,
			mockReturn: nil,
			mockError:  nil,
			wantErr:    true,
			errMsg:     "invalid technology ID",
		},
		{
			name:       "Technology not found",
			id:         999,
			mockReturn: nil,
			mockError:  errors.New("technology not found"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Only setup mock expectation if ID is valid
			if tt.id != 0 {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockReturn, tt.mockError)
			}

			service := NewTechnologyService(mockRepo)
			technology, err := service.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, technology)
				assert.Equal(t, tt.mockReturn.Name, technology.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTechnologyService_GetByCountry(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		page          int
		limit         int
		mockReturn    []models.Technology
		mockTotal     int64
		mockError     error
		expectedPage  int
		expectedLimit int
		wantErr       bool
		errMsg        string
	}{
		{
			name:  "Valid country code",
			code:  "TUR",
			page:  1,
			limit: 10,
			mockReturn: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
				{ID: 2, Name: "Bayraktar TB2"},
			},
			mockTotal:     2,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10,
			wantErr:       false,
		},
		{
			name:          "Empty country code",
			code:          "",
			page:          1,
			limit:         10,
			mockReturn:    nil,
			mockTotal:     0,
			mockError:     nil,
			expectedPage:  0,
			expectedLimit: 0,
			wantErr:       true,
			errMsg:        "country code cannot be empty",
		},
		{
			name:          "Invalid pagination",
			code:          "USA",
			page:          0,
			limit:         0,
			mockReturn:    []models.Technology{},
			mockTotal:     0,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Only setup mock if code is valid
			if tt.code != "" {
				mockRepo.On("GetByCountry", mock.Anything, tt.code, tt.expectedPage, tt.expectedLimit).
					Return(tt.mockReturn, tt.mockTotal, tt.mockError)
			}

			service := NewTechnologyService(mockRepo)
			technologies, total, err := service.GetByCountry(context.Background(), tt.code, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTotal, total)
				assert.Equal(t, len(tt.mockReturn), len(technologies))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTechnologyService_GetByCategory(t *testing.T) {
	tests := []struct {
		name          string
		categoryName  string
		page          int
		limit         int
		mockReturn    []models.Technology
		mockTotal     int64
		mockError     error
		expectedPage  int
		expectedLimit int
		wantErr       bool
		errMsg        string
	}{
		{
			name:         "Valid category",
			categoryName: "Aircraft",
			page:         1,
			limit:        10,
			mockReturn: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
			},
			mockTotal:     1,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10,
			wantErr:       false,
		},
		{
			name:          "Empty category name",
			categoryName:  "",
			page:          1,
			limit:         10,
			mockReturn:    nil,
			mockTotal:     0,
			mockError:     nil,
			expectedPage:  0,
			expectedLimit: 0,
			wantErr:       true,
			errMsg:        "category name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			if tt.categoryName != "" {
				mockRepo.On("GetByCategory", mock.Anything, tt.categoryName, tt.expectedPage, tt.expectedLimit).
					Return(tt.mockReturn, tt.mockTotal, tt.mockError)
			}

			service := NewTechnologyService(mockRepo)
			technologies, total, err := service.GetByCategory(context.Background(), tt.categoryName, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTotal, total)
				assert.Equal(t, len(tt.mockReturn), len(technologies))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTechnologyService_GetByStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		page       int
		limit      int
		mockReturn []models.Technology
		mockTotal  int64
		mockError  error
		wantErr    bool
		errMsg     string
	}{
		{
			name:   "Valid status - current",
			status: "current",
			page:   1,
			limit:  10,
			mockReturn: []models.Technology{
				{ID: 1, Name: "F-35", Status: "current"},
			},
			mockTotal: 1,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:   "Valid status - prototype",
			status: "prototype",
			page:   1,
			limit:  10,
			mockReturn: []models.Technology{
				{ID: 1, Name: "KAAN Fighter", Status: "prototype"},
			},
			mockTotal: 1,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:       "Invalid status",
			status:     "invalid",
			page:       1,
			limit:      10,
			mockReturn: nil,
			mockTotal:  0,
			mockError:  nil,
			wantErr:    true,
			errMsg:     "invalid status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Only setup mock for valid statuses
			validStatuses := map[string]bool{
				"historical": true, "current": true, "future": true,
				"concept": true, "prototype": true,
			}

			if validStatuses[tt.status] {
				mockRepo.On("GetByStatus", mock.Anything, tt.status, 1, 10).
					Return(tt.mockReturn, tt.mockTotal, tt.mockError)
			}

			service := NewTechnologyService(mockRepo)
			technologies, total, err := service.GetByStatus(context.Background(), tt.status, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTotal, total)
				assert.Equal(t, len(tt.mockReturn), len(technologies))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTechnologyService_GetByYearRange(t *testing.T) {
	tests := []struct {
		name       string
		startYear  int
		endYear    int
		page       int
		limit      int
		mockReturn []models.Technology
		mockTotal  int64
		mockError  error
		wantErr    bool
		errMsg     string
	}{
		{
			name:      "Valid year range",
			startYear: 2010,
			endYear:   2020,
			page:      1,
			limit:     10,
			mockReturn: []models.Technology{
				{ID: 1, Name: "F-35"},
			},
			mockTotal: 1,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:       "Start year greater than end year",
			startYear:  2020,
			endYear:    2010,
			page:       1,
			limit:      10,
			mockReturn: nil,
			mockTotal:  0,
			mockError:  nil,
			wantErr:    true,
			errMsg:     "start year cannot be greater than end year",
		},
		{
			name:       "Year out of range (too early)",
			startYear:  1800,
			endYear:    1900,
			page:       1,
			limit:      10,
			mockReturn: nil,
			mockTotal:  0,
			mockError:  nil,
			wantErr:    true,
			errMsg:     "year range must be between 1900 and 2100",
		},
		{
			name:       "Year out of range (too late)",
			startYear:  2100,
			endYear:    2200,
			page:       1,
			limit:      10,
			mockReturn: nil,
			mockTotal:  0,
			mockError:  nil,
			wantErr:    true,
			errMsg:     "year range must be between 1900 and 2100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Only setup mock for valid ranges
			if tt.startYear <= tt.endYear && tt.startYear >= 1900 && tt.endYear <= 2100 {
				mockRepo.On("GetByYearRange", mock.Anything, tt.startYear, tt.endYear, 1, 10).
					Return(tt.mockReturn, tt.mockTotal, tt.mockError)
			}

			service := NewTechnologyService(mockRepo)
			technologies, total, err := service.GetByYearRange(context.Background(), tt.startYear, tt.endYear, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTotal, total)
				assert.Equal(t, len(tt.mockReturn), len(technologies))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTechnologyService_Search(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		page       int
		limit      int
		mockReturn []models.Technology
		mockTotal  int64
		mockError  error
		wantErr    bool
		errMsg     string
	}{
		{
			name:  "Valid search query",
			query: "Fighter",
			page:  1,
			limit: 10,
			mockReturn: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
			},
			mockTotal: 1,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:       "Empty query",
			query:      "",
			page:       1,
			limit:      10,
			mockReturn: nil,
			mockTotal:  0,
			mockError:  nil,
			wantErr:    true,
			errMsg:     "search query cannot be empty",
		},
		{
			name:       "Query too short",
			query:      "F",
			page:       1,
			limit:      10,
			mockReturn: nil,
			mockTotal:  0,
			mockError:  nil,
			wantErr:    true,
			errMsg:     "search query must be at least 2 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Only setup mock for valid queries
			if len(tt.query) >= 2 {
				mockRepo.On("Search", mock.Anything, tt.query, 1, 10).
					Return(tt.mockReturn, tt.mockTotal, tt.mockError)
			}

			service := NewTechnologyService(mockRepo)
			technologies, total, err := service.Search(context.Background(), tt.query, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockTotal, total)
				assert.Equal(t, len(tt.mockReturn), len(technologies))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTechnologyService_Create(t *testing.T) {
	year2023 := 2023
	year2025 := 2025
	year1800 := 1800
	year2200 := 2200

	tests := []struct {
		name       string
		technology *models.Technology
		mockError  error
		wantErr    bool
		errMsg     string
	}{
		{
			name: "Valid technology",
			technology: &models.Technology{
				Name:          "Test Aircraft",
				CountryID:     1,
				CategoryID:    1,
				YearDeveloped: &year2023,
				Status:        "concept",
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name: "Missing name",
			technology: &models.Technology{
				CountryID:  1,
				CategoryID: 1,
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "technology name is required",
		},
		{
			name: "Missing country ID",
			technology: &models.Technology{
				Name:       "Test",
				CategoryID: 1,
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "country ID is required",
		},
		{
			name: "Missing category ID",
			technology: &models.Technology{
				Name:      "Test",
				CountryID: 1,
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "category ID is required",
		},
		{
			name: "Invalid status",
			technology: &models.Technology{
				Name:       "Test",
				CountryID:  1,
				CategoryID: 1,
				Status:     "invalid",
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "invalid status",
		},
		{
			name: "Year developed out of range (too early)",
			technology: &models.Technology{
				Name:          "Test",
				CountryID:     1,
				CategoryID:    1,
				YearDeveloped: &year1800,
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "year developed must be between 1900 and 2100",
		},
		{
			name: "Year developed out of range (too late)",
			technology: &models.Technology{
				Name:          "Test",
				CountryID:     1,
				CategoryID:    1,
				YearDeveloped: &year2200,
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "year developed must be between 1900 and 2100",
		},
		{
			name: "Year developed greater than year deployed",
			technology: &models.Technology{
				Name:          "Test",
				CountryID:     1,
				CategoryID:    1,
				YearDeveloped: &year2025,
				YearDeployed:  &year2023,
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "year developed cannot be greater than year deployed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Only setup mock for valid requests
			if tt.name == "Valid technology" {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Technology")).
					Return(tt.mockError)
			}

			service := NewTechnologyService(mockRepo)
			err := service.Create(context.Background(), tt.technology)

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

func TestTechnologyService_Update(t *testing.T) {
	year2023 := 2023

	tests := []struct {
		name           string
		id             uint
		technology     *models.Technology
		mockGetReturn  *models.Technology
		mockGetError   error
		mockUpdateError error
		wantErr        bool
		errMsg         string
	}{
		{
			name: "Valid update",
			id:   1,
			technology: &models.Technology{
				Name:          "Updated KAAN",
				CountryID:     1,
				CategoryID:    1,
				YearDeveloped: &year2023,
			},
			mockGetReturn: &models.Technology{
				ID:         1,
				Name:       "KAAN Fighter",
				CountryID:  1,
				CategoryID: 1,
			},
			mockGetError:    nil,
			mockUpdateError: nil,
			wantErr:         false,
		},
		{
			name: "Zero ID",
			id:   0,
			technology: &models.Technology{
				Name:       "Test",
				CountryID:  1,
				CategoryID: 1,
			},
			mockGetReturn:   nil,
			mockGetError:    nil,
			mockUpdateError: nil,
			wantErr:         true,
			errMsg:          "invalid technology ID",
		},
		{
			name: "Technology not found",
			id:   999,
			technology: &models.Technology{
				Name:       "Test",
				CountryID:  1,
				CategoryID: 1,
			},
			mockGetReturn:   nil,
			mockGetError:    errors.New("technology not found"),
			mockUpdateError: nil,
			wantErr:         true,
			errMsg:          "technology not found",
		},
		{
			name: "Missing name",
			id:   1,
			technology: &models.Technology{
				CountryID:  1,
				CategoryID: 1,
			},
			mockGetReturn: &models.Technology{ID: 1},
			mockGetError:  nil,
			wantErr:       true,
			errMsg:        "technology name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Setup GetByID mock if ID is valid
			if tt.id != 0 {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockGetReturn, tt.mockGetError)
			}

			// Setup Update mock if get succeeds and validation passes
			if tt.mockGetReturn != nil && tt.technology.Name != "" && tt.technology.CountryID != 0 && tt.technology.CategoryID != 0 {
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Technology")).
					Return(tt.mockUpdateError)
			}

			service := NewTechnologyService(mockRepo)
			err := service.Update(context.Background(), tt.id, tt.technology)

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

func TestTechnologyService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockGetReturn *models.Technology
		mockGetError  error
		mockDelError  error
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "Valid delete",
			id:            1,
			mockGetReturn: &models.Technology{ID: 1, Name: "KAAN Fighter"},
			mockGetError:  nil,
			mockDelError:  nil,
			wantErr:       false,
		},
		{
			name:          "Zero ID",
			id:            0,
			mockGetReturn: nil,
			mockGetError:  nil,
			mockDelError:  nil,
			wantErr:       true,
			errMsg:        "invalid technology ID",
		},
		{
			name:          "Technology not found",
			id:            999,
			mockGetReturn: nil,
			mockGetError:  errors.New("technology not found"),
			mockDelError:  nil,
			wantErr:       true,
			errMsg:        "technology not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTechnologyRepository)

			// Setup GetByID mock if ID is valid
			if tt.id != 0 {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockGetReturn, tt.mockGetError)
			}

			// Setup Delete mock if get succeeds
			if tt.mockGetReturn != nil {
				mockRepo.On("Delete", mock.Anything, tt.id).
					Return(tt.mockDelError)
			}

			service := NewTechnologyService(mockRepo)
			err := service.Delete(context.Background(), tt.id)

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
