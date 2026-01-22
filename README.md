# Accounts API

A RESTful API for managing accounts and transactions built with Go, following Clean Architecture principles.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
  - [Quick Start](#quick-start)
  - [Manual Setup](#manual-setup)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Database Migrations](#database-migrations)
- [Docker](#docker)
- [Makefile Commands](#makefile-commands)
- [Project Structure](#project-structure)

## Overview

This API provides functionality for:

- **Account Management**: Create and retrieve accounts with document validation
- **Transaction Processing**: Create transactions with different operation types
  - Purchase (debit)
  - Installment Purchase (debit)
  - Withdrawal (debit)
  - Payment (credit)

## Architecture

The project follows **Clean Architecture** principles with clear separation of concerns, check [Architecture](docs/architecture.mmd) file.

### Layer Responsibilities

| Layer | Responsibility |
|-------|----------------|
| **Infrastructure** | HTTP routing, middleware (logging, error handling, request ID), handlers |
| **Application** | Use cases (business orchestration), DTOs, mappers |
| **Domain** | Entities, business rules, value objects |
| **Data** | Repository implementations, database access |

## Tech Stack

- **Language**: Go 1.25+
- **Router**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **Dependency Injection**: [Wire](https://github.com/google/wire)
- **API Documentation**: [Swagger](https://github.com/swaggo/swag)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Logging**: [slog](https://go.dev/blog/slog)
- **Testing**: Go testing + testify + [Testcontainers](https://testcontainers.com/)
- **Linting**: [golangci-lint](https://golangci-lint.run/)
- **SQL connection + driver**: [sqlx](https://github.com/jmoiron/sqlx) / [pgx](https://github.com/jackc/pgx)
- **Config**: [viper](https://github.com/spf13/viper)

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose
- Make (optional but recommended)

## Getting Started

### Quick Start

The fastest way to get the project running:

```bash
# Clone the repository
git clone https://github.com/alexduzi/challengepismo.git
cd challengepismo

# Complete first-time setup (installs tools, configures env, starts DB)
make quickstart

# Apply database migrations
make db-migrate-up

# Run the application
make run
```

The API will be available at `http://localhost:8080`

### Manual Setup

If you prefer step-by-step setup:

```bash
# 1. Install development tools
make install-tools

# 2. Setup project (creates .env, downloads dependencies)
make setup

# 3. Start the database
make db-up

# 4. Apply migrations
make db-migrate-up

# 5. Run the application
make run
```

## Configuration

Copy `.env.example` to `.env` and configure as needed:

```bash
cp .env.example .env
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | `api-accounts` |
| `APP_ENV` | Environment (development/production) | `development` |
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin framework mode (debug/release) | `debug` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `accounts_db` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `DB_MAX_CONNS` | Maximum DB connections | `25` |
| `DB_MIN_CONNS` | Minimum DB connections | `5` |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `debug` |
| `LOG_FORMAT` | Log format (json/text) | `text` |

## Running the Application

### Local Development

```bash
# Start database only
make db-up

# Run the application
make run

# Or start development environment (DB + ready message)
make dev
```

### Using Docker Compose

```bash
# Start all services (API + PostgreSQL)
make compose-up

# View logs
make compose-logs

# Stop all services
make compose-down
```

### Using the Launcher Script

```bash
./scripts/run.sh
```

This interactive script will:
- Check dependencies
- Start all containers
- Wait for services to be ready
- Display available endpoints

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `GET` | `/readiness` | Readiness probe |
| `POST` | `/api/v1/accounts` | Create account |
| `GET` | `/api/v1/accounts/:accountId` | Get account by ID |
| `POST` | `/api/v1/transactions` | Create transaction |

### Swagger Documentation

Access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

### Example Requests

**Create Account:**
```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "document_number": "12345678900",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "11987654321",
    "account_type": "checking"
  }'
```

**Get Account:**
```bash
curl http://localhost:8080/api/v1/accounts/1
```

**Create Transaction:**
```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 1,
    "operation_type_id": 4,
    "amount": 100.50
  }'
```

### Operation Types

| ID | Type | Effect |
|----|------|--------|
| 1 | Purchase | Debit (negative) |
| 2 | Installment Purchase | Debit (negative) |
| 3 | Withdrawal | Debit (negative) |
| 4 | Payment | Credit (positive) |

## Testing

The project includes both unit and integration tests.

### Run All Tests

```bash
make test
```

### Run Unit Tests Only

Fast tests that don't require a database:

```bash
make test-unit
```

### Run Integration Tests

Requires Docker (uses Testcontainers):

```bash
make test-integration
```

### Test Coverage

```bash
# Generate coverage report
make test-coverage

# Generate HTML coverage report and open in browser
make test-coverage-html
```

### Code Quality

```bash
# Run linter
make lint

# Format code
make fmt

# Run go vet
make vet
```

## Database Migrations

### Create a New Migration

```bash
make db-create-migration
# Enter migration name when prompted
```

### Apply Migrations

```bash
make db-migrate-up
```

### Rollback Last Migration

```bash
make db-migrate-down
```

### Check Migration Status

```bash
make db-migrate-status
```

### Reset Database

```bash
make db-reset
```

### Access Database Shell

```bash
make db-shell
```

## Docker

### Build Image

```bash
make docker-build
```

### Run Container

```bash
make docker-run
```

### View Logs

```bash
make docker-logs
```

### Stop Container

```bash
make docker-stop
```

### Clean Up

```bash
make docker-clean
```

## Makefile Commands

Run `make help` to see all available commands:

### Setup Commands
| Command | Description |
|---------|-------------|
| `make setup` | Initial project setup (copy .env, install deps) |
| `make install-tools` | Install development tools (wire, swag, golangci-lint) |
| `make quickstart` | Complete first-time setup (tools + setup + db) |

### Development Commands
| Command | Description |
|---------|-------------|
| `make run` | Run the application locally |
| `make build` | Build the application binary |
| `make wire` | Generate Wire dependency injection code |
| `make swagger` | Generate Swagger documentation |
| `make dev` | Start development environment |

### Testing Commands
| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-unit` | Run unit tests only |
| `make test-integration` | Run integration tests |
| `make test-coverage` | Run tests with coverage |
| `make test-coverage-html` | Generate HTML coverage report |

### Quality Commands
| Command | Description |
|---------|-------------|
| `make lint` | Run golangci-lint |
| `make fmt` | Format code |
| `make vet` | Run go vet |

### Database Commands
| Command | Description |
|---------|-------------|
| `make db-up` | Start PostgreSQL |
| `make db-down` | Stop PostgreSQL |
| `make db-shell` | Connect to PostgreSQL shell |
| `make db-reset` | Reset database |
| `make db-migrate-up` | Apply migrations |
| `make db-migrate-down` | Rollback last migration |
| `make db-create-migration` | Create new migration |

### Docker Commands
| Command | Description |
|---------|-------------|
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |
| `make docker-stop` | Stop container |
| `make compose-up` | Start all services |
| `make compose-down` | Stop all services |
| `make compose-logs` | View logs |

### Utility Commands
| Command | Description |
|---------|-------------|
| `make clean` | Clean build artifacts |
| `make deps` | Download dependencies |
| `make all` | Run full pipeline (setup, build, test, lint) |

## Project Structure

```
.
├── cmd/
│   └── api/                    # Application entrypoint
│       ├── main.go             # Main function
│       ├── wire.go             # Wire dependency injection
│       └── wire_gen.go         # Generated Wire code
├── internal/
│   ├── domain/                 # Domain entities and business rules
│   ├── dto/
│   │   ├── request/            # Request DTOs
│   │   └── response/           # Response DTOs
│   ├── repository/             # Repository interfaces
│   ├── usecase/                # Use cases (business logic)
│   │   └── mapper/             # DTO <-> Entity mappers
│   └── infrastructure/
│       ├── config/             # Configuration
│       ├── http/
│       │   ├── handler/        # HTTP handlers
│       │   ├── middleware/     # HTTP middleware
│       │   └── router/         # Router setup
│       ├── logger/             # Logging
│       ├── persistence/
│       │   └── postgres/       # PostgreSQL repositories
│       ├── validator/          # Request validation
│       └── exception/          # Error handling
├── test/
│   └── integration/            # Integration tests
├── migrations/                 # Database migrations
├── docs/                       # Swagger documentation
├── scripts/                    # Helper scripts
├── Dockerfile                  # Docker build
├── docker-compose.yml          # Docker Compose config
├── Makefile                    # Make commands
├── .env.example                # Environment template
└── go.mod                      # Go modules
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Author

Alex Duzi - [duzihd@gmail.com](mailto:duzihd@gmail.com)
