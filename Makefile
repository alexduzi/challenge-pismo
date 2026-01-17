.PHONY: help setup run build swagger test test-unit test-integration test-coverage test-coverage-html lint clean deps wire install-tools docker-build docker-run docker-stop docker-logs docker-compose-up docker-compose-down docker-compose-logs docker-clean db-up db-down db-logs db-shell db-reset db-migrate all

# Variables
APP_NAME := api-accounts
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_CONTAINER := $(APP_NAME)-container
GO_FILES := $(shell find . -name '*.go' -type f -not -path "./vendor/*")
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Help target
help:
	@echo "$(BLUE)═══════════════════════════════════════════════════════$(NC)"
	@echo "$(BLUE)  $(APP_NAME) - Available Commands$(NC)"
	@echo "$(BLUE)═══════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(GREEN)Setup:$(NC)"
	@echo "  make setup               - Initial project setup (copy .env, install deps)"
	@echo "  make install-tools       - Install development tools (wire, swag, golangci-lint)"
	@echo ""
	@echo "$(GREEN)Local Development:$(NC)"
	@echo "  make run                 - Run the application locally"
	@echo "  make build               - Build the application binary"
	@echo "  make wire                - Generate Wire dependency injection code"
	@echo "  make swagger             - Generate/regenerate Swagger documentation"
	@echo ""
	@echo "$(GREEN)Testing & Quality:$(NC)"
	@echo "  make test                - Run all tests (unit + integration)"
	@echo "  make test-unit           - Run only unit tests (fast, no DB required)"
	@echo "  make test-integration    - Run only integration tests (requires DB)"
	@echo "  make test-coverage       - Run tests with coverage report"
	@echo "  make test-coverage-html  - Generate HTML coverage report"
	@echo "  make lint                - Run golangci-lint analysis"
	@echo "  make fmt                 - Format code with gofmt"
	@echo "  make vet                 - Run go vet"
	@echo ""
	@echo "$(GREEN)Database:$(NC)"
	@echo "  make db-up               - Start PostgreSQL container"
	@echo "  make db-down             - Stop PostgreSQL container"
	@echo "  make db-logs             - View PostgreSQL logs"
	@echo "  make db-shell            - Connect to PostgreSQL shell"
	@echo "  make db-reset            - Reset database (delete and recreate)"
	@echo "  make db-migrate          - Run database migrations"
	@echo "  make db-create-migration - Create a new migration file"
	@echo ""
	@echo "$(GREEN)Docker:$(NC)"
	@echo "  make docker-build        - Build Docker image"
	@echo "  make docker-run          - Run Docker container"
	@echo "  make docker-stop         - Stop and remove Docker container"
	@echo "  make docker-logs         - View Docker container logs"
	@echo "  make docker-shell        - Open shell in running container"
	@echo "  make docker-clean        - Remove Docker images and containers"
	@echo ""
	@echo "$(GREEN)Docker Compose:$(NC)"
	@echo "  make compose-up          - Start all services with Docker Compose"
	@echo "  make compose-down        - Stop all services"
	@echo "  make compose-restart     - Restart all services"
	@echo "  make compose-logs        - View logs from all services"
	@echo "  make compose-ps          - List running services"
	@echo "  make compose-build       - Rebuild all services"
	@echo ""
	@echo "$(GREEN)Utilities:$(NC)"
	@echo "  make clean               - Clean build artifacts and test cache"
	@echo "  make deps                - Download and tidy dependencies"
	@echo "  make all                 - Run full pipeline (setup, build, test, lint)"
	@echo "  make dev                 - Start development environment (DB + App)"
	@echo ""
	@echo "$(BLUE)═══════════════════════════════════════════════════════$(NC)"

# Install development tools
install-tools:
	@echo "$(YELLOW)Installing development tools...$(NC)"
	@echo "$(BLUE)→ Installing Wire...$(NC)"
	@go install github.com/google/wire/cmd/wire@latest
	@echo "$(BLUE)→ Installing Swag...$(NC)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(BLUE)→ Installing golangci-lint...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(BLUE)→ Installing migrate...$(NC)"
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)✓ All tools installed successfully!$(NC)"

# Initial project setup
setup:
	@echo "$(YELLOW)Setting up project...$(NC)"
	@if [ ! -f .env ]; then \
		echo "$(BLUE)→ Creating .env file from .env.example...$(NC)"; \
		cp .env.example .env; \
		echo "$(GREEN)✓ .env file created successfully!$(NC)"; \
	else \
		echo "$(YELLOW)⚠ .env file already exists, skipping...$(NC)"; \
	fi
	@echo "$(BLUE)→ Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Setup complete! Run 'make run' to start the application.$(NC)"

