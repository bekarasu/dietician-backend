# Auth Service

Handles user authentication and authorization.

## Endpoints

| Method | Path                      | Auth | Description          |
| ------ | ------------------------- | ---- | -------------------- |
| POST   | /api/v1/auth/register     | No   | Register new user    |
| POST   | /api/v1/auth/login        | No   | Login user           |
| POST   | /api/v1/auth/refresh      | No   | Refresh access token |
| GET    | /api/v1/auth/me           | Yes  | Get current user     |
| GET    | /health                   | No   | Health check         |

## Database Tables

- `users` - User accounts (id, email, password_hash, first_name, last_name)
- `refresh_tokens` - JWT refresh tokens (user_id, token, expires_at)

## Environment Variables

- `AUTH_SERVICE_PORT` - Service port (default: 8081)
- `JWT_SECRET` - Secret key for JWT signing
- `AUTH_DB_NAME` - PostgreSQL database name
- `POSTGRES_*` - PostgreSQL connection settings
- `REDIS_*` - Redis connection settings

## Testing

```bash
cd services/account-service
go test ./...
```

## Security

- Passwords hashed with bcrypt (default cost)
- Access tokens expire in 15 minutes
- Refresh tokens expire in 7 days
- Refresh tokens stored in PostgreSQL and cleaned on refresh
