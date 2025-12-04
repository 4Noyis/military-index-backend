# Testing Guide - Military Index Backend

## Overview

This document describes the testing strategy and how to run tests for the Military Index Backend project.

## Test Coverage Status

### Country Service ✅
- **Repository Layer**: 73.8% coverage
- **Service Layer**: 92.2% coverage
- **Handler Layer**: 97.3% coverage
- **Total**: 87.8% coverage (excluding cmd and router)

### Technology Service ✅
- **Repository Layer**: 1.8% coverage (SQLite datetime limitations with Preload)
- **Service Layer**: 77.4% coverage ✅
- **Handler Layer**: 91.8% coverage ✅
- **Total**: 84.6% coverage (service + handler layers)

### Testing Strategy

We follow the testing pyramid approach:
1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test interactions between components
3. **E2E Tests**: Test complete API flows (future milestone)

---

## Running Tests

### Run All Tests
```bash
# From project root
cd services/country-service
go test ./...
```

### Run Specific Package Tests
```bash
# Repository tests
go test ./internal/repository -v

# Service tests
go test ./internal/service -v

# Handler tests (when available)
go test ./internal/handler -v
```

### Run Tests with Coverage
```bash
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # View in browser
```

### Run Tests with Verbose Output
```bash
go test ./... -v
```

---

## Test Structure

### Repository Tests (`*_repository_test.go`)

Repository tests use an **in-memory SQLite database** for fast, isolated testing:

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    err = db.AutoMigrate(&models.Country{})
    require.NoError(t, err)

    return db
}
```

**What we test:**
- CRUD operations (Create, Read, Update, Delete)
- Pagination logic
- Search functionality (PostgreSQL-specific features skipped in SQLite)
- Error handling
- Edge cases (invalid IDs, non-existent records)

**Example test:**
```go
func TestCountryRepository_GetByID(t *testing.T) {
    db := setupTestDB(t)
    repo := NewCountryRepository(db)
    ctx := context.Background()

    countries := seedCountries(t, db)

    tests := []struct {
        name    string
        id      uint
        want    string
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
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            country, err := repo.GetByID(ctx, tt.id)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, country.Name)
            }
        })
    }
}
```

---

### Service Tests (`*_service_test.go`)

Service tests use **mocks** to isolate business logic from database:

```go
type MockCountryRepository struct {
    mock.Mock
}

func (m *MockCountryRepository) GetByID(ctx context.Context, id uint) (*models.Country, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.Country), args.Error(1)
}
```

**What we test:**
- Business logic validation
- Input sanitization and normalization
- Pagination defaults
- Error handling and logging
- Edge cases (empty strings, nil pointers, invalid ranges)

**Example test:**
```go
func TestCountryService_Create(t *testing.T) {
    tests := []struct {
        name      string
        country   *models.Country
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
            wantErr: false,
        },
        {
            name: "Missing name",
            country: &models.Country{
                Code: "TST",
            },
            wantErr: true,
            errMsg:  "country name is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockCountryRepository)

            if !tt.wantErr {
                mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Country")).
                    Return(nil)
            }

            service := NewCountryService(mockRepo)
            err := service.Create(context.Background(), tt.country)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

---

### Handler Tests (`*_handler_test.go`)

Handler tests use `httptest` to test HTTP endpoints:

```go
func TestCountryHandler_GetAll(t *testing.T) {
    mockService := new(MockCountryService)
    handler := NewCountryHandler(mockService)

    req := httptest.NewRequest("GET", "/api/v1/countries", nil)
    rr := httptest.NewRecorder()

    mockService.On("GetAll", mock.Anything, 1, 10).
        Return([]models.Country{{Name: "Turkey"}}, int64(1), nil)

    handler.GetAll(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    mockService.AssertExpectations(t)
}
```

---

## Test Dependencies

The project uses these testing libraries:

### testify
```bash
go get github.com/stretchr/testify
```

**Features:**
- `assert`: Assertions (Equal, Error, NoError, etc.)
- `require`: Assertions that stop test on failure
- `mock`: Mock objects for interfaces

### SQLite Driver (for repository tests)
```bash
go get gorm.io/driver/sqlite
```

Used for fast, in-memory database testing without Docker/PostgreSQL.

---

## Writing New Tests

### 1. Repository Tests

```go
func TestCountryRepository_YourFunction(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    repo := NewCountryRepository(db)
    ctx := context.Background()

    // Seed data if needed
    seedCountries(t, db)

    // Define test cases
    tests := []struct {
        name    string
        input   interface{}
        want    interface{}
        wantErr bool
    }{
        // Add test cases
    }

    // Run tests
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Call function
            result, err := repo.YourFunction(ctx, tt.input)

            // Assert results
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, result)
            }
        })
    }
}
```

### 2. Service Tests

```go
func TestCountryService_YourFunction(t *testing.T) {
    tests := []struct {
        name       string
        input      interface{}
        mockReturn interface{}
        mockError  error
        wantErr    bool
    }{
        // Add test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mock
            mockRepo := new(MockCountryRepository)
            mockRepo.On("MethodName", mock.Anything, mock.Anything).
                Return(tt.mockReturn, tt.mockError)

            // Create service
            service := NewCountryService(mockRepo)

            // Call function
            result, err := service.YourFunction(context.Background(), tt.input)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                // Add more assertions
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

---

## Test Best Practices

### 1. Use Table-Driven Tests
```go
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"Valid input", "input", "output", false},
    {"Invalid input", "", "", true},
}
```

### 2. Test Edge Cases
- Empty strings
- Nil pointers
- Zero values
- Invalid IDs (0, negative, very large)
- Boundary conditions

### 3. Use Descriptive Test Names
```go
// Good
t.Run("Valid_country_code", ...)
t.Run("Empty_code_returns_error", ...)

