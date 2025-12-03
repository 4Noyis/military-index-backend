# Architecture Documentation - Military Index Backend

## Table of Contents
1. [Overview](#overview)
2. [System Architecture](#system-architecture)
3. [Service Architecture](#service-architecture)
4. [Data Models](#data-models)
5. [API Design](#api-design)
6. [Technology Stack](#technology-stack)
7. [Design Decisions](#design-decisions)
8. [Security Considerations](#security-considerations)
9. [Scalability & Performance](#scalability--performance)

---

## Overview

The Military Index Backend is a **microservices-based API** for managing and querying military technology data. The system provides a RESTful API for accessing information about military equipment, organized by country and category.

### Key Features
- RESTful API with comprehensive filtering and search
- Microservices architecture for scalability
- PostgreSQL database with optimized indexing
- Rate limiting and API gateway protection
- Comprehensive test coverage (>85%)
- Docker containerization for easy deployment

### Design Philosophy
- **Separation of Concerns**: Each service handles a specific domain
- **API Gateway Pattern**: Single entry point for all clients
- **Repository Pattern**: Clean separation between business logic and data access
- **RESTful Design**: Standard HTTP methods and status codes
- **Test-Driven Development**: High test coverage ensures reliability

---

## System Architecture

### High-Level Architecture

```
┌─────────────┐
│   Clients   │
│ (Web/Mobile)│
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│   API Gateway       │◄─── Rate Limiting
│   (Port 8080)       │◄─── CORS
│                     │◄─── Logging
└──────┬──────┬───────┘
       │      │
       ▼      ▼
┌──────────┐  ┌──────────┐
│ Country  │  │Technology│
│ Service  │  │ Service  │
│(Port 8082)  │(Port 8081)
└─────┬────┘  └─────┬────┘
      │             │
      └─────┬───────┘
            ▼
    ┌───────────────┐
    │  PostgreSQL   │
    │  (Port 5432)  │
    └───────────────┘
```

### Components

#### 1. API Gateway
- **Technology**: Go (Gin framework)
- **Responsibilities**:
  - Request routing to appropriate microservices
  - Rate limiting (100 requests/minute per IP)
  - CORS handling
  - Request logging
  - Health check aggregation
- **Port**: 8080

#### 2. Country Service
- **Technology**: Go (Gin framework)
- **Responsibilities**:
  - CRUD operations for countries
  - Country search and filtering
  - Country code lookups
- **Port**: 8082
- **Endpoints**: 6 endpoints

#### 3. Technology Service
- **Technology**: Go (Gin framework)
- **Responsibilities**:
  - CRUD operations for military technologies
  - Complex filtering (by country, category, status, year range)
  - Full-text search
  - Technology statistics
- **Port**: 8081
- **Endpoints**: 10 endpoints

#### 4. Database
- **Technology**: PostgreSQL 16
- **Features**:
  - ACID compliance
  - Full-text search (GIN indexes)
  - Connection pooling
  - Automatic timestamps
- **Port**: 5432

---

## Service Architecture

Each microservice follows a **3-layer architecture** for clean separation of concerns:

### Layer Structure

```
┌─────────────────────────────────────┐
│         HTTP Handler Layer          │
│  (Request validation, Response      │
│   formatting, HTTP error handling)  │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│        Service/Business Layer       │
│  (Business logic, Validation,       │
│   Coordination, Logging)            │
└────────────────┬────────────────────┘
                 │
┌────────────────▼────────────────────┐
│        Repository/Data Layer        │
│  (Database operations, Queries,     │
│   Transaction management)           │
└────────────────┬────────────────────┘
                 │
            ┌────▼────┐
            │Database │
            └─────────┘
```

### Country Service Structure

```
services/country-service/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   ├── handler/
│   │   ├── country_handler.go     # HTTP handlers
│   │   └── country_handler_test.go
│   ├── service/
│   │   ├── country_service.go     # Business logic
│   │   └── country_service_test.go
│   ├── repository/
│   │   ├── country_repository.go  # Data access
│   │   └── country_repository_test.go
│   └── router/
│       └── router.go              # Route definitions
├── Dockerfile                     # Container configuration
└── go.mod                        # Dependencies
```

### Technology Service Structure

```
services/technology-service/
├── cmd/
│   └── main.go                    # Entry point
├── internal/
│   ├── handler/
│   │   ├── technology_handler.go  # HTTP handlers
│   │   └── technology_handler_test.go
│   ├── service/
│   │   ├── technology_service.go  # Business logic
│   │   └── technology_service_test.go
│   ├── repository/
│   │   ├── technology_repository.go  # Data access
│   │   └── technology_repository_test.go
│   └── router/
│       └── router.go              # Route definitions
├── Dockerfile
└── go.mod
```

### Shared Code Structure

```
shared/
├── config/
│   └── database.go               # Database connection setup
├── models/
│   ├── country.go                # Country model
│   ├── technology.go             # Technology model
│   └── category.go               # Category model
├── middleware/
│   ├── cors.go                   # CORS middleware
│   ├── logger.go                 # Request logging
│   └── error.go                  # Error handling
└── utils/
    └── response.go               # Standardized responses
```

---

## Data Models

### Entity Relationship Diagram

```
┌──────────────┐
│  countries   │
├──────────────┤
│ id (PK)      │◄─────┐
│ name         │      │
│ code (UQ)    │      │
│ flag_url     │      │
│ created_at   │      │
│ updated_at   │      │
└──────────────┘      │
                      │
                      │ country_id (FK)
                      │
┌──────────────────┐  │      ┌─────────────────┐
│ technologies     │──┘      │ tech_categories │
├──────────────────┤         ├─────────────────┤
│ id (PK)          │────────►│ id (PK)         │
│ country_id (FK)  │         │ name (UQ)       │
│ category_id (FK) │         │ description     │
│ name             │         │ icon_url        │
│ description      │         │ created_at      │
│ designer         │         └─────────────────┘
│ year_developed   │
│ year_deployed    │
│ manufacturer     │
│ unit_cost        │
│ mass             │
│ length           │
│ width            │
│ height           │
│ status           │
│ image_url        │
│ created_at       │
│ updated_at       │
└──────────────────┘
```

### Country Model
```go
type Country struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`      // Unique
    Code      string    `json:"code"`      // Unique, 2-3 chars
    FlagURL   string    `json:"flag_url"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Technology Model
```go
type Technology struct {
    ID            uint      `json:"id"`
    CountryID     uint      `json:"country_id"`
    CategoryID    uint      `json:"category_id"`
    Name          string    `json:"name"`
    Description   string    `json:"description"`
    Designer      string    `json:"designer"`
    YearDeveloped *int      `json:"year_developed"`
    YearDeployed  *int      `json:"year_deployed"`
    Manufacturer  string    `json:"manufacturer"`
    UnitCost      *int64    `json:"unit_cost"`
    Mass          *int      `json:"mass"`        // kg
    Length        *int      `json:"length"`      // mm
    Width         *int      `json:"width"`       // mm
    Height        *int      `json:"height"`      // mm
    Status        string    `json:"status"`      // historical, current, future, concept, prototype
    ImageURL      string    `json:"image_url"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`

    // Relationships
    Country  Country  `json:"country,omitempty"`
    Category Category `json:"category,omitempty"`
}
```

### Category Model
```go
type Category struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`        // Unique
    Description string    `json:"description"`
    IconURL     string    `json:"icon_url"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### Database Constraints
- **Foreign Keys**:
  - `technologies.country_id` → `countries.id` (CASCADE on delete)
  - `technologies.category_id` → `tech_categories.id` (RESTRICT on delete)
- **Unique Constraints**:
  - `countries.name`
  - `countries.code`
  - `tech_categories.name`
- **Indexes**:
  - Primary keys (all tables)
  - Foreign keys (technologies)
  - `technologies.year_developed`
  - `technologies.year_deployed`
  - `technologies.status`
  - GIN index on `technologies.name` (full-text search)

---

## API Design

### RESTful Principles

The API follows REST best practices:
- **Resource-based URLs**: `/api/v1/countries`, `/api/v1/technologies`
- **HTTP verbs**: GET, POST, PUT, DELETE
- **HTTP status codes**: 200, 201, 400, 404, 500
- **JSON responses**: Standardized format
- **Pagination**: Query parameters `?page=1&limit=10`
- **Filtering**: Query parameters for search and filters
- **Versioning**: `/api/v1/` prefix for future compatibility

### Standard Response Format

#### Success Response
```json
{
  "success": true,
  "data": { ... },
  "message": "Optional success message"
}
```

#### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": "Additional error details (optional)"
  }
}
```

#### Paginated Response
```json
{
  "success": true,
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### HTTP Status Codes

| Code | Usage |
|------|-------|
| 200  | Successful GET, PUT, DELETE |
| 201  | Successful POST (resource created) |
| 204  | Successful DELETE (no content) |
| 400  | Bad request (validation error) |
| 404  | Resource not found |
| 429  | Rate limit exceeded |
| 500  | Internal server error |

### Pagination

All list endpoints support pagination:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)

### Sorting

Technology endpoints support sorting:
- `sort`: Field to sort by (`name`, `year_developed`, `created_at`, etc.)
- `order`: Sort order (`asc` or `desc`)

---

## Technology Stack

### Backend
- **Language**: Go 1.23
- **Framework**: Gin (HTTP router and middleware)
- **ORM**: GORM (Go Object-Relational Mapping)
- **Database**: PostgreSQL 16
- **Testing**: testify (assertions and mocks)

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Database Admin**: pgAdmin 4 (development)

### Development Tools
- **Package Manager**: Go modules
- **Environment Variables**: godotenv
- **Version Control**: Git
- **CI/CD**: Ready for GitHub Actions

### Key Dependencies

```go
// Web framework
github.com/gin-gonic/gin

// Database
gorm.io/gorm
gorm.io/driver/postgres
github.com/lib/pq

// Testing
github.com/stretchr/testify

// Rate limiting
golang.org/x/time/rate

// Environment
github.com/joho/godotenv
```

---

## Design Decisions

### 1. Microservices Architecture

**Decision**: Split into separate services (Country, Technology) instead of a monolith.

**Rationale**:
- **Scalability**: Services can be scaled independently
- **Maintainability**: Easier to understand and modify
- **Technology flexibility**: Can use different tech stacks per service
- **Team structure**: Different teams can own different services

**Trade-offs**:
- More complex deployment
- Network latency between services
- Data consistency challenges

### 2. API Gateway Pattern

**Decision**: Use a gateway as the single entry point.

**Rationale**:
- **Single point of entry**: Clients only need one URL
- **Cross-cutting concerns**: Rate limiting, logging, CORS in one place
- **Service abstraction**: Can change backend services without affecting clients
- **Security**: Centralized authentication/authorization point

**Implementation**: Reverse proxy using `httputil.ReverseProxy`

### 3. Repository Pattern

**Decision**: Separate data access logic from business logic.

**Rationale**:
- **Testability**: Easy to mock data layer
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Can swap database implementations
- **Single responsibility**: Each layer has one job

**Structure**:
```
Handler → Service → Repository → Database
```

### 4. PostgreSQL Database

**Decision**: Use PostgreSQL instead of NoSQL or other SQL databases.

**Rationale**:
- **ACID compliance**: Strong consistency guarantees
- **Relationships**: Natural fit for relational data
- **Full-text search**: Built-in GIN indexes
- **Mature ecosystem**: Well-supported, battle-tested
- **JSON support**: Can store semi-structured data if needed

### 5. GORM ORM

**Decision**: Use GORM instead of raw SQL.

**Rationale**:
- **Productivity**: Faster development
- **Type safety**: Compile-time checks
- **Migrations**: Built-in migration support
- **Associations**: Automatic relationship handling
- **SQL injection protection**: Parameterized queries

**Trade-offs**:
- Slight performance overhead
- Learning curve for complex queries
- Less control over generated SQL

### 6. Docker Containerization

**Decision**: Containerize all services.

**Rationale**:
- **Consistency**: Same environment everywhere
- **Isolation**: Services don't interfere
- **Portability**: Easy to deploy anywhere
- **Development**: Quick setup for new developers

### 7. Rate Limiting

**Decision**: Implement rate limiting at the API gateway.

**Rationale**:
- **Protection**: Prevent abuse and DDoS
- **Fair usage**: Ensure all clients get service
- **Cost control**: Limit resource consumption
- **Gateway level**: Single implementation for all services

**Implementation**: 100 requests/minute per IP using token bucket algorithm

### 8. Shared Code Module

**Decision**: Create a `shared` module for common code.

**Rationale**:
- **DRY principle**: Don't repeat code
- **Consistency**: Same models and utilities everywhere
- **Maintainability**: Single place to update common logic

**Contents**: Models, database config, middleware, utilities

---

## Security Considerations

### Current Security Measures

1. **SQL Injection Protection**
   - GORM uses parameterized queries
   - No raw SQL string concatenation

2. **CORS Configuration**
   - Properly configured CORS headers
   - Restricts cross-origin requests

3. **Rate Limiting**
   - 100 requests/minute per IP
   - Prevents brute force and DDoS

4. **Input Validation**
   - Service layer validates all inputs
   - Type checking via Go's type system
   - Range validation for years, limits, etc.

5. **Error Handling**
   - No sensitive information in error messages
   - Structured error responses
   - Logging for debugging

### Recommended Enhancements

1. **Authentication & Authorization**
   - JWT tokens for API access
   - Role-based access control (RBAC)
   - API keys for external clients

2. **HTTPS/TLS**
   - Encrypt data in transit
   - Use Let's Encrypt for certificates

3. **Database Security**
   - Encrypt data at rest
   - Use strong passwords
   - Limit database user permissions
   - Regular backups

4. **API Security**
   - Request signing
   - API versioning
   - Input sanitization
   - Output encoding

5. **Monitoring**
   - Log all API requests
   - Alert on suspicious patterns
   - Track failed authentication attempts

---

## Scalability & Performance

### Current Optimizations

1. **Database Indexing**
   - Primary key indexes
   - Foreign key indexes
   - Full-text search indexes (GIN)
   - Query optimization with proper indexes

2. **Connection Pooling**
   - Database connection pool (max 25 connections)
   - Reuse connections across requests
   - Automatic connection management

3. **Pagination**
   - Limit results per request
   - Prevent large data transfers
   - Configurable page sizes

4. **Docker Multi-Stage Builds**
   - Smaller image sizes
   - Faster deployments
   - Reduced attack surface

### Scalability Strategy

#### Horizontal Scaling
```
           ┌──────────────┐
           │ Load Balancer│
           └───────┬──────┘
                   │
        ┏━━━━━━━━━━┻━━━━━━━━━━┓
        ▼                     ▼
┌───────────────┐     ┌───────────────┐
│ API Gateway 1 │     │ API Gateway 2 │
└───────┬───────┘     └───────┬───────┘
        │                     │
    ┌───┴───┐             ┌───┴───┐
    ▼       ▼             ▼       ▼
 Service Service      Service Service
 Instance Instance    Instance Instance
```

**Benefits**:
- Add more instances as load increases
- No single point of failure
- Better resource utilization

#### Vertical Scaling
- Increase CPU/memory per instance
- Simpler than horizontal scaling
- Limited by hardware constraints

### Performance Recommendations

1. **Caching**
   - Redis for frequently accessed data
   - Response caching at gateway level
   - Query result caching

2. **Database Optimization**
   - Query optimization
   - Index tuning
   - Database replication (read replicas)
   - Partitioning for large tables

3. **CDN**
   - Cache static assets
   - Reduce latency globally
   - Offload traffic from origin

4. **Monitoring**
   - Prometheus metrics
   - Grafana dashboards
   - Alert on performance degradation

5. **Load Balancing**
   - Distribute traffic across instances
   - Health checks for automatic failover
   - Session affinity if needed

---

## Monitoring & Observability

### Recommended Metrics

1. **Application Metrics**
   - Request rate (requests/second)
   - Response time (p50, p95, p99)
   - Error rate
   - Active connections

2. **Business Metrics**
   - API calls per endpoint
   - Most queried technologies/countries
   - Search queries
   - User patterns

3. **Infrastructure Metrics**
   - CPU usage
   - Memory usage
   - Disk I/O
   - Network throughput

4. **Database Metrics**
   - Query performance
   - Connection pool usage
   - Slow queries
   - Lock contention

### Logging Strategy

**Log Levels**:
- **ERROR**: Service failures, critical issues
- **WARN**: Degraded performance, recoverable errors
- **INFO**: Normal operations, API calls
- **DEBUG**: Detailed debugging information

**Log Format**: Structured JSON logs for easy parsing

---

## Future Enhancements

1. **GraphQL API**
   - Allow clients to request exactly what they need
   - Reduce over-fetching and under-fetching

2. **WebSocket Support**
   - Real-time updates
   - Live data feeds

3. **Message Queue**
   - Asynchronous processing
   - Event-driven architecture
   - Better decoupling

4. **Service Mesh**
   - Istio or Linkerd
   - Advanced traffic management
   - Observability

5. **API Analytics**
   - Usage tracking
   - Popular endpoints
   - Client insights

---

**Last Updated**: December 3, 2025
**Version**: 1.0
**Maintained By**: Backend Team
