# Progress Service

Tracks weight logs, habit logs, and weekly progress summaries.

## Endpoints

| Method | Path                                      | Description          |
| ------ | ----------------------------------------- | -------------------- |
| GET    | /api/v1/progress/{userId}                 | Get progress history |
| POST   | /api/v1/progress/{userId}/weight          | Add weight log       |
| GET    | /api/v1/progress/{userId}/weekly-summary  | Get weekly summary   |
| POST   | /api/v1/progress/{userId}/habits          | Add habit log        |
| GET    | /health                                   | Health check         |

## Database Tables

- `weight_logs` - Weight measurements (weight_kg, notes, logged_at)
- `habit_logs` - Habit tracking (habit_name, completed, notes)
- `weekly_progress_summaries` - Weekly aggregation data

## Environment Variables

- `PROGRESS_SERVICE_PORT` - Service port (default: 8083)
- `PROGRESS_DB_NAME` - PostgreSQL database name
- `POSTGRES_*` - PostgreSQL connection settings
- `REDIS_*` - Redis connection settings

## Testing

```bash
cd services/progress-service
go test ./...
```