// Bad
t.Run("Test1", ...)
t.Run("Test2", ...)
```

### 4. Keep Tests Independent
- Each test should setup its own data
- Don't rely on test execution order
- Clean up after tests if needed

### 5. Test Both Success and Failure Paths
```go
tests := []struct {
    // ...
}{
    {"Success case", ...},
    {"Error case", ...},
}
```

---

## Known Limitations

### SQLite vs PostgreSQL
Repository tests use SQLite for speed, but some PostgreSQL features aren't available:

- **ILIKE** (case-insensitive search) → Use LIKE in SQLite or skip test
- **Full-text search (GIN index)** → Not available in SQLite
- **Array types** → Not supported
- **JSON operators** → Limited support
- **DateTime with Preload** → SQLite has datetime scanning issues when using GORM's Preload with models that have `time.Time` fields (like Category.CreatedAt)

**Workaround**: Skip PostgreSQL-specific tests or use integration tests with real PostgreSQL. For Technology Service, repository tests with Preload relationships are skipped due to SQLite datetime limitations.

```go
func TestCountryRepository_Search(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping PostgreSQL-specific test in short mode")
    }
    // Test code...
}
```

Run without PostgreSQL-specific tests:
```bash
go test ./... -short
```

---

## Continuous Integration

### GitHub Actions (Future)
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.23
      - run: go test ./... -v -cover
```

---

## Coverage Goals

- **Repository Layer**: > 70% ✅ (73.8% achieved)
- **Service Layer**: > 80% ✅ (92.2% achieved)
- **Handler Layer**: > 60% ✅ (97.3% achieved)
- **Overall Project**: > 70% ✅ (87.8% achieved for Country Service)

---

## Troubleshooting

### Tests Fail Locally
```bash
# Clean build cache
go clean -testcache

# Update dependencies
go mod tidy

# Run with verbose output
go test ./... -v
```

### Coverage Report Not Generating
```bash
# Ensure coverage file is created
go test ./... -coverprofile=coverage.out

# Check file exists
ls -la coverage.out

# Generate HTML report
go tool cover -html=coverage.out
```

### Mock Not Working
```bash
# Ensure testify is installed
go get github.com/stretchr/testify/mock

# Check mock implementation matches interface
# Run with verbose to see mock expectations
go test ./... -v
```

---

## Next Steps

1. ✅ Repository tests (Country Service) - 73.8% coverage
2. ✅ Service tests (Country Service) - 92.2% coverage
3. ✅ Handler tests (Country Service) - 97.3% coverage
4. ⚠️ Repository tests (Technology Service) - 1.8% coverage (limited by SQLite)
5. ✅ Service tests (Technology Service) - 77.4% coverage
6. ✅ Handler tests (Technology Service) - 91.8% coverage
7. ⏳ Integration tests (E2E API tests)
8. ⏳ Performance tests (Load testing)

---

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [GORM Testing](https://gorm.io/docs/testing.html)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

---

**Last Updated:** December 3, 2025
