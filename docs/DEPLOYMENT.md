# Deployment Guide - Military Index Backend

## Table of Contents
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Local Development Deployment](#local-development-deployment)
4. [Production Deployment](#production-deployment)
5. [Deployment Options](#deployment-options)
6. [Environment Configuration](#environment-configuration)
7. [Database Setup](#database-setup)
8. [Monitoring & Logging](#monitoring--logging)
9. [Backup & Recovery](#backup--recovery)
10. [Troubleshooting](#troubleshooting)

---

## Overview

This guide covers deployment strategies for the Military Index Backend across different environments. The application is designed to run in Docker containers for consistency and portability.

### Deployment Architecture

```
┌─────────────────────────────────────────┐
│         Load Balancer (Optional)        │
└───────────────┬─────────────────────────┘
                │
┌───────────────▼─────────────────────────┐
│           API Gateway                    │
│         (Docker Container)               │
└───────┬──────────────────┬───────────────┘
        │                  │
┌───────▼────────┐  ┌──────▼──────────────┐
│ Country        │  │ Technology          │
│ Service        │  │ Service             │
│ (Container)    │  │ (Container)         │
└────────┬───────┘  └──────┬──────────────┘
         │                 │
         └────────┬────────┘
                  │
         ┌────────▼────────┐
         │  PostgreSQL     │
         │  Database       │
         └─────────────────┘
```

---

## Prerequisites

### Required Software
- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher
- **Git**: For cloning the repository
- **Make**: For using Makefile commands (optional)

### System Requirements

**Minimum**:
- CPU: 2 cores
- RAM: 4 GB
- Disk: 10 GB free space
- Network: Outbound internet access

**Recommended** (Production):
- CPU: 4+ cores
- RAM: 8+ GB
- Disk: 50+ GB SSD
- Network: Stable connection with low latency

---

## Local Development Deployment

### Quick Start

1. **Clone the repository**
\`\`\`bash
git clone https://github.com/your-org/military-index-backend.git
cd military-index-backend
\`\`\`

2. **Configure environment variables**
\`\`\`bash
# Copy example env file
cp .env.example .env

# Edit .env with your settings
nano .env
\`\`\`

3. **Start all services with development profile**
\`\`\`bash
# Using Make
make up-dev

# Or using Docker Compose directly
docker-compose --profile dev up -d
\`\`\`

4. **Verify services are running**
\`\`\`bash
# Check container status
docker-compose ps

# Check logs
make logs

# Or
docker-compose logs -f
\`\`\`

5. **Access the services**
- API Gateway: http://localhost:8080
- Country Service: http://localhost:8082
- Technology Service: http://localhost:8081
- pgAdmin: http://localhost:5050 (dev profile only)
- PostgreSQL: localhost:5432

### Development Commands

\`\`\`bash
# Start services
make up-dev              # Start all including pgAdmin
make up                  # Start core services only

# Stop services
make down                # Stop all services

# View logs
make logs                # View all logs
docker-compose logs api-gateway -f  # Follow specific service

# Database operations
make db-shell            # Open PostgreSQL shell
make migrate-up          # Run migrations
make seed                # Seed database

# Restart a service
docker-compose restart country-service

# Rebuild after code changes
docker-compose up -d --build country-service
\`\`\`

---

## Production Deployment

### Pre-Deployment Checklist

- [ ] Environment variables configured
- [ ] Database backups enabled
- [ ] SSL/TLS certificates obtained
- [ ] Domain DNS configured
- [ ] Firewall rules set
- [ ] Monitoring tools configured
- [ ] Log aggregation setup
- [ ] Rate limits tuned
- [ ] Health checks tested
- [ ] Disaster recovery plan documented

### Production Environment Variables

Create a \`.env.production\` file:

\`\`\`bash
# Database
DB_HOST=your-db-host.com
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=STRONG_PASSWORD_HERE
DB_NAME=military_index_prod
DB_SSLMODE=require

# Services
API_GATEWAY_PORT=8080
COUNTRY_SERVICE_PORT=8082
TECHNOLOGY_SERVICE_PORT=8081

# Environment
ENV=production
LOG_LEVEL=info

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
\`\`\`

---

## Deployment Options

### Option 1: Single Server Deployment (Simple)

**Best for**: Small to medium workloads, budget-conscious deployments

\`\`\`bash
# On your server
git clone <repository>
cd military-index-backend
cp .env.example .env
# Edit .env with production values
docker-compose up -d
\`\`\`

**Pros**:
- Simple setup
- Low cost
- Easy to manage

**Cons**:
- Single point of failure
- Limited scalability
- No automatic failover

### Option 2: DigitalOcean App Platform

**Best for**: Managed deployments, automatic scaling

**Steps**:
1. Create new App on DigitalOcean
2. Connect to GitHub repository
3. Configure build settings for each service
4. Add managed PostgreSQL database
5. Set environment variables
6. Deploy!

**Pros**:
- Fully managed
- Auto-scaling
- Zero-downtime deployments
- Built-in SSL

**Cons**:
- Higher cost
- Less control

---

## Environment Configuration

### Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| \`DB_HOST\` | PostgreSQL host | localhost | Yes |
| \`DB_PORT\` | PostgreSQL port | 5432 | Yes |
| \`DB_USER\` | Database username | postgres | Yes |
| \`DB_PASSWORD\` | Database password | - | Yes |
| \`DB_NAME\` | Database name | military_index | Yes |
| \`DB_SSLMODE\` | SSL mode | disable | No |
| \`API_GATEWAY_PORT\` | Gateway port | 8080 | No |
| \`COUNTRY_SERVICE_PORT\` | Country service port | 8082 | No |
| \`TECHNOLOGY_SERVICE_PORT\` | Technology service port | 8081 | No |

---

## Database Setup

### Initial Database Setup

\`\`\`bash
# Run migrations
make migrate-up

# Or manually
docker-compose exec db psql -U postgres -d military_index -f /docker-entrypoint-initdb.d/000001_init_schema.up.sql
\`\`\`

### Seed Initial Data

\`\`\`bash
# Seed categories and sample countries
make seed
\`\`\`

---

## Monitoring & Logging

### Health Checks

All services expose health check endpoints:

\`\`\`bash
# API Gateway
curl http://localhost:8080/health

# Country Service
curl http://localhost:8082/health

# Technology Service
curl http://localhost:8081/health
\`\`\`

### Application Logs

\`\`\`bash
# View logs
docker-compose logs -f

# Filter by service
docker-compose logs -f api-gateway

# View last 100 lines
docker-compose logs --tail=100 country-service
\`\`\`

---

## Backup & Recovery

### Database Backups

\`\`\`bash
# Manual backup
make db-backup

# Automated backups (add to crontab)
0 2 * * * cd /path/to/project && make db-backup
\`\`\`

### Database Restore

\`\`\`bash
# Restore from backup
make db-restore file=backups/backup_20251203.sql
\`\`\`

---

## Troubleshooting

### Common Issues

#### 1. Services Won't Start

\`\`\`bash
# Check logs
docker-compose logs service-name

# Check if ports are in use
sudo lsof -i :8080

# Rebuild images
docker-compose build --no-cache
\`\`\`

#### 2. Database Connection Errors

\`\`\`bash
# Check database is running
docker-compose ps db

# Test connection
docker-compose exec db psql -U postgres -c "SELECT 1"

# Verify environment variables
docker-compose exec country-service env | grep DB_
\`\`\`

#### 3. Slow API Responses

\`\`\`bash
# Check system resources
docker stats

# Check database performance
docker-compose exec db psql -U postgres -d military_index
\`\`\`

---

**Last Updated**: December 3, 2025
**Version**: 1.0
**Maintained By**: DevOps Team
