.PHONY: help setup up down logs clean migrate-up migrate-down migrate-create seed test

help: ## Shıw this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'


setup: ## Initial setup - copy .env.example to .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo " Created .env file from .env.example"; \
		echo " Please update .env with your configuration"; \
	else \
		echo " .env file already exists"; \
	fi

up: ## Start all services
	docker-compose up -d
	@echo "Services Started"
	@echo "PostgreSQL: localhost:5432"
	@echo "pgAdmin: http://localhost:5050"

up-dev: ## Start Services with dev profile (includes pgAdmin)
	docker-compose --profile dev up -d
	@echo "Development Services Started"

down: ## Stop all Services
	docker-compose down
	@echo "Services stopped"

logs: ## Show logs
	docker-compose logs -f

clean: ## Remove all containers and volumes
	docker-compose down -v
	@echo "All data cleaned"

migrate-up: ## Run database migrations up
	docker-compose --profile tools run --rm migrate
	@echo "migrations applied"

migrate-down: ## Rollback last migration
	docker-compose --profile tools run --rm migrate -path=/migrations -database="postgresql://$$DB_USER:$$DB_PASSWORD@postgres:5432/$$DB_NAME?sslmode=disable" down 1
	@echo "Migration rolled back"

migrate-create: ## Create new migration (usage: make migrate-create name=add_users_table)
	echo "❌ Error: name parameter is required"; \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@timestamp=$$(date +%s); \
	touch migrations/$${timestamp}_$(name).up.sql; \
	touch migrations/$${timestamp}_$(name).down.sql; \
	echo " Created migration files:"; \
	echo "  - migrations/$${timestamp}_$(name).up.sql"; \
	echo "  - migrations/$${timestamp}_$(name).down.sql"

seed: ## Seed database with sample data
	docker exec -i military_index_db psql -U militaryindex -d military_index_db < scripts/seed-data.sql
	@echo "Database seeded"

db-shell: ## Open PostgreSQL shell
	docker exec -it military_index_db psql -U militaryindex -d military_index_db

db-backup: ## Backup database
	@mkdir -p backups
		docker exec military_index_db pg_dump -U militaryindex military_index_db > backups/backup_$$(date +%Y%m%d_%H%M%S).sql
		@echo "✅ Database backed up to backups/"

test: ## Run database connection test
	go run scripts/test-db.go

db-restore: ## Restore database from backup (usage: make db-restore file=backups/backup_20231114.sql)
	@if [ -z "$(file)" ]; then \
		echo "❌ Error: file parameter is required"; \
		echo "Usage: make db-restore file=backups/your_backup.sql"; \
		exit 1; \
	fi
	docker exec -i military_index_db psql -U militaryindex -d military_index_db < $(file)
	@echo " Database restored from $(file)"
