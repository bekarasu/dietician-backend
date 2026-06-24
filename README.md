# Dietician Backend

A microservice-based backend for a dietician/nutrition application built with Go.

## Architecture

```
┌──────────────────────────────────────────────────┐
│                  API Gateway (:80)               │
│      CORS, Rate Limit, JWT, Request ID, Logging  │
└──────┬────────────┬────────────┬────────────┬────┘
       │            │            │            │
  ┌────┘            │            │            └────┐
  ▼                 ▼            ▼                 ▼
┌────┐            ┌────┐       ┌────┐            ┌────┐
│Acct│            │Reco│       │Medi│            │Prog│
│8081│            │8082│       │8083│            │8084│
└──┬─┘            └──┬─┘       └──┬─┘            └──┬─┘
   └─────────┬───────┴────────────┴──────┬──────────┘
             │                           │
         ┌───┼───────────────────────────┼───┐
         ▼   ▼                           ▼   ▼
   ┌──────────┐ ┌─────────┐ ┌─────────────┐
   │PostgreSQL│ │  Redis  │ │ Event Bus   │
   │  :5432   │ │  :6379  │ │ (planned)   │
   └──────────┘ └─────────┘ └─────────────┘
```

The Recommendation Service consumes data from Account, Progress, and Medical services via HTTP client abstractions, allowing it to generate context-aware nutrition recommendations.

## Services

| Service                | Port | Description                                 |
| ---------------------- | ---- | ------------------------------------------- |
| API Gateway            | 80   | Request routing, CORS, rate limit, JWT auth |
| Account Service        | 8081 | Authentication, user profiles & preferences |
| Recommendation Service | 8082 | AI-powered nutrition advice                 |
| Medical Service        | 8083 | Medical file uploads and metadata           |
| Progress Service       | 8084 | Weight, habit tracking and logs             |

## Prerequisites

- Docker and Docker Compose
- Go 1.26+ (for local development and testing)

## Quick Start

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Start all services:
   ```bash
   docker compose up --build
   ```
   or
   ```bash
   make up
   ```

3. Verify everything is running:
   ```bash
   curl http://localhost:8080/health
   ```

## Development

### Live Reload

All Go services use [Air](https://github.com/air-verse/air) for live reload. When you modify `.go` files locally, the corresponding service automatically rebuilds and restarts inside Docker.

### Running Tests

Run all tests locally (no Docker required):

```bash
make test        # Run all tests
make test-v      # Run all tests with verbose output
go test ./...    # From individual service directory
```

### Makefile Commands

The root `Makefile` provides several commands to manage the application's lifecycle, testing, and maintenance.

| Command | Description | Underlying Command / Details |
|---------|-------------|------------------------------|
| `make up` | Builds and starts all services in the foreground. | `docker compose up --build` |
| `make up-d` | Builds and starts all services in the background (detached mode). | `docker compose up --build -d` |
| `make down` | Stops all running services and removes containers/networks. | `docker compose down` |
| `make down-v` | Stops services and removes containers, networks, and all named volumes (clears database data). | `docker compose down -v` |
| `make build` | Builds the Docker images for all services without starting them. | `docker compose build` |
| `make logs` | Follows the log output of all running services. | `docker compose logs -f` |
| `make health` | Checks the health endpoint of the API Gateway. | `curl -s http://localhost:8080/health` |
| `make health-all`| Checks the health endpoints of the Gateway and all individual microservices. | Hits `:8080` to `:8084` health endpoints |
| `make ps` | Lists all running Docker containers for the project. | `docker compose ps` |
| `make test` | Runs all Go tests across the packages and services. | `go test dietician.local/...` |
| `make test-v` | Runs all Go tests with verbose output. | `go test -v dietician.local/...` |
| `make swag` | Generates Swagger API documentation for all microservices. | Loops through `services/*` and runs `make swag` |
| `make tidy` | Cleans and updates Go module dependencies in all services. | Loops through `services/*` and runs `go mod tidy` |

### Health Checks

Each service exposes `GET /health`:

```bash
# Individual services
curl http://localhost:8081/health  # Account
curl http://localhost:8082/health  # Recommendation
curl http://localhost:8083/health  # Medical
curl http://localhost:8084/health  # Progress
```

### Database Migrations

Migrations run automatically when each service starts. Each service manages its own database tables in its own PostgreSQL database.

### Environment Variables

See [.env.example](.env.example) for all available configuration options.

## Project Structure

```
.
├── services/
│   ├── account-service/         # Authentication, user profiles & dietary preferences
│   ├── medical-service/         # Medical file upload and metadata
│   ├── progress-service/        # Weight logs and habit tracking
│   ├── recommendation-service/  # AI recommendations
├── packages/             # Shared Go packages
│   ├── config/           # Configuration structures
│   ├── constants/        # Shared constants
│   ├── errors/           # Structured error types
│   ├── events/           # Event bus / message definitions
│   ├── healthcheck/      # Health check utilities
│   ├── httpclient/       # HTTP client factory
│   ├── localizer/        # i18n & localization
│   ├── logging/          # Structured logging (zap)
│   ├── middleware/       # Common HTTP middleware
│   ├── openai/           # OpenAI API client wrapper
│   ├── response/         # JSON response helpers
│   ├── smtp/             # Email sending utilities
│   ├── swagger/          # Swagger documentation helpers
│   ├── tokenizer/        # JWT & token utilities
│   ├── utils/            # General utility functions
│   ├── validation/       # Request validation utilities
│   └── viperconfig/      # Viper-based configuration loader
├── docker-compose.dev.yml
├── docker-compose.infra.yml
├── docker-compose.prod.yml
├── Makefile
├── go.work
└── .env.example
```

## Security

- JWT authentication through the API Gateway
- Password hashing with bcrypt
- Request body size limit
- File upload content type validation placeholder
- Secrets via environment variables (never committed)
- Request ID propagation across services
