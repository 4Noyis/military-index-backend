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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCountryService is a mock implementation of CountryService
type MockCountryService struct {
	mock.Mock
}

func (m *MockCountryService) GetAll(ctx context.Context, page, limit int) ([]models.Country, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]models.Country), args.Get(1).(int64), args.Error(2)
}

func (m *MockCountryService) GetByID(ctx context.Context, id uint) (*models.Country, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryService) GetByCode(ctx context.Context, code string) (*models.Country, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryService) Search(ctx context.Context, query string, page, limit int) ([]models.Country, int64, error) {
	args := m.Called(ctx, query, page, limit)
	return args.Get(0).([]models.Country), args.Get(1).(int64), args.Error(2)
}

func (m *MockCountryService) Create(ctx context.Context, country *models.Country) error {
	args := m.Called(ctx, country)
	return args.Error(0)
}

func (m *MockCountryService) Update(ctx context.Context, id uint, country *models.Country) error {
	args := m.Called(ctx, id, country)
	return args.Error(0)
}

func (m *MockCountryService) Delete(ctx context.Context, id uint) error {
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

func TestCountryHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockCountries  []models.Country
		mockTotal      int64
		mockError      error
		expectedStatus int
		expectedPage   int
		expectedLimit  int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "Success - Get all countries",
			queryParams: "",
			mockCountries: []models.Country{
				{ID: 1, Name: "Turkey", Code: "TUR", FlagURL: "/flags/turkey.svg"},
				{ID: 2, Name: "USA", Code: "USA", FlagURL: "/flags/usa.svg"},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedPage:   0, // Handler passes 0 when not provided
			expectedLimit:  0, // Handler passes 0 when not provided
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)

				pagination := resp["pagination"].(map[string]interface{})
				assert.Equal(t, float64(1), pagination["page"])  // Normalized in response
				assert.Equal(t, float64(10), pagination["limit"]) // Normalized in response
				assert.Equal(t, float64(2), pagination["total"])
			},
		},
		{
			name:        "Success - With pagination",
			queryParams: "?page=2&limit=5",
			mockCountries: []models.Country{
				{ID: 6, Name: "Germany", Code: "DEU"},
			},
			mockTotal:      10,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedPage:   2,
			expectedLimit:  5,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				pagination := resp["pagination"].(map[string]interface{})
				assert.Equal(t, float64(2), pagination["page"])
				assert.Equal(t, float64(5), pagination["limit"])
				assert.Equal(t, float64(10), pagination["total"])
			},
		},
		{
			name:        "Success - Search query",
			queryParams: "?q=Turkey",
			mockCountries: []models.Country{
				{ID: 1, Name: "Turkey", Code: "TUR"},
			},
			mockTotal:      1,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedPage:   0, // Handler passes 0 when not provided
			expectedLimit:  0, // Handler passes 0 when not provided
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].([]interface{})
				assert.Len(t, data, 1)
			},
		},
		{
			name:           "Error - Service error",
			queryParams:    "",
			mockCountries:  []models.Country{},
			mockTotal:      0,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedPage:   0, // Handler passes 0 when not provided
			expectedLimit:  0, // Handler passes 0 when not provided
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "Failed to fetch countries")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCountryService)
			handler := NewCountryHandler(mockService)

			// Setup mock expectations
			if tt.queryParams != "" && len(tt.queryParams) > 3 && tt.queryParams[1:3] == "q=" {
				// Search request
				mockService.On("Search", mock.Anything, "Turkey", tt.expectedPage, tt.expectedLimit).
					Return(tt.mockCountries, tt.mockTotal, tt.mockError)
			} else {
				// GetAll request
				mockService.On("GetAll", mock.Anything, tt.expectedPage, tt.expectedLimit).
					Return(tt.mockCountries, tt.mockTotal, tt.mockError)
			}

			// Create request
			req := httptest.NewRequest("GET", "/api/v1/countries"+tt.queryParams, nil)
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

