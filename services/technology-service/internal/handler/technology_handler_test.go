package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4Noyis/military-index-backend/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTechnologyService is a mock implementation of TechnologyService
type MockTechnologyService struct {
	mock.Mock
}

func (m *MockTechnologyService) GetAll(ctx context.Context, page, limit int, sortBy, sortOrder string) ([]models.Technology, int64, error) {
	args := m.Called(ctx, page, limit, sortBy, sortOrder)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyService) GetByID(ctx context.Context, id uint) (*models.Technology, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Technology), args.Error(1)
}

func (m *MockTechnologyService) GetByCountry(ctx context.Context, code string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, code, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyService) GetByCategory(ctx context.Context, categoryName string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, categoryName, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyService) GetByStatus(ctx context.Context, status string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, status, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyService) GetByYearRange(ctx context.Context, startYear, endYear int, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, startYear, endYear, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyService) Search(ctx context.Context, query string, page, limit int) ([]models.Technology, int64, error) {
	args := m.Called(ctx, query, page, limit)
	return args.Get(0).([]models.Technology), args.Get(1).(int64), args.Error(2)
}

func (m *MockTechnologyService) Create(ctx context.Context, technology *models.Technology) error {
	args := m.Called(ctx, technology)
	return args.Error(0)
}

func (m *MockTechnologyService) Update(ctx context.Context, id uint, technology *models.Technology) error {
	args := m.Called(ctx, id, technology)
	return args.Error(0)
}

