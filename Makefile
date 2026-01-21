# Force Bash shell for all commands
SHELL := /bin/bash

.PHONY: help setup run build swagger test test-unit test-integration test-coverage test-coverage-html lint clean deps wire install-tools docker-build docker-run docker-stop docker-logs docker-compose-up docker-compose-down docker-compose-logs docker-clean db-up db-down db-logs db-shell db-reset db-create-migration db-migrate-up db-migrate-down db-migrate-status db-migrate-force compose-up compose-down compose-restart compose-logs compose-ps compose-build dev all quickstart fmt vet

# Variables
APP_NAME := api-accounts
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_CONTAINER := $(APP_NAME)-container
GO_FILES := $(shell find . -name '*.go' -type f -not -path "./vendor/*")
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html
MIGRATIONS_DIR := ./migrations
DB_URL := postgresql://postgres:postgres@localhost:5432/accounts_db?sslmode=disable

# Colors for output (use with echo -e or printf)
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
CYAN := \033[0;36m
BOLD := \033[1m
NC := \033[0m

# Default target
.DEFAULT_GOAL := help

# Help target
help:
	@echo -e "$(BLUE)═══════════════════════════════════════════════════════$(NC)"
	@echo -e "$(BLUE)  $(APP_NAME) - Available Commands$(NC)"
	@echo -e "$(BLUE)═══════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo -e "$(GREEN)Setup:$(NC)"
	@echo "  make setup               - Initial project setup (copy .env, install deps)"
	@echo "  make install-tools       - Install development tools (wire, swag, golangci-lint)"
	@echo "  make quickstart          - Complete first-time setup (tools + setup + db)"
	@echo ""
	@echo -e "$(GREEN)Local Development:$(NC)"
	@echo "  make run                 - Run the application locally"
	@echo "  make build               - Build the application binary"
	@echo "  make wire                - Generate Wire dependency injection code"
	@echo "  make swagger             - Generate/regenerate Swagger documentation"
	@echo "  make dev                 - Start development environment (DB + ready to run)"
	@echo ""
	@echo -e "$(GREEN)Testing & Quality:$(NC)"
	@echo "  make test                - Run all tests (unit + integration)"
	@echo "  make test-unit           - Run only unit tests (fast, no DB required)"
	@echo "  make test-integration    - Run only integration tests (requires DB)"
	@echo "  make test-coverage       - Run tests with coverage report"
	@echo "  make test-coverage-html  - Generate HTML coverage report"
	@echo "  make lint                - Run golangci-lint analysis"
	@echo "  make fmt                 - Format code with gofmt"
	@echo "  make vet                 - Run go vet"
	@echo ""
	@echo -e "$(GREEN)Database:$(NC)"
	@echo "  make db-up               - Start PostgreSQL container"
	@echo "  make db-down             - Stop PostgreSQL container"
	@echo "  make db-logs             - View PostgreSQL logs"
	@echo "  make db-shell            - Connect to PostgreSQL shell"
	@echo "  make db-reset            - Reset database (delete and recreate)"
	@echo "  make db-create-migration - Create a new migration file"
	@echo "  make db-migrate-up       - Apply all pending migrations"
	@echo "  make db-migrate-down     - Rollback last migration"
	@echo "  make db-migrate-status   - Check migration status"
	@echo "  make db-migrate-force    - Force migration to specific version"
	@echo ""
	@echo -e "$(GREEN)Docker:$(NC)"
	@echo "  make docker-build        - Build Docker image"
	@echo "  make docker-run          - Run Docker container"
	@echo "  make docker-stop         - Stop and remove Docker container"
	@echo "  make docker-logs         - View Docker container logs"
	@echo "  make docker-shell        - Open shell in running container"
	@echo "  make docker-clean        - Remove Docker images and containers"
	@echo ""
	@echo -e "$(GREEN)Docker Compose:$(NC)"
	@echo "  make compose-up          - Start all services with Docker Compose"
	@echo "  make compose-down        - Stop all services"
	@echo "  make compose-restart     - Restart all services"
	@echo "  make compose-logs        - View logs from all services"
	@echo "  make compose-ps          - List running services"
	@echo "  make compose-build       - Rebuild all services"
	@echo ""
	@echo -e "$(GREEN)Utilities:$(NC)"
	@echo "  make clean               - Clean build artifacts and test cache"
	@echo "  make deps                - Download and tidy dependencies"
	@echo "  make all                 - Run full pipeline (setup, build, test, lint)"
	@echo ""
	@echo -e "$(BLUE)═══════════════════════════════════════════════════════$(NC)"

# Install development tools
install-tools:
	@echo -e "$(YELLOW)Installing development tools...$(NC)"
	@echo -e "$(BLUE)→ Installing Wire...$(NC)"
	@go install github.com/google/wire/cmd/wire@latest
	@echo -e "$(BLUE)→ Installing Swag...$(NC)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo -e "$(BLUE)→ Installing golangci-lint...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo -e "$(BLUE)→ Installing migrate...$(NC)"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo -e "$(GREEN)✓ All tools installed successfully!$(NC)"

# Initial project setup
setup:
	@echo -e "$(YELLOW)Setting up project...$(NC)"
	@if [ ! -f .env ]; then \
		echo -e "$(BLUE)→ Creating .env file from .env.example...$(NC)"; \
		cp .env.example .env; \
		echo -e "$(GREEN)✓ .env file created successfully!$(NC)"; \
	else \
		echo -e "$(YELLOW)⚠ .env file already exists, skipping...$(NC)"; \
	fi
	@echo -e "$(BLUE)→ Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo -e "$(GREEN)✓ Setup complete! Run 'make run' to start the application.$(NC)"

# Wire - Dependency Injection
wire:
	@echo -e "$(YELLOW)Generating Wire code...$(NC)"
	@which wire > /dev/null || (echo -e "$(RED)✗ Wire not found. Run 'make install-tools'$(NC)" && exit 1)
	@cd cmd/api && wire
	@echo -e "$(GREEN)✓ Wire code generated successfully!$(NC)"

# Run the application locally
run:
	@echo -e "$(YELLOW)Starting application...$(NC)"
	@go run ./cmd/api

# Build the application
build:
	@echo -e "$(YELLOW)Building application...$(NC)"
	@mkdir -p bin
	@go build -ldflags="-s -w" -o bin/api ./cmd/api
	@echo -e "$(GREEN)✓ Build complete: bin/api$(NC)"

# Generate/regenerate Swagger documentation
swagger:
	@echo -e "$(YELLOW)Generating Swagger documentation...$(NC)"
	@which swag > /dev/null || (echo -e "$(RED)✗ Swag not found. Run 'make install-tools'$(NC)" && exit 1)
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo -e "$(GREEN)✓ Swagger docs generated in docs/$(NC)"

# Format code
fmt:
	@echo -e "$(YELLOW)Formatting code...$(NC)"
	@gofmt -s -w $(GO_FILES)
	@echo -e "$(GREEN)✓ Code formatted successfully!$(NC)"

# Run go vet
vet:
	@echo -e "$(YELLOW)Running go vet...$(NC)"
	@go vet ./...
	@echo -e "$(GREEN)✓ No issues found!$(NC)"

# Run all tests (unit + integration)
test:
	@echo -e "$(YELLOW)Running all tests...$(NC)"
	@go test -v ./... -count=1
	@echo -e "$(GREEN)✓ All tests passed!$(NC)"

# Run only unit tests (fast, no database required)
test-unit:
	@echo -e "$(YELLOW)Running unit tests...$(NC)"
	@go test -v -short ./internal/... -count=1
	@echo -e "$(GREEN)✓ Unit tests passed!$(NC)"

# Run only integration tests (requires database)
test-integration:
	@echo -e "$(YELLOW)Running integration tests...$(NC)"
	@echo -e "$(BLUE)→ Make sure database is running (make db-up)$(NC)"
	@go test -v ./test/integration/... -count=1
	@echo -e "$(GREEN)✓ Integration tests passed!$(NC)"

# Run tests with coverage report
test-coverage:
	@echo -e "$(YELLOW)Running tests with coverage...$(NC)"
	@go test -short -coverprofile=$(COVERAGE_FILE) -coverpkg=./internal/... ./internal/...
	@echo ""
	@echo -e "$(BLUE)Coverage summary:$(NC)"
	@go tool cover -func=$(COVERAGE_FILE) | grep -v "mocks.go" | grep total | awk '{print "Total coverage: " $$3}'
	@rm -f $(COVERAGE_FILE)

# Generate HTML coverage report and open in browser
test-coverage-html:
	@echo -e "$(YELLOW)Generating HTML coverage report...$(NC)"
	@go test -short -coverprofile=$(COVERAGE_FILE) -coverpkg=./internal/... ./internal/...
	@cat $(COVERAGE_FILE) | grep -v "mocks.go" > $(COVERAGE_FILE).tmp && mv $(COVERAGE_FILE).tmp $(COVERAGE_FILE)
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo -e "$(GREEN)✓ Coverage report generated: $(COVERAGE_HTML)$(NC)"
	@echo -e "$(BLUE)→ Opening in browser...$(NC)"
	@which xdg-open > /dev/null && xdg-open $(COVERAGE_HTML) || open $(COVERAGE_HTML) || echo -e "$(YELLOW)Please open $(COVERAGE_HTML) manually$(NC)"

# Run linter
lint:
	@echo -e "$(YELLOW)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo -e "$(RED)✗ golangci-lint not found. Run 'make install-tools'$(NC)" && exit 1)
	@golangci-lint run ./...
	@echo -e "$(GREEN)✓ Linter checks passed!$(NC)"

# Clean build artifacts and test cache
clean:
	@echo -e "$(YELLOW)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@go clean -testcache
	@go clean -cache
	@echo -e "$(GREEN)✓ Clean complete$(NC)"

# Download and tidy dependencies
deps:
	@echo -e "$(YELLOW)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@go mod verify
	@echo -e "$(GREEN)✓ Dependencies updated$(NC)"

# Database commands

# Start PostgreSQL container
db-up:
	@echo -e "$(YELLOW)Starting PostgreSQL...$(NC)"
	@docker-compose up -d postgres
	@echo -e "$(BLUE)→ Waiting for PostgreSQL to be ready...$(NC)"
	@sleep 5
	@docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1 || sleep 3
	@echo -e "$(GREEN)✓ PostgreSQL is ready!$(NC)"

# Stop PostgreSQL container
db-down:
	@echo -e "$(YELLOW)Stopping PostgreSQL...$(NC)"
	@docker-compose stop postgres
	@echo -e "$(GREEN)✓ PostgreSQL stopped$(NC)"

# View PostgreSQL logs
db-logs:
	@echo -e "$(BLUE)Showing PostgreSQL logs (Ctrl+C to exit)...$(NC)"
	@docker-compose logs -f postgres

# Connect to PostgreSQL shell
db-shell:
	@echo -e "$(BLUE)Connecting to PostgreSQL shell...$(NC)"
	@docker-compose exec postgres psql -U postgres -d accounts_db

# Reset database (delete and recreate)
db-reset:
	@echo -e "$(YELLOW)Resetting database...$(NC)"
	@echo -e "$(RED)⚠ This will delete ALL data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker-compose up -d postgres; \
		sleep 5; \
		echo -e "$(GREEN)✓ Database reset complete!$(NC)"; \
	else \
		echo -e "$(YELLOW)Cancelled.$(NC)"; \
	fi

# Create a new migration file
db-create-migration:
	@echo -e "$(YELLOW)Creating new migration...$(NC)"
	@which migrate > /dev/null || (echo -e "$(RED)✗ migrate not found. Run 'make install-tools'$(NC)" && exit 1)
	@read -p "Enter migration name: " name; \
	if [ -z "$$name" ]; then \
		echo -e "$(RED)✗ Migration name cannot be empty$(NC)"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name && \
	echo -e "$(GREEN)✓ Migration created in $(MIGRATIONS_DIR)$(NC)"

# Run migrations up
db-migrate-up:
	@echo -e "$(YELLOW)Running migrations up...$(NC)"
	@which migrate > /dev/null || (echo -e "$(RED)✗ migrate not found. Run 'make install-tools'$(NC)" && exit 1)
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up
	@echo -e "$(GREEN)✓ Migrations applied!$(NC)"

# Run migrations down
db-migrate-down:
	@echo -e "$(YELLOW)Rolling back last migration...$(NC)"
	@which migrate > /dev/null || (echo -e "$(RED)✗ migrate not found. Run 'make install-tools'$(NC)" && exit 1)
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1
	@echo -e "$(GREEN)✓ Migration rolled back!$(NC)"

# Check migration status
db-migrate-status:
	@echo -e "$(BLUE)Migration status:$(NC)"
	@which migrate > /dev/null || (echo -e "$(RED)✗ migrate not found. Run 'make install-tools'$(NC)" && exit 1)
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Force migration version
db-migrate-force:
	@which migrate > /dev/null || (echo -e "$(RED)✗ migrate not found. Run 'make install-tools'$(NC)" && exit 1)
	@read -p "Enter version to force: " version; \
	if [ -z "$$version" ]; then \
		echo -e "$(RED)✗ Version cannot be empty$(NC)"; \
		exit 1; \
	fi; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $$version && \
	echo -e "$(GREEN)✓ Migration forced to version $$version$(NC)"

# Docker commands

# Build Docker image
docker-build:
	@echo -e "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo -e "$(GREEN)✓ Docker image built successfully: $(DOCKER_IMAGE)$(NC)"

# Run Docker container
docker-run:
	@echo -e "$(YELLOW)Starting Docker container...$(NC)"
	@docker run -d -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)
	@echo -e "$(GREEN)✓ Container started!$(NC)"
	@echo -e "$(BLUE)→ Access at http://localhost:8080$(NC)"
	@echo -e "$(BLUE)→ Health check: http://localhost:8080/health$(NC)"

# Stop and remove Docker container
docker-stop:
	@echo -e "$(YELLOW)Stopping Docker container...$(NC)"
	@docker stop $(DOCKER_CONTAINER) 2>/dev/null || true
	@docker rm $(DOCKER_CONTAINER) 2>/dev/null || true
	@echo -e "$(GREEN)✓ Container stopped and removed$(NC)"

# View Docker container logs
docker-logs:
	@echo -e "$(BLUE)Showing container logs (Ctrl+C to exit)...$(NC)"
	@docker logs -f $(DOCKER_CONTAINER)

# Open shell in running container
docker-shell:
	@echo -e "$(BLUE)Opening shell in container...$(NC)"
	@docker exec -it $(DOCKER_CONTAINER) /bin/sh

# Clean up Docker resources
docker-clean: docker-stop
	@echo -e "$(YELLOW)Cleaning up Docker resources...$(NC)"
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@docker system prune -f
	@echo -e "$(GREEN)✓ Docker cleanup complete$(NC)"

# Docker Compose commands

# Start all services with Docker Compose
compose-up:
	@echo -e "$(YELLOW)Starting all services with Docker Compose...$(NC)"
	@docker-compose up -d
	@echo -e "$(GREEN)✓ All services started!$(NC)"
	@echo -e "$(BLUE)→ Application: http://localhost:8080$(NC)"
	@echo -e "$(BLUE)→ PostgreSQL: localhost:5432$(NC)"
	@echo -e "$(BLUE)→ View logs: make compose-logs$(NC)"

# Stop all services
compose-down:
	@echo -e "$(YELLOW)Stopping all services...$(NC)"
	@docker-compose down
	@echo -e "$(GREEN)✓ All services stopped$(NC)"

# Restart all services
compose-restart:
	@echo -e "$(YELLOW)Restarting all services...$(NC)"
	@docker-compose restart
	@echo -e "$(GREEN)✓ All services restarted$(NC)"

# View logs from all services
compose-logs:
	@echo -e "$(BLUE)Showing logs from all services (Ctrl+C to exit)...$(NC)"
	@docker-compose logs -f

# List running services
compose-ps:
	@echo -e "$(BLUE)Running services:$(NC)"
	@docker-compose ps

# Rebuild all services
compose-build:
	@echo -e "$(YELLOW)Rebuilding all services...$(NC)"
	@docker-compose build --no-cache
	@echo -e "$(GREEN)✓ All services rebuilt!$(NC)"

# Development environment
dev: db-up
	@echo -e "$(YELLOW)Starting development environment...$(NC)"
	@sleep 2
	@echo -e "$(GREEN)✓ Development environment ready!$(NC)"
	@echo -e "$(BLUE)→ Database is running$(NC)"
	@echo -e "$(BLUE)→ Run 'make run' to start the application$(NC)"

# Run full pipeline
all: clean deps fmt vet lint test build
	@echo ""
	@echo -e "$(GREEN)═══════════════════════════════════════════════════════$(NC)"
	@echo -e "$(GREEN)  ✓ All tasks completed successfully!$(NC)"
	@echo -e "$(GREEN)═══════════════════════════════════════════════════════$(NC)"

# Quick start (for first time setup)
quickstart: install-tools setup db-up
	@echo ""
	@echo -e "$(GREEN)═══════════════════════════════════════════════════════$(NC)"
	@echo -e "$(GREEN)  ✓ Project setup complete!$(NC)"
	@echo -e "$(GREEN)═══════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo -e "$(BLUE)Next steps:$(NC)"
	@echo "  1. Review and update .env file"
	@echo "  2. Run 'make db-migrate-up' to apply migrations"
	@echo "  3. Run 'make run' to start the application"
	@echo "  4. Access http://localhost:8080/health"
	@echo ""