func TestCountryHandler_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		urlID          string
		mockCountry    *models.Country
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:  "Success - Get country by ID",
			urlID: "1",
			mockCountry: &models.Country{
				ID:      1,
				Name:    "Turkey",
				Code:    "TUR",
				FlagURL: "/flags/turkey.svg",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "Turkey", data["name"])
				assert.Equal(t, "TUR", data["code"])
			},
		},
		{
			name:           "Error - Invalid ID format",
			urlID:          "abc",
			mockCountry:    nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "Invalid country ID")
			},
		},
		{
			name:           "Error - Country not found",
			urlID:          "999",
			mockCountry:    nil,
			mockError:      errors.New("country not found"),
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "Country not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCountryService)
			handler := NewCountryHandler(mockService)

			// Setup mock expectations (only if ID is valid)
			if tt.urlID != "abc" {
				id := uint(1)
				if tt.urlID == "999" {
					id = 999
				}
				mockService.On("GetByID", mock.Anything, id).
					Return(tt.mockCountry, tt.mockError)
			}

			// Create request with mux vars
			req := httptest.NewRequest("GET", "/api/v1/countries/"+tt.urlID, nil)
			rr := httptest.NewRecorder()

			// Setup mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.urlID})

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

func TestCountryHandler_GetByCode(t *testing.T) {
	tests := []struct {
		name           string
		urlCode        string
		mockCountry    *models.Country
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:    "Success - Get country by code",
			urlCode: "TUR",
			mockCountry: &models.Country{
				ID:      1,
				Name:    "Turkey",
				Code:    "TUR",
				FlagURL: "/flags/turkey.svg",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "Turkey", data["name"])
				assert.Equal(t, "TUR", data["code"])
			},
		},
		{
			name:           "Error - Country not found",
			urlCode:        "XXX",
			mockCountry:    nil,
			mockError:      errors.New("country not found"),
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "Country not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCountryService)
			handler := NewCountryHandler(mockService)

			// Setup mock expectations
			mockService.On("GetByCode", mock.Anything, tt.urlCode).
				Return(tt.mockCountry, tt.mockError)

			// Create request with mux vars
			req := httptest.NewRequest("GET", "/api/v1/countries/code/"+tt.urlCode, nil)
			rr := httptest.NewRecorder()

			// Setup mux vars
			req = mux.SetURLVars(req, map[string]string{"code": tt.urlCode})

			// Call handler
			handler.GetByCode(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Parse and check response
			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestCountryHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Success - Create country",
			requestBody: models.Country{
				Name:    "Japan",
				Code:    "JPN",
				FlagURL: "/flags/japan.svg",
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "Japan", data["name"])
				assert.Equal(t, "JPN", data["code"])
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
				assert.Contains(t, errorObj["message"], "Invalid request body")
			},
		},
		{
			name: "Error - Validation error",
			requestBody: models.Country{
				Name: "Test",
				Code: "T", // Too short
			},
			mockError:      errors.New("country code must be 2-3 characters"),
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "Failed to create country")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCountryService)
			handler := NewCountryHandler(mockService)

			// Setup mock expectations (only for valid requests)
			if _, ok := tt.requestBody.(models.Country); ok && tt.name != "Error - Validation error" {
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Country")).
					Return(tt.mockError)
			} else if tt.name == "Error - Validation error" {
				mockService.On("Create", mock.Anything, mock.AnythingOfType("*models.Country")).
					Return(tt.mockError)
			}

			// Create request body
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			// Create request
			req := httptest.NewRequest("POST", "/api/v1/countries", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Call handler
			handler.Create(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Parse and check response
			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestCountryHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		urlID          string
		requestBody    interface{}
		mockError      error
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:  "Success - Update country",
			urlID: "1",
			requestBody: models.Country{
				Name:    "Republic of Turkey",
				Code:    "TUR",
				FlagURL: "/flags/turkey-new.svg",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "Republic of Turkey", data["name"])
			},
		},
		{
			name:           "Error - Invalid ID",
			urlID:          "abc",
			requestBody:    models.Country{Name: "Test", Code: "TST"},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "Invalid country ID")
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
				assert.Contains(t, errorObj["message"], "Invalid request body")
			},
		},
		{
			name:  "Error - Country not found",
			urlID: "999",
			requestBody: models.Country{
				Name: "Test",
				Code: "TST",
			},
			mockError:      errors.New("country not found"),
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.Contains(t, resp, "error")
				errorObj := resp["error"].(map[string]interface{})
				assert.Contains(t, errorObj["message"], "Failed to update country")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCountryService)
			handler := NewCountryHandler(mockService)

			// Setup mock expectations (only for valid requests)
			if tt.urlID != "abc" && tt.name != "Error - Invalid JSON" {
				id := uint(1)
				if tt.urlID == "999" {
					id = 999
				}
				mockService.On("Update", mock.Anything, id, mock.AnythingOfType("*models.Country")).
					Return(tt.mockError)
			}

			// Create request body
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			// Create request
			req := httptest.NewRequest("PUT", "/api/v1/countries/"+tt.urlID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Setup mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.urlID})

			// Call handler
			handler.Update(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Parse and check response
			resp := parseJSONResponse(t, rr.Body.Bytes())
			tt.checkResponse(t, resp)

			mockService.AssertExpectations(t)
		})
	}
}

func TestCountryHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		urlID          string
		mockError      error
		expectedStatus int
		checkBody      bool
	}{
		{
			name:           "Success - Delete country",
			urlID:          "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			checkBody:      false,
		},
		{
			name:           "Error - Invalid ID",
			urlID:          "abc",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			checkBody:      true,
		},
		{
			name:           "Error - Country not found",
			urlID:          "999",
			mockError:      errors.New("country not found"),
			expectedStatus: http.StatusNotFound,
			checkBody:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCountryService)
			handler := NewCountryHandler(mockService)

			// Setup mock expectations (only for valid IDs)
			if tt.urlID != "abc" {
				id := uint(1)
				if tt.urlID == "999" {
					id = 999
				}
				mockService.On("Delete", mock.Anything, id).
					Return(tt.mockError)
			}

			// Create request
			req := httptest.NewRequest("DELETE", "/api/v1/countries/"+tt.urlID, nil)
			rr := httptest.NewRecorder()

			// Setup mux vars
			req = mux.SetURLVars(req, map[string]string{"id": tt.urlID})

			// Call handler
			handler.Delete(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Check response body if needed
			if tt.checkBody {
				resp := parseJSONResponse(t, rr.Body.Bytes())
				assert.False(t, resp["success"].(bool))
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCountryHandler_GetAll_EdgeCases(t *testing.T) {
	t.Run("Empty result set", func(t *testing.T) {
		mockService := new(MockCountryService)
		handler := NewCountryHandler(mockService)

		mockService.On("GetAll", mock.Anything, 0, 0).
			Return([]models.Country{}, int64(0), nil)

		req := httptest.NewRequest("GET", "/api/v1/countries", nil)
		rr := httptest.NewRecorder()

		handler.GetAll(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		resp := parseJSONResponse(t, rr.Body.Bytes())
		assert.True(t, resp["success"].(bool))
		data := resp["data"].([]interface{})
		assert.Len(t, data, 0)

		pagination := resp["pagination"].(map[string]interface{})
		assert.Equal(t, float64(0), pagination["total"])
	})

	t.Run("Large page number", func(t *testing.T) {
		mockService := new(MockCountryService)
		handler := NewCountryHandler(mockService)

		mockService.On("GetAll", mock.Anything, 100, 10).
			Return([]models.Country{}, int64(50), nil)

		req := httptest.NewRequest("GET", "/api/v1/countries?page=100&limit=10", nil)
		rr := httptest.NewRecorder()

		handler.GetAll(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		mockService.AssertExpectations(t)
	})
}