func (m *MockTechnologyService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to parse JSON response
func parseJSONResponse(t *testing.T, body []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(body, &result)
	assert.NoError(t, err, "Failed to parse JSON response")
	return result
}

func TestTechnologyHandler_GetAll(t *testing.T) {
	tests := []struct {
		name              string
		queryParams       string
		mockTechnologies  []models.Technology
		mockTotal         int64
		mockError         error
		expectedStatus    int
		expectedPage      int
		expectedLimit     int
		expectedSortBy    string
		expectedSortOrder string
		checkResponse     func(*testing.T, map[string]interface{})
	}{
		{
			name:        "Success - Get all technologies",
			queryParams: "",
			mockTechnologies: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
				{ID: 2, Name: "Bayraktar TB2"},
			},
			mockTotal:         2,
			mockError:         nil,
			expectedStatus:    http.StatusOK,
			expectedPage:      0,
			expectedLimit:     0,
			expectedSortBy:    "",
			expectedSortOrder: "",
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)

				pagination := resp["pagination"].(map[string]interface{})
				assert.Equal(t, float64(1), pagination["page"])
				assert.Equal(t, float64(10), pagination["limit"])
				assert.Equal(t, float64(2), pagination["total"])
			},
		},
		{
			name:        "Success - With pagination and sorting",
			queryParams: "?page=2&limit=5&sort=year_developed&order=desc",
			mockTechnologies: []models.Technology{
				{ID: 6, Name: "Su-57"},
			},
			mockTotal:         10,
			mockError:         nil,
			expectedStatus:    http.StatusOK,
			expectedPage:      2,
			expectedLimit:     5,
			expectedSortBy:    "year_developed",
			expectedSortOrder: "desc",
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				pagination := resp["pagination"].(map[string]interface{})
				assert.Equal(t, float64(2), pagination["page"])
				assert.Equal(t, float64(5), pagination["limit"])
				assert.Equal(t, float64(10), pagination["total"])
			},
		},
		{
			name:              "Success - Search query",
			queryParams:       "?q=Fighter",
			mockTechnologies:  []models.Technology{{ID: 1, Name: "KAAN Fighter"}},
			mockTotal:         1,
			mockError:         nil,
			expectedStatus:    http.StatusOK,
			expectedPage:      0,
			expectedLimit:     0,
			expectedSortBy:    "",
			expectedSortOrder: "",
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 1)
			},
		},
		{
			name:              "Error - Service error",
			queryParams:       "",
			mockTechnologies:  []models.Technology{},
			mockTotal:         0,
			mockError:         errors.New("database error"),
			expectedStatus:    http.StatusInternalServerError,
			expectedPage:      0,
			expectedLimit:     0,
			expectedSortBy:    "",
			expectedSortOrder: "",
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "failed to fetch technologies")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			// Setup mock expectations
			if tt.queryParams != "" && len(tt.queryParams) > 3 && tt.queryParams[1:3] == "q=" {
				// Search request
				mockService.On("Search", mock.Anything, "Fighter", tt.expectedPage, tt.expectedLimit).
					Return(tt.mockTechnologies, tt.mockTotal, tt.mockError)
			} else {
				// GetAll request
				mockService.On("GetAll", mock.Anything, tt.expectedPage, tt.expectedLimit, tt.expectedSortBy, tt.expectedSortOrder).
					Return(tt.mockTechnologies, tt.mockTotal, tt.mockError)
			}

			// Create request
			req := httptest.NewRequest("GET", "/api/v1/technologies"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			// Call handler
			handler.GetAll(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Parse and check response
			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		urlID          string
		mockTechnology *models.Technology
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:  "Success - Get technology by ID",
			urlID: "1",
			mockTechnology: &models.Technology{
				ID:   1,
				Name: "KAAN Fighter",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "KAAN Fighter", data["name"])
			},
		},
		{
			name:           "Error - Invalid ID format",
			urlID:          "abc",
			mockTechnology: nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "invalid technology ID")
			},
		},
		{
			name:           "Error - Technology not found",
			urlID:          "999",
			mockTechnology: nil,
			mockError:      errors.New("technology not found"),
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "technology not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			// Setup mock expectations (only if ID is valid)
			if tt.urlID != "abc" {
				id := uint(1)
				if tt.urlID == "999" {
					id = 999
				}
				mockService.On("GetByID", mock.Anything, id).
					Return(tt.mockTechnology, tt.mockError)
			}

			// Create request
			req := httptest.NewRequest("GET", "/api/v1/technologies/"+tt.urlID, nil)
			rr := httptest.NewRecorder()

			// Call handler
			handler.GetByID(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Parse and check response
			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_GetByCountry(t *testing.T) {
	tests := []struct {
		name             string
		urlCode          string
		queryParams      string
		mockTechnologies []models.Technology
		mockTotal        int64
		mockError        error
		expectedStatus   int
		expectedPage     int
		expectedLimit    int
		checkResponse    func(*testing.T, map[string]interface{})
	}{
		{
			name:        "Success - Get technologies by country",
			urlCode:     "TUR",
			queryParams: "",
			mockTechnologies: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
				{ID: 2, Name: "Bayraktar TB2"},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedPage:   0,
			expectedLimit:  0,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:             "Error - Service error",
			urlCode:          "USA",
			queryParams:      "",
			mockTechnologies: []models.Technology{},
			mockTotal:        0,
			mockError:        errors.New("database error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedPage:     0,
			expectedLimit:    0,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "failed to fetch technologies")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			// Setup mock expectations
			mockService.On("GetByCountry", mock.Anything, tt.urlCode, tt.expectedPage, tt.expectedLimit).
				Return(tt.mockTechnologies, tt.mockTotal, tt.mockError)

			// Create request
			req := httptest.NewRequest("GET", "/api/v1/technologies/country/"+tt.urlCode+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			// Call handler
			handler.GetByCountry(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Parse and check response
			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_GetByCategory(t *testing.T) {
	tests := []struct {
		name             string
		urlCategory      string
		mockTechnologies []models.Technology
		mockTotal        int64
		mockError        error
		expectedStatus   int
		checkResponse    func(*testing.T, map[string]interface{})
	}{
		{
			name:        "Success - Get technologies by category",
			urlCategory: "Aircraft",
			mockTechnologies: []models.Technology{
				{ID: 1, Name: "KAAN Fighter"},
			},
			mockTotal:      1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 1)
			},
		},
		{
			name:             "Error - Service error",
			urlCategory:      "Drones",
			mockTechnologies: []models.Technology{},
			mockTotal:        0,
			mockError:        errors.New("database error"),
			expectedStatus:   http.StatusInternalServerError,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			mockService.On("GetByCategory", mock.Anything, tt.urlCategory, 0, 0).
				Return(tt.mockTechnologies, tt.mockTotal, tt.mockError)

			req := httptest.NewRequest("GET", "/api/v1/technologies/category/"+tt.urlCategory, nil)
			rr := httptest.NewRecorder()

			handler.GetByCategory(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_GetByStatus(t *testing.T) {
	tests := []struct {
		name             string
		urlStatus        string
		mockTechnologies []models.Technology
		mockTotal        int64
		mockError        error
		expectedStatus   int
		checkResponse    func(*testing.T, map[string]interface{})
	}{
		{
			name:      "Success - Get technologies by status",
			urlStatus: "current",
			mockTechnologies: []models.Technology{
				{ID: 1, Name: "F-35", Status: "current"},
			},
			mockTotal:      1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 1)
			},
		},
		{
			name:             "Error - Invalid status",
			urlStatus:        "invalid",
			mockTechnologies: []models.Technology{},
			mockTotal:        0,
			mockError:        errors.New("invalid status"),
			expectedStatus:   http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "invalid status")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			mockService.On("GetByStatus", mock.Anything, tt.urlStatus, 0, 0).
				Return(tt.mockTechnologies, tt.mockTotal, tt.mockError)

			req := httptest.NewRequest("GET", "/api/v1/technologies/status/"+tt.urlStatus, nil)
			rr := httptest.NewRecorder()

			handler.GetByStatus(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_GetByYearRange(t *testing.T) {
	tests := []struct {
		name             string
		queryParams      string
		mockTechnologies []models.Technology
		mockTotal        int64
		mockError        error
		expectedStatus   int
		expectedStart    int
		expectedEnd      int
		checkResponse    func(*testing.T, map[string]interface{})
	}{
		{
			name:        "Success - Valid year range",
			queryParams: "?start=2010&end=2020",
			mockTechnologies: []models.Technology{
				{ID: 1, Name: "F-35"},
			},
			mockTotal:      1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedStart:  2010,
			expectedEnd:    2020,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 1)
			},
		},
		{
			name:             "Error - Missing start year",
			queryParams:      "?end=2020",
			mockTechnologies: []models.Technology{},
			mockTotal:        0,
			mockError:        nil,
			expectedStatus:   http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "start and end year are required")
			},
		},
		{
			name:             "Error - Invalid start year format",
			queryParams:      "?start=abc&end=2020",
			mockTechnologies: []models.Technology{},
			mockTotal:        0,
			mockError:        nil,
			expectedStatus:   http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "invalid start year")
			},
		},
		{
			name:             "Error - Invalid year range from service",
			queryParams:      "?start=2020&end=2010",
			mockTechnologies: []models.Technology{},
			mockTotal:        0,
			mockError:        errors.New("start year cannot be greater than end year"),
			expectedStatus:   http.StatusBadRequest,
			expectedStart:    2020,
			expectedEnd:      2010,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			// Only setup mock for valid requests
			if tt.expectedStart != 0 && tt.expectedEnd != 0 {
				mockService.On("GetByYearRange", mock.Anything, tt.expectedStart, tt.expectedEnd, 0, 0).
					Return(tt.mockTechnologies, tt.mockTotal, tt.mockError)
			}

			req := httptest.NewRequest("GET", "/api/v1/technologies/year-range"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			handler.GetByYearRange(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Success - Create technology",
			requestBody: models.Technology{
				Name:       "Test Aircraft",
				CountryID:  1,
				CategoryID: 1,
				Status:     "concept",
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "Test Aircraft", data["name"])
			},
		},
		{
			name:           "Error - Invalid JSON",
			requestBody:    "invalid json",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "invalid request body")
			},
		},
		{
			name: "Error - Validation error",
			requestBody: models.Technology{
				Name: "Test",
				// Missing CountryID and CategoryID
			},
			mockError:      errors.New("country ID is required"),
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			// Setup mock expectations (only for valid requests)
			if _, ok := tt.requestBody.(models.Technology); ok && tt.name != "Error - Validation error" {
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Technology")).
					Return(tt.mockError)
			} else if tt.name == "Error - Validation error" {
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Technology")).
					Return(tt.mockError)
			}

			// Create request body
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/api/v1/technologies", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Create(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		urlID          string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:  "Success - Update technology",
			urlID: "1",
			requestBody: models.Technology{
				Name:       "Updated KAAN",
				CountryID:  1,
				CategoryID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "Updated KAAN", data["name"])
			},
		},
		{
			name:  "Error - Invalid ID",
			urlID: "abc",
			requestBody: models.Technology{
				Name:       "Test",
				CountryID:  1,
				CategoryID: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "invalid technology ID")
			},
		},
		{
			name:           "Error - Invalid JSON",
			urlID:          "1",
			requestBody:    "invalid json",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "invalid request body")
			},
		},
		{
			name:  "Error - Technology not found",
			urlID: "999",
			requestBody: models.Technology{
				Name:       "Test",
				CountryID:  1,
				CategoryID: 1,
			},
			mockError:      errors.New("technology not found"),
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "technology not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			// Setup mock expectations (only for valid requests)
			if tt.urlID != "abc" && tt.name != "Error - Invalid JSON" {
				id := uint(1)
				if tt.urlID == "999" {
					id = 999
				}
				mockService.On("Update", mock.Anything, id, mock.AnythingOfType("*models.Technology")).
					Return(tt.mockError)
			}

			// Create request body
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("PUT", "/api/v1/technologies/"+tt.urlID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Update(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestTechnologyHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		urlID          string
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "Success - Delete technology",
			urlID:          "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "technology deleted successfully", data["message"])
			},
		},
		{
			name:           "Error - Invalid ID",
			urlID:          "abc",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "invalid technology ID")
			},
		},
		{
			name:           "Error - Technology not found",
			urlID:          "999",
			mockError:      errors.New("technology not found"),
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "technology not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTechnologyService)
			handler := NewTechnologyHandler(mockService)

			// Setup mock expectations (only for valid IDs)
			if tt.urlID != "abc" {
				id := uint(1)
				if tt.urlID == "999" {
					id = 999
				}
				mockService.On("Delete", mock.Anything, id).
					Return(tt.mockError)
			}

			req := httptest.NewRequest("DELETE", "/api/v1/technologies/"+tt.urlID, nil)
			rr := httptest.NewRecorder()

			handler.Delete(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}
