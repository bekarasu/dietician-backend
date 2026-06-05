# Dietician Backend

A microservice-based backend for a dietician/nutrition application built with Go.

## Architecture

```
┌──────────────────────────────────────────────────┐
│                  API Gateway (:8080)              │
│      CORS, Rate Limit, JWT, Request ID, Logging  │
└──────┬───┬───┬───┬───┬───┬───────────────────────┘
       │   │   │   │   │   │
  ┌────┘   │   │   │   │   └────┐
  ▼        ▼   ▼   ▼   ▼       ▼
┌────┐ ┌────┐┌────┐┌────┐┌────┐┌────┐
│Auth│ │Prof││Prog││Reco││Medi││Trac│
│8081│ │8082││8083││8084││8085││8086│
└──┬─┘ └──┬─┘└──┬─┘└──┬─┘└──┬─┘└──┬─┘
   └───────┴─────┴──┬──┴─────┴─────┘
                    │
         ┌──────────┼──────────┐
         ▼          ▼          ▼
   ┌──────────┐ ┌─────────┐ ┌─────────────┐
   │PostgreSQL│ │  Redis  │ │ Event Bus   │
   │  :5432   │ │  :6379  │ │ (planned)   │
   └──────────┘ └─────────┘ └─────────────┘
```

The Recommendation Service consumes data from Profile, Progress, Tracking, and Medical services via HTTP client abstractions, allowing it to generate context-aware nutrition recommendations.

## Services

| Service                | Port | Description                                 |
| ---------------------- | ---- | ------------------------------------------- |
| API Gateway            | 8080 | Request routing, CORS, rate limit, JWT auth |
| Auth Service           | 8081 | Authentication and authorization            |
| Profile Service        | 8082 | User profiles and preferences               |
| Progress Service       | 8083 | Weight and habit tracking                   |
| Recommendation Service | 8084 | AI-powered nutrition advice                 |
| Medical Service        | 8085 | Medical file uploads and metadata           |
| Tracking Service       | 8086 | Meal, hydration, and coffee logs            |

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

### Useful Commands

```bash
make up          # Build and start all services
make up-d        # Build and start in detached mode
make down        # Stop all services
make down-v      # Stop and remove volumes
make logs        # Follow all service logs
make health      # Check gateway health
make health-all  # Check all service health endpoints
make ps          # List running containers
make test        # Run all tests
make test-v      # Run all tests (verbose)
```

### Health Checks

Each service exposes `GET /health`:

```bash
# Gateway (aggregates all services)
curl http://localhost:8080/health

# Individual services
curl http://localhost:8081/health  # Auth
curl http://localhost:8082/health  # Profile
curl http://localhost:8083/health  # Progress
curl http://localhost:8084/health  # Recommendation
curl http://localhost:8085/health  # Medical
curl http://localhost:8086/health  # Tracking
```

### Database Migrations

Migrations run automatically when each service starts. Each service manages its own database tables in its own PostgreSQL database.

### Environment Variables

See [.env.example](.env.example) for all available configuration options.

## Project Structure

```
.
├── gateway/              # API Gateway (CORS, rate limit, JWT, routing)
├── services/
│   ├── account-service/     # Authentication (register, login, JWT) and User profiles and dietary preferences
│   ├── progress-service/ # Weight logs and habit tracking
│   ├── recommendation-service/  # AI recommendations with service clients
│   ├── medical-service/  # Medical file upload metadata
│   └── tracking-service/ # Meal, hydration, coffee tracking
├── packages/             # Shared Go packages
│   ├── auth/             # Auth middleware helpers
│   ├── config/           # Database and env config
│   ├── errors/           # Structured error types
│   ├── events/           # Event publisher abstraction
│   ├── health/           # Health check utilities
│   ├── httpclient/       # HTTP client factory
│   ├── logger/           # Structured logging (zap)
│   └── response/         # JSON response helpers
├── infrastructure/       # DB init scripts
├── docker-compose.yml
├── Makefile
├── go.work
└── .env.example
```

## Security

- JWT authentication through the API Gateway
- Password hashing with bcrypt
- CORS headers configured on the gateway
- Rate limiting (120 requests/minute per IP)
- Request body size limit (10 MB)
- File upload content type validation placeholder
- Secrets via environment variables (never committed)
- Request ID propagation across services

## Event System

