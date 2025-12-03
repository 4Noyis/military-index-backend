# Frequently Asked Questions (FAQ)

## Table of Contents
1. [General Questions](#general-questions)
2. [Setup & Installation](#setup--installation)
3. [Development](#development)
4. [Database](#database)
5. [API Usage](#api-usage)
6. [Testing](#testing)
7. [Deployment](#deployment)
8. [Troubleshooting](#troubleshooting)
9. [Performance](#performance)
10. [Security](#security)

---

## General Questions

### What is the Military Index Backend?

The Military Index Backend is a RESTful API microservices system for managing and querying military technology data. It provides comprehensive information about military equipment organized by country and category.

### What technologies does it use?

- **Backend**: Go 1.23
- **Framework**: Gin
- **Database**: PostgreSQL 16
- **ORM**: GORM
- **Containerization**: Docker
- **Testing**: testify

### Is it production-ready?

Yes! The project includes:
- ✅ Comprehensive test coverage (>85%)
- ✅ Docker containerization
- ✅ API documentation
- ✅ Rate limiting
- ✅ Error handling
- ✅ Logging
- ⚠️ Authentication (to be added)

### What's the license?

[Check the LICENSE file in the repository]

---

## Setup & Installation

### How do I install the project?

\`\`\`bash
# Clone repository
git clone https://github.com/your-org/military-index-backend.git
cd military-index-backend

# Start with Docker Compose
make up-dev
\`\`\`

See [SETUP.md](./SETUP.md) for detailed instructions.

### What are the minimum system requirements?

- **CPU**: 2 cores
- **RAM**: 4 GB
- **Disk**: 10 GB free
- **Docker**: 20.10+
- **Docker Compose**: 2.0+

### Do I need to install Go locally?

No! Everything runs in Docker containers. You only need Docker and Docker Compose installed.

However, if you want to develop locally outside Docker:
- Go 1.23 or higher
- PostgreSQL 16 (or use Docker for DB only)

### How do I update to the latest version?

\`\`\`bash
git pull
docker-compose down
docker-compose up -d --build
\`\`\`

---

## Development

### How do I run tests?

\`\`\`bash
# Country Service
cd services/country-service
go test ./...

# With coverage
go test ./... -cover

# Technology Service
cd services/technology-service
go test ./...
\`\`\`

See [TESTING.md](./TESTING.md) for comprehensive testing guide.

### How do I add a new endpoint?

1. **Add handler method** in `internal/handler/`
2. **Add service method** in `internal/service/`
3. **Add repository method** in `internal/repository/` (if needed)
4. **Register route** in `internal/router/router.go`
5. **Write tests** for all layers
6. **Update API documentation**

Example structure:
\`\`\`go
// Handler
func (h *Handler) NewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Handle request
}

// Service
func (s *Service) NewMethod(ctx context.Context) error {
    // Business logic
}

// Repository
func (r *Repository) NewQuery(ctx context.Context) error {
    // Database operation
}
\`\`\`

### How do I add a new microservice?

1. Create directory in `services/`
2. Copy structure from existing service
3. Update `docker-compose.yml`
4. Add routes to API gateway
5. Update documentation

### How do I debug a service?

\`\`\`bash
# View logs
docker-compose logs -f service-name

# Enable debug mode
# In .env
LOG_LEVEL=debug

# Run locally (outside Docker)
cd services/country-service
go run cmd/main.go
\`\`\`

### What's the code structure?

Each service follows clean architecture:
\`\`\`
cmd/           # Entry point
internal/
  handler/     # HTTP handlers
  service/     # Business logic
  repository/  # Data access
  router/      # Route definitions
\`\`\`

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed architecture.

---

## Database

### How do I access the database?

\`\`\`bash
# Using make
make db-shell

# Or directly
docker-compose exec db psql -U postgres -d military_index
\`\`\`

### How do I run migrations?

\`\`\`bash
# Up migrations
make migrate-up

# Down migrations
make migrate-down

# Create new migration
make migrate-create name=add_new_table
\`\`\`

### How do I backup the database?

\`\`\`bash
# Manual backup
make db-backup

# Automated (add to crontab)
0 2 * * * cd /path/to/project && make db-backup
\`\`\`

Backups are stored in `backups/` directory.

### How do I restore a backup?

\`\`\`bash
make db-restore file=backups/backup_20251203.sql
\`\`\`

### Can I use a different database?

The application is designed for PostgreSQL. Using a different database would require:
- Changing the GORM driver
- Updating connection configuration
- Testing compatibility (full-text search, etc.)
- Updating migrations

### What's the database schema?

See [ARCHITECTURE.md](./ARCHITECTURE.md#data-models) for complete schema and ERD.

Key tables:
- `countries`: Country information
- `tech_categories`: Technology categories
- `technologies`: Military equipment data

---

## API Usage

### Where is the API documentation?

See [API.md](./API.md) for complete API reference with examples.

### What's the base URL?

**Local**: `http://localhost:8080/api/v1/`
**Production**: `https://your-domain.com/api/v1/`

### How do I authenticate?

Currently, the API is public (no authentication required).

For production, implement:
- JWT tokens
- API keys
- OAuth 2.0

### What's the rate limit?

**Default**: 100 requests per minute per IP

Configure in `.env`:
\`\`\`bash
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
\`\`\`

### How do I handle pagination?

All list endpoints support pagination:
\`\`\`bash
# Default: page=1, limit=10
GET /api/v1/technologies

# Custom pagination
GET /api/v1/technologies?page=2&limit=20
\`\`\`

Response includes pagination info:
\`\`\`json
{
  "pagination": {
    "page": 2,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
\`\`\`

### How do I search?

Use the `q` query parameter:
\`\`\`bash
# Search countries
GET /api/v1/countries?q=Turkey

# Search technologies
GET /api/v1/technologies?q=Fighter
\`\`\`

### How do I filter technologies?

Multiple filter options:
\`\`\`bash
# By country
GET /api/v1/technologies/country/TUR

# By category
GET /api/v1/technologies/category/Aircraft

# By status
GET /api/v1/technologies/status/current

# By year range
GET /api/v1/technologies/year-range?start=2010&end=2020
\`\`\`

### What are the response formats?

All responses are JSON:

**Success**:
\`\`\`json
{
  "success": true,
  "data": { ... }
}
\`\`\`

**Error**:
\`\`\`json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  }
}
\`\`\`

### Can I use the API from a browser?

Yes! CORS is enabled. The API can be called from any origin (configurable in production).

---

## Testing

### How comprehensive is the test coverage?

- **Country Service**: 87.8%
- **Technology Service**: 84.6%
- **Overall**: >85%

See [TESTING.md](./TESTING.md) for details.

### What testing frameworks are used?

- **testify**: Assertions and mocks
- **httptest**: HTTP handler testing
- **SQLite**: In-memory database for tests
- **GORM**: ORM testing

### How do I run specific tests?

\`\`\`bash
# Run specific test function
go test -run TestCountryService_Create

# Run specific package
go test ./internal/service

# Run with coverage
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
\`\`\`

### Why do some repository tests fail?

SQLite doesn't support some PostgreSQL features:
- ILIKE (case-insensitive search)
- DateTime with Preload relationships

These tests are skipped. They work fine with real PostgreSQL.

### How do I write new tests?

Follow the existing patterns:

**Service tests** (use mocks):
\`\`\`go
func TestService_Method(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("Method", mock.Anything).Return(nil)

    service := NewService(mockRepo)
    err := service.Method(context.Background())

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
\`\`\`

**Handler tests** (use httptest):
\`\`\`go
func TestHandler_Endpoint(t *testing.T) {
    mockService := new(MockService)
    handler := NewHandler(mockService)

    req := httptest.NewRequest("GET", "/endpoint", nil)
    rr := httptest.NewRecorder()

    handler.Endpoint(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}
\`\`\`

---

## Deployment

### How do I deploy to production?

See [DEPLOYMENT.md](./DEPLOYMENT.md) for comprehensive deployment guide.

Quick steps:
1. Configure production `.env`
2. Setup database
3. Deploy containers
4. Run migrations
5. Verify health checks

### What are the deployment options?

1. **Single Server**: Simple, budget-friendly
2. **DigitalOcean App Platform**: Fully managed
3. **AWS ECS/Fargate**: Enterprise scale
4. **Kubernetes**: Maximum flexibility

### Do I need a domain?

For production, yes. For development, localhost is fine.

### How do I setup SSL/HTTPS?

**Option 1**: Use nginx reverse proxy with Let's Encrypt
**Option 2**: Use cloud provider's load balancer (handles SSL)
**Option 3**: Use Caddy (automatic HTTPS)

### How do I scale the application?

**Vertical** scaling (easier):
- Increase CPU/RAM of servers
- Upgrade database instance

**Horizontal** scaling (better):
- Deploy multiple instances of each service
- Use load balancer
- Add database read replicas

### What's the estimated cost?

Depends on deployment option:

**DigitalOcean** (managed):
- Droplet: $12/month
- Database: $15/month
- **Total**: ~$30/month

**AWS** (enterprise):
- ECS: $50-100/month
- RDS: $30-60/month
- Load Balancer: $20/month
- **Total**: ~$100-180/month

**VPS** (budget):
- Single server: $5-10/month
- **Total**: ~$10/month

---

## Troubleshooting

### Services won't start

\`\`\`bash
# Check logs
docker-compose logs

# Check ports
sudo lsof -i :8080

# Rebuild
docker-compose down -v
docker-compose up -d --build
\`\`\`

### Database connection fails

\`\`\`bash
# Verify database is running
docker-compose ps db

# Check connection
docker-compose exec db psql -U postgres -c "SELECT 1"

# Check environment variables
docker-compose exec country-service env | grep DB_
\`\`\`

### API returns 500 errors

\`\`\`bash
# Check service logs
docker-compose logs country-service

# Enable debug mode
LOG_LEVEL=debug

# Check database connectivity
make db-shell
\`\`\`

### Tests are failing

\`\`\`bash
# Clean test cache
go clean -testcache

# Update dependencies
go mod tidy

# Run with verbose output
go test ./... -v
\`\`\`

### Port already in use

\`\`\`bash
# Find process using port
sudo lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
API_GATEWAY_PORT=8081
\`\`\`

### Docker out of space

\`\`\`bash
# Remove unused containers
docker system prune -a

# Remove volumes (WARNING: deletes data)
docker system prune -a --volumes
\`\`\`

---

## Performance

### How fast is the API?

Typical response times:
- Simple queries: 10-50ms
- Complex queries: 50-200ms
- Search queries: 100-300ms

### How many requests can it handle?

Depends on resources, but typical:
- **Single server**: 100-500 req/s
- **Scaled (3 instances)**: 300-1500 req/s
- **Rate limit**: 100 req/min per IP

### How do I optimize performance?

1. **Add caching** (Redis)
2. **Optimize queries** (add indexes)
3. **Scale horizontally** (more instances)
4. **Use CDN** (for static assets)
5. **Database tuning** (connection pool, read replicas)

### Why are queries slow?

Common causes:
- Missing database indexes
- Large result sets (use pagination)
- Complex joins
- Network latency

**Debug**:
\`\`\`sql
-- Check slow queries
SELECT query, mean_time
FROM pg_stat_statements
WHERE mean_time > 100
ORDER BY mean_time DESC;
\`\`\`

### Can I cache responses?

Yes! Recommended caching strategies:
- Redis for frequently accessed data
- HTTP cache headers
- Gateway-level caching

---

## Security

### Is the API secure?

Current security measures:
- ✅ SQL injection protection (GORM)
- ✅ Rate limiting
- ✅ CORS configuration
- ✅ Input validation
- ✅ Error handling
- ⚠️ No authentication (add for production)

### How do I add authentication?

Recommended approach:
1. Implement JWT authentication
2. Add middleware to verify tokens
3. Create user management endpoints
4. Add role-based access control (RBAC)

### Should I use HTTPS?

**YES** for production! Use:
- Let's Encrypt (free SSL certificates)
- Cloud provider SSL
- Reverse proxy (nginx/Caddy)

### How do I secure the database?

1. Strong passwords
2. Firewall rules (allow only app servers)
3. SSL connections (`DB_SSLMODE=require`)
4. Regular backups
5. Limit user permissions

### Are there known vulnerabilities?

The codebase is regularly updated. Check:
- `go mod` for dependency updates
- GitHub security alerts
- OWASP Top 10 best practices

---

## Contributing

### How do I contribute?

See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

Quick steps:
1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Submit pull request

### What's the branch strategy?

- `main`: Production-ready code
- `stage`: Staging/testing
- `feature/*`: New features
- `fix/*`: Bug fixes

### How do I report bugs?

Open an issue on GitHub with:
- Description of the bug
- Steps to reproduce
- Expected vs actual behavior
- System information
- Logs (if applicable)

---

## Additional Resources

- **Setup Guide**: [SETUP.md](./SETUP.md)
- **API Documentation**: [API.md](./API.md)
- **Architecture**: [ARCHITECTURE.md](./ARCHITECTURE.md)
- **Deployment**: [DEPLOYMENT.md](./DEPLOYMENT.md)
- **Testing**: [TESTING.md](./TESTING.md)
- **Contributing**: [CONTRIBUTING.md](./CONTRIBUTING.md)
- **Roadmap**: [ROADMAP.md](./ROADMAP.md)

---

**Last Updated**: December 3, 2025
**Version**: 1.0

**Have a question not answered here?** Open an issue on GitHub!