# Wire - Dependency Injection
wire:
	@echo "$(YELLOW)Generating Wire code...$(NC)"
	@which wire > /dev/null || (echo "$(RED)✗ Wire not found. Run 'make install-tools'$(NC)" && exit 1)
	@cd cmd/api && wire
	@echo "$(GREEN)✓ Wire code generated successfully!$(NC)"

# Run the application locally
run:
	@echo "$(YELLOW)Starting application...$(NC)"
	@go run cmd/api/main.go

# Build the application
build:
	@echo "$(YELLOW)Building application...$(NC)"
	@mkdir -p bin
	@go build -ldflags="-s -w" -o bin/api cmd/api/main.go
	@echo "$(GREEN)✓ Build complete: bin/api$(NC)"

# Generate/regenerate Swagger documentation
swagger:
	@echo "$(YELLOW)Generating Swagger documentation...$(NC)"
	@which swag > /dev/null || (echo "$(RED)✗ Swag not found. Run 'make install-tools'$(NC)" && exit 1)
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "$(GREEN)✓ Swagger docs generated in docs/$(NC)"

# Format code
fmt:
	@echo "$(YELLOW)Formatting code...$(NC)"
	@gofmt -s -w $(GO_FILES)
	@echo "$(GREEN)✓ Code formatted successfully!$(NC)"

# Run go vet
vet:
	@echo "$(YELLOW)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ No issues found!$(NC)"

# Run all tests (unit + integration)
test:
	@echo "$(YELLOW)Running all tests...$(NC)"
	@go test -v ./... -count=1
	@echo "$(GREEN)✓ All tests passed!$(NC)"

# Run only unit tests (fast, no database required)
test-unit:
	@echo "$(YELLOW)Running unit tests...$(NC)"
	@go test -v -short ./internal/... -count=1
	@echo "$(GREEN)✓ Unit tests passed!$(NC)"

# Run only integration tests (requires database)
test-integration:
	@echo "$(YELLOW)Running integration tests...$(NC)"
	@echo "$(BLUE)→ Make sure database is running (make db-up)$(NC)"
	@go test -v ./test/integration/... -count=1
	@echo "$(GREEN)✓ Integration tests passed!$(NC)"

# Run tests with coverage report
test-coverage:
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	@go test -short -coverprofile=$(COVERAGE_FILE) ./internal/...
	@echo ""
	@echo "$(BLUE)Coverage summary:$(NC)"
	@go tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print "Total coverage: " $$3}'
	@rm -f $(COVERAGE_FILE)

# Generate HTML coverage report and open in browser
test-coverage-html:
	@echo "$(YELLOW)Generating HTML coverage report...$(NC)"
	@go test -short -coverprofile=$(COVERAGE_FILE) ./internal/...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)✓ Coverage report generated: $(COVERAGE_HTML)$(NC)"
	@echo "$(BLUE)→ Opening in browser...$(NC)"
	@which xdg-open > /dev/null && xdg-open $(COVERAGE_HTML) || open $(COVERAGE_HTML) || echo "$(YELLOW)Please open $(COVERAGE_HTML) manually$(NC)"

# Run linter
lint:
	@echo "$(YELLOW)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)✗ golangci-lint not found. Run 'make install-tools'$(NC)" && exit 1)
	@golangci-lint run ./...
	@echo "$(GREEN)✓ Linter checks passed!$(NC)"

# Clean build artifacts and test cache
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@go clean -testcache
	@go clean -cache
	@echo "$(GREEN)✓ Clean complete$(NC)"

# Download and tidy dependencies
deps:
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

# Database commands

# Start PostgreSQL container
db-up:
	@echo "$(YELLOW)Starting PostgreSQL...$(NC)"
	@docker-compose up -d postgres
	@echo "$(BLUE)→ Waiting for PostgreSQL to be ready...$(NC)"
	@sleep 5
	@echo "$(GREEN)✓ PostgreSQL is ready!$(NC)"

# Stop PostgreSQL container
db-down:
	@echo "$(YELLOW)Stopping PostgreSQL...$(NC)"
	@docker-compose stop postgres
	@echo "$(GREEN)✓ PostgreSQL stopped$(NC)"

# View PostgreSQL logs
db-logs:
	@echo "$(BLUE)Showing PostgreSQL logs (Ctrl+C to exit)...$(NC)"
	@docker-compose logs -f postgres

# Connect to PostgreSQL shell
db-shell:
	@echo "$(BLUE)Connecting to PostgreSQL shell...$(NC)"
	@docker-compose exec postgres psql -U postgres -d pismo_challenge

# Reset database (delete and recreate)
db-reset:
	@echo "$(YELLOW)Resetting database...$(NC)"
	@echo "$(RED)⚠ This will delete ALL data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker-compose up -d postgres; \
		sleep 5; \
		echo "$(GREEN)✓ Database reset complete!$(NC)"; \
	else \
		echo "$(YELLOW)Cancelled.$(NC)"; \
	fi

