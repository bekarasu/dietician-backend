# Medical Service

Manages medical file upload metadata with a placeholder storage abstraction.

## Endpoints

| Method | Path                                           | Description         |
| ------ | ---------------------------------------------- | ------------------- |
| POST   | /api/v1/medical/{userId}/uploads               | Create upload       |
| GET    | /api/v1/medical/{userId}/uploads               | List uploads        |
| GET    | /api/v1/medical/{userId}/uploads/{uploadId}    | Get upload detail   |
| DELETE | /api/v1/medical/{userId}/uploads/{uploadId}    | Delete upload       |
| GET    | /health                                        | Health check        |

## Database Tables

- `medical_uploads` - Upload records (upload_type, title, description, status)
- `medical_file_metadata` - File metadata (file_name, file_size, content_type, storage_key)

## Storage Provider

The service uses a `StorageProvider` interface for S3-compatible storage:

```go
type StorageProvider interface {
    Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
    Delete(ctx context.Context, key string) error
}
```

Currently using a no-op implementation. Real S3/MinIO can be added later.

## Environment Variables

- `APP_PORT` - Service port (default: 8085)
- `APP_NAME` - Service name (default: medical-service)
- `POSTGRES_*` - PostgreSQL connection settings
- `REDIS_*` - Redis connection settings

## Testing

```bash
cd services/medical-service
go test ./...
```
