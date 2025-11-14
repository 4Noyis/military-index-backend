 Setup Guide for Contributors

## 🎯 Prerequisites

### Required
- **Docker Desktop**: [Download](https://www.docker.com/products/docker-desktop)
- **Git**: [Download](https://git-scm.com/downloads)
- **Go 1.21+**: [Download](https://golang.org/dl/)
- **Make**: Pre-installed on macOS/Linux, [Install on Windows](https://gnuwin32.sourceforge.net/packages/make.htm)

### Optional
- **VS Code** or **GoLand**: IDE
- **Postman** or **Insomnia**: API testing
- **TablePlus** or **DBeaver**: Database GUI

---

## 🚀 Step 1: Clone Repository
```bash
# Clone repository
git clone https://github.com/yourusername/military-index-backend.git
cd military-index-backend

# Switch to stage branch for development
git checkout stage
```

---

## 🔧 Step 2: Environment Setup

### Create `.env` file
```bash
# Copy example environment file
cp .env.example .env

# Edit with your preferred editor
nano .env
# or
code .env
```

### Minimum required configuration for local development:
```env
# Leave most defaults, just ensure these are set:
DB_USER=militaryindex
DB_PASSWORD=SecurePassword123
DB_NAME=military_index_db

# These should work as-is for local development
DB_HOST=postgres
DB_PORT=5432
API_GATEWAY_PORT=8080
```

---

## 🐳 Step 3: Start Services

### Option A: Using Make (Recommended)
```bash
# Start all services (database + all microservices)
make up-dev

# Wait 10-15 seconds for database initialization
# Watch logs to see when ready
make logs
```

### Option B: Using Docker Compose Directly
```bash
# Start with dev profile (includes pgAdmin)
docker-compose --profile dev up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

---

## ✅ Step 4: Verify Installation
```bash
# Check if all services are running
make verify

# Or manually verify:

# 1. Check database
docker exec -it military_index_db psql -U militaryindex -d military_index_db -c "\dt"

# 2. Check API Gateway
curl http://localhost:8080/health

# 3. Check Technology Service
curl http://localhost:8081/api/v1/technologies

# 4. Check Country Service
curl http://localhost:8082/api/v1/countries
```

Expected output:
```json
{
  "status": "ok",
  "database": "connected",
  "services": ["technology", "country"]
}
```

---

## 🗄️ Step 5: Access Database

### Option A: pgAdmin (Web Interface)

1. Open browser: http://localhost:5050
2. Login:
   - Email: `admin@admin.com`
   - Password: `admin`
3. Add server:
   - Host: `postgres`
   - Port: `5432`
   - Database: `military_index_db`
   - Username: `militaryindex`
   - Password: `SecurePassword123`

### Option B: Command Line
```bash
# Open PostgreSQL shell
make db-shell

# Inside psql:
\dt                          # List tables
SELECT * FROM countries;     # View countries
\q                          # Exit
```

### Option C: Database GUI Tools

**Connection Details:**
- Host: `localhost`
- Port: `5432`
- Database: `military_index_db`
- Username: `militaryindex`
- Password: `SecurePassword123`

---

## 🧪 Step 6: Run Tests
```bash
# Run all tests
make test

# Run specific service tests
cd services/technology-service
go test ./...

# Run with coverage
make test-coverage
```

---

## 🔄 Step 7: Database Operations

### Migrations
```bash
# Run migrations
make db-migrate-up

# Rollback last migration
make db-migrate-down

# Check migration status
make db-migrate-status
```

### Seed Data
```bash
# Seed sample data
make db-seed

# Reset database (drop all + recreate + seed)
make db-reset
```

---

## 💻 Step 8: Local Development

### Running Individual Services
```bash
# Terminal 1: API Gateway
cd services/api-gateway
go run cmd/main.go

# Terminal 2: Technology Service
cd services/technology-service
go run cmd/main.go

# Terminal 3: Country Service
cd services/country-service
go run cmd/main.go
```

### Hot Reload (Optional)
```bash
# Install Air for hot reload
go install github.com/cosmtrek/air@latest

# Run with hot reload
cd services/technology-service
air
```

---

## 🐛 Troubleshooting

### Issue: Port already in use
```bash
# Check what's using the port
lsof -i :5432  # macOS/Linux
netstat -ano | findstr :5432  # Windows

# Stop the service using the port or change port in .env
```

### Issue: Database connection failed
```bash
# Check if PostgreSQL container is running
docker-compose ps

# Restart database
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### Issue: Tables not found
```bash
# Reset database
make db-reset

# Or manually run migrations
make db-migrate-up
make db-seed
```

### Issue: Cannot connect to Docker
```bash
# Restart Docker Desktop
# On macOS/Windows: Quit and restart Docker Desktop

# Check Docker is running
docker ps
```

---

## 🧹 Cleanup

### Stop services (keep data)
```bash
make down
```

### Stop and remove all data
```bash
make clean
```

### Remove everything including images
```bash
docker-compose down -v --rmi all
```

---

## 📚 Next Steps

1. ✅ **Read API Documentation**: `docs/API.md`
2. ✅ **Understand Architecture**: `docs/ARCHITECTURE.md`
3. ✅ **Check Contributing Guide**: `docs/CONTRIBUTING.md`
4. ✅ **Create your first feature**: `git checkout -b feature/your-feature-name`

---

## 💡 Tips

1. **Use Makefile commands**: They're shortcuts for common tasks
2. **Check logs frequently**: `make logs` or `docker-compose logs -f [service]`
3. **Keep .env updated**: Don't commit real credentials
4. **Use branches**: Never commit directly to `main` or `stage`
5. **Write tests**: Before creating PR

---

## 🆘 Getting Help

- **Documentation**: Check `docs/` folder
- **Issues**: Create GitHub issue
- **Discord/Slack**: [Your team chat link]
- **Email**: dev@military-index.com

---

## ✅ Checklist

Before starting development, ensure:

- [ ] Docker is running
- [ ] All containers are up: `docker-compose ps`
- [ ] Database has tables: `make db-verify`
- [ ] API Gateway responds: `curl http://localhost:8080/health`
- [ ] You're on `stage` branch: `git branch`
- [ ] `.env` file exists with correct values
- [ ] You can access pgAdmin at http://localhost:5050

If all checked ✅, you're ready to develop! 🚀
