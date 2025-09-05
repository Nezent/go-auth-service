# Makefile for Go Auth Service

.PHONY: help build run test clean docker migrate dev prod logs

# Variables
BINARY_NAME=auth-service
DOCKER_IMAGE=auth-service
VERSION=1.0.0

# Default target
help: ## Show this help message
	@echo "Go Auth Service - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
dev: ## Start development server with hot reload
	air

build: ## Build the application binary
	go build -o bin/$(BINARY_NAME) cmd/auth-service/main.go

run: build ## Build and run the application
	./bin/$(BINARY_NAME) server

test: ## Run all tests
	go test -v ./...

test-cover: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

benchmark: ## Run benchmark tests
	go test -bench=. -benchmem ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Docker Commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

docker-up: ## Start all services with Docker Compose
	docker-compose up --build

docker-down: ## Stop all Docker Compose services
	docker-compose down

docker-logs: ## Show Docker Compose logs
	docker-compose logs -f

# Database Commands
migrate-up: ## Apply all database migrations
	go run cmd/auth-service/main.go migrate up

migrate-down: ## Rollback last migration
	go run cmd/auth-service/main.go migrate down

migrate-status: ## Show migration status
	go run cmd/auth-service/main.go migrate status

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Please provide a migration name using NAME=migration_name"; \
		exit 1; \
	fi
	go run cmd/auth-service/main.go migrate create $(NAME) sql

# Database Management
db-up: ## Start only the database
	docker-compose up -d db

db-shell: ## Connect to database shell
	docker-compose exec db psql -U nezent -d authdb

adminer: ## Start Adminer for database management
	docker-compose up -d adminer
	@echo "Adminer available at: http://localhost:8081"

# Monitoring
metrics: ## Show available metrics endpoint
	@echo "Metrics available at: http://localhost:8080/metrics"
	@curl -s http://localhost:8080/metrics | head -20

health: ## Check service health
	@echo "Checking service health..."
	@curl -s http://localhost:8080/health | jq '.' || echo "Service not running or jq not installed"

# Code Quality
lint: ## Run linting tools
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

fmt: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

# Dependencies
deps: ## Download dependencies
	go mod download
	go mod tidy

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

# Security
security-scan: ## Run security scan with gosec
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

# Production
prod-build: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o bin/$(BINARY_NAME) cmd/auth-service/main.go

prod-docker: ## Build production Docker image
	docker build -f Dockerfile.prod -t $(DOCKER_IMAGE)-prod:$(VERSION) .

# Setup
setup: ## Initial project setup
	@echo "Setting up Go Auth Service..."
	@cp .env.example .env || echo ".env already exists"
	@echo "1. Update .env file with your configuration"
	@echo "2. Generate RSA keys: make generate-keys"
	@echo "3. Start services: make docker-up"
	@echo "4. Run migrations: make migrate-up"

generate-keys: ## Generate RSA keys for JWT
	@mkdir -p internal/infrastructure/keys
	@openssl genrsa -out internal/infrastructure/keys/private.pem 2048
	@openssl rsa -in internal/infrastructure/keys/private.pem -pubout -out internal/infrastructure/keys/public.pem
	@echo "RSA keys generated successfully!"

# Utility
logs: ## Show application logs
	docker-compose logs -f auth

restart: ## Restart the application
	docker-compose restart auth

ps: ## Show running containers
	docker-compose ps

# All-in-one commands
install: deps generate-keys setup ## Complete installation setup

start: docker-up migrate-up ## Start everything (services + migrations)

stop: docker-down ## Stop all services