The project includes an event publisher abstraction (`packages/events`) with topic definitions for all domain events:

- `UserRegistered`, `ProfileUpdated`, `WeightLogged`
- `MealLogged`, `HydrationLogged`, `MedicalFileUploaded`
- `RecommendationGenerated`, `WeeklyProgressEvaluated`

Currently using a no-op publisher. The interface supports future integration with Kafka, RabbitMQ, NATS, or Redis Streams.

## Service-to-Service Communication

The Recommendation Service uses clean client interfaces to fetch data from other services:

- `ProfileClient` - User profile and dietary preferences
- `ProgressClient` - Weekly progress summaries
- `TrackingClient` - Daily tracking summaries
- `MedicalClient` - Medical upload metadata

Both HTTP and mock implementations are provided.

## Example API Requests

All requests go through the gateway at `http://localhost:8080`. Authenticated endpoints require `Authorization: Bearer <token>`.

### Auth

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123","first_name":"John","last_name":"Doe"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123"}'

# Refresh Token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'

# Get Current User
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <token>"
```

### Profile

```bash
# Get Profile
curl http://localhost:8080/api/v1/profiles/<userId> \
  -H "Authorization: Bearer <token>"

# Update Profile
curl -X PUT http://localhost:8080/api/v1/profiles/<userId> \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe","date_of_birth":"1990-01-15","gender":"male","height_cm":180,"weight_kg":75}'

# Update Preferences
curl -X PUT http://localhost:8080/api/v1/profiles/<userId>/preferences \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"preferences":["vegetarian"],"allergies":[{"allergy":"peanuts","severity":"severe"}],"disliked_foods":["liver"]}'
```

### Tracking

```bash
# Add Meal
curl -X POST http://localhost:8080/api/v1/tracking/<userId>/meals \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"meal_type":"lunch","description":"Grilled Chicken","calories":550,"protein_g":40,"carbs_g":20,"fat_g":15}'

# Add Hydration
curl -X POST http://localhost:8080/api/v1/tracking/<userId>/hydration \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"amount_ml":500}'

# Add Coffee
curl -X POST http://localhost:8080/api/v1/tracking/<userId>/coffee \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"coffee_type":"espresso","amount_ml":30}'

# Daily Summary
curl http://localhost:8080/api/v1/tracking/<userId>/daily-summary \
  -H "Authorization: Bearer <token>"
```

### Progress

```bash
# Add Weight Log
curl -X POST http://localhost:8080/api/v1/progress/<userId>/weight \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"weight_kg":84.5}'

# Add Habit
curl -X POST http://localhost:8080/api/v1/progress/<userId>/habits \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"habit_name":"drink_water","completed":true}'

# Weekly Summary
curl http://localhost:8080/api/v1/progress/<userId>/weekly-summary \
  -H "Authorization: Bearer <token>"
```

### Recommendations

```bash
# Generate Daily Recommendations
curl -X POST http://localhost:8080/api/v1/recommendations/daily \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"user_id":"<userId>"}'

# Get Recommendation History
curl http://localhost:8080/api/v1/recommendations/<userId>/history \
  -H "Authorization: Bearer <token>"
```

### Medical

```bash
# Create Upload
curl -X POST http://localhost:8080/api/v1/medical/<userId>/uploads \
  -H "Authorization: Bearer <token>" -H "Content-Type: application/json" \
  -d '{"upload_type":"blood_test","title":"Annual Blood Work","description":"Fasting blood test","files":[{"file_name":"results.pdf","file_size":102400,"content_type":"application/pdf"}]}'

# List Uploads
curl http://localhost:8080/api/v1/medical/<userId>/uploads \
  -H "Authorization: Bearer <token>"

# Get Upload Detail
curl http://localhost:8080/api/v1/medical/<userId>/uploads/<uploadId> \
  -H "Authorization: Bearer <token>"

# Delete Upload
curl -X DELETE http://localhost:8080/api/v1/medical/<userId>/uploads/<uploadId> \
  -H "Authorization: Bearer <token>"
```

## Future Roadmap

- Real AI provider integrations (OpenAI, Claude)
- S3-compatible file storage for medical uploads
- Event bus integration (Kafka/NATS/Redis Streams)
- Kubernetes deployment manifests
- API rate limiting per user
- Caching layer with Redis
- WebSocket support for real-time updates
- Admin dashboard APIs