# Run database migrations
db-migrate:
	@echo "$(YELLOW)Running migrations...$(NC)"
	@docker-compose exec postgres psql -U postgres -d pismo_challenge -f /docker-entrypoint-initdb.d/001_create_tables.sql
	@echo "$(GREEN)✓ Migrations complete!$(NC)"

# Create a new migration file
db-create-migration:
	@echo "$(YELLOW)Creating new migration...$(NC)"
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	filename="migrations/$${timestamp}_$${name}.sql"; \
	touch $$filename; \
	echo "-- Migration: $$name" > $$filename; \
	echo "-- Created at: $$(date)" >> $$filename; \
	echo "" >> $$filename; \
	echo "$(GREEN)✓ Migration created: $$filename$(NC)"

# Docker commands

# Build Docker image
docker-build:
	@echo "$(YELLOW)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)✓ Docker image built successfully: $(DOCKER_IMAGE)$(NC)"

# Run Docker container
docker-run:
	@echo "$(YELLOW)Starting Docker container...$(NC)"
	@docker run -d -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)
	@echo "$(GREEN)✓ Container started!$(NC)"
	@echo "$(BLUE)→ Access at http://localhost:8080$(NC)"
	@echo "$(BLUE)→ Health check: http://localhost:8080/health$(NC)"

# Stop and remove Docker container
docker-stop:
	@echo "$(YELLOW)Stopping Docker container...$(NC)"
	@docker stop $(DOCKER_CONTAINER) 2>/dev/null || true
	@docker rm $(DOCKER_CONTAINER) 2>/dev/null || true
	@echo "$(GREEN)✓ Container stopped and removed$(NC)"

# View Docker container logs
docker-logs:
	@echo "$(BLUE)Showing container logs (Ctrl+C to exit)...$(NC)"
	@docker logs -f $(DOCKER_CONTAINER)

# Open shell in running container
docker-shell:
	@echo "$(BLUE)Opening shell in container...$(NC)"
	@docker exec -it $(DOCKER_CONTAINER) /bin/sh

# Clean up Docker resources
docker-clean: docker-stop
	@echo "$(YELLOW)Cleaning up Docker resources...$(NC)"
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@docker system prune -f
	@echo "$(GREEN)✓ Docker cleanup complete$(NC)"

# Docker Compose commands

# Start all services with Docker Compose
compose-up:
	@echo "$(YELLOW)Starting all services with Docker Compose...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ All services started!$(NC)"
	@echo "$(BLUE)→ Application: http://localhost:8080$(NC)"
	@echo "$(BLUE)→ PostgreSQL: localhost:5432$(NC)"
	@echo "$(BLUE)→ View logs: make compose-logs$(NC)"

# Stop all services
compose-down:
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ All services stopped$(NC)"

# Restart all services
compose-restart:
	@echo "$(YELLOW)Restarting all services...$(NC)"
	@docker-compose restart
	@echo "$(GREEN)✓ All services restarted$(NC)"

# View logs from all services
compose-logs:
	@echo "$(BLUE)Showing logs from all services (Ctrl+C to exit)...$(NC)"
	@docker-compose logs -f

# List running services
compose-ps:
	@echo "$(BLUE)Running services:$(NC)"
	@docker-compose ps

# Rebuild all services
compose-build:
	@echo "$(YELLOW)Rebuilding all services...$(NC)"
	@docker-compose build --no-cache
	@echo "$(GREEN)✓ All services rebuilt!$(NC)"

# Development environment
dev: db-up
	@echo "$(YELLOW)Starting development environment...$(NC)"
	@sleep 2
	@echo "$(GREEN)✓ Development environment ready!$(NC)"
	@echo "$(BLUE)→ Database is running$(NC)"
	@echo "$(BLUE)→ Run 'make run' to start the application$(NC)"

# Run full pipeline
all: clean deps fmt vet lint test build
	@echo ""
	@echo "$(GREEN)═══════════════════════════════════════════════════════$(NC)"
	@echo "$(GREEN)  ✓ All tasks completed successfully!$(NC)"
	@echo "$(GREEN)═══════════════════════════════════════════════════════$(NC)"

# Quick start (for first time setup)
quickstart: install-tools setup db-up
	@echo ""
	@echo "$(GREEN)═══════════════════════════════════════════════════════$(NC)"
	@echo "$(GREEN)  ✓ Project setup complete!$(NC)"
	@echo "$(GREEN)═══════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(BLUE)Next steps:$(NC)"
	@echo "  1. Review and update .env file"
	@echo "  2. Run 'make run' to start the application"
	@echo "  3. Access http://localhost:8080/health"
	@echo ""