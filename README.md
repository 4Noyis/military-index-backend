# Military Index Backend

Military technology database and API service built with Go microservices.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Make

### Local Setup

1. **Clone the repository**
```bash
   git clone https://github.com/yourusername/military-index-backend.git
   cd military-index-backend
```

2. **Setup environment variables**
```bash
   cp .env.example .env
   # Edit .env with your configuration
```

3. **Start all services**
```bash
   make up-dev
```

4. **Verify installation**
```bash
   make verify
```

5. **Access services**
   - API Gateway: http://localhost:8080
   - pgAdmin: http://localhost:5050
   - API Docs: http://localhost:8080/swagger

## 📚 Documentation

- [Setup Guide](docs/SETUP.md)
- [API Documentation](docs/API.md)
- [Contributing Guide](docs/CONTRIBUTING.md)
- [Deployment Guide](docs/DEPLOYMENT.md)

## 🏗️ Architecture
```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
┌──────▼──────────┐
│  API Gateway    │
│   Port: 8080    │
└──────┬──────────┘
       │
       ├─────────────────┬─────────────────┐
       │                 │                 │
┌──────▼──────────┐ ┌───▼──────────┐ ┌────▼─────────┐
│ Tech Service    │ │Country Service│ │ PostgreSQL   │
│  Port: 8081     │ │  Port: 8082   │ │  Port: 5432  │
└─────────────────┘ └───────────────┘ └──────────────┘
```

## 🛠️ Development

### Run specific service
```bash
make run-gateway
make run-tech-service
make run-country-service
```

### Run tests
```bash
make test
```

### Database operations
```bash
make db-migrate-up      # Run migrations
make db-migrate-down    # Rollback migrations
make db-seed            # Seed data
make db-reset           # Reset database
```

## 🌍 Environments

- **Development**: Local Docker containers
- **Staging**: https://api-staging.military-index.com
- **Production**: https://api.military-index.com

## 📦 Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Container**: Docker
- **CI/CD**: GitHub Actions

## 🤝 Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md)
