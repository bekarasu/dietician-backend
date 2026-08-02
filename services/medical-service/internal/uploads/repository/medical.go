package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

type MedicalUpload struct {
	ID            string                `db:"id" json:"id"`
	UserID        string                `db:"user_id" json:"userId"`
	UploadType    string                `db:"upload_type" json:"uploadType"`
	Title         string                `db:"title" json:"title"`
	Description   string                `db:"description" json:"description"`
	Status        string                `db:"status" json:"status"`
	ParsedResults *json.RawMessage      `db:"parsed_results" json:"parsedResults,omitempty"`
	ParsedAt      *time.Time            `db:"parsed_at" json:"parsedAt,omitempty"`
	CreatedAt     time.Time             `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time             `db:"updated_at" json:"updatedAt"`
	Metadata      []MedicalFileMetadata `db:"-" json:"metadata"`
}

type MedicalFileMetadata struct {
	ID          string    `db:"id" json:"id"`
	UploadID    string    `db:"upload_id" json:"uploadId"`
	FileName    string    `db:"file_name" json:"fileName"`
	FileSize    int64     `db:"file_size" json:"fileSize"`
	ContentType string    `db:"content_type" json:"contentType"`
	StorageKey  string    `db:"storage_key" json:"storageKey"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type IMedicalRepository interface {
	CreateUpload(ctx context.Context, upload *MedicalUpload) error
	GetUploadsByUserID(ctx context.Context, userID string) ([]MedicalUpload, error)
	GetUploadByID(ctx context.Context, uploadID string) (*MedicalUpload, error)
	DeleteUpload(ctx context.Context, uploadID string) error
	CreateFileMetadata(ctx context.Context, meta *MedicalFileMetadata) error
	GetFileMetadataByUploadID(ctx context.Context, uploadID string) ([]MedicalFileMetadata, error)
	UpdateParsedResults(ctx context.Context, uploadID string, results json.RawMessage) error
	UpdateStatus(ctx context.Context, uploadID string, status string) error
}

type medicalRepository struct {
	db *sqlx.DB
}

func NewMedicalRepository(db *sqlx.DB) IMedicalRepository {
	return &medicalRepository{db: db}
}

func (r *medicalRepository) CreateUpload(ctx context.Context, upload *MedicalUpload) error {
	query := `INSERT INTO medical_uploads (user_id, upload_type, title, description, status) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, upload.UserID, upload.UploadType, upload.Title, upload.Description, upload.Status).Scan(&upload.ID, &upload.CreatedAt, &upload.UpdatedAt)
}

func (r *medicalRepository) GetUploadsByUserID(ctx context.Context, userID string) ([]MedicalUpload, error) {
	var uploads []MedicalUpload
	query := `SELECT id, user_id, upload_type, title, description, status, parsed_results, parsed_at, created_at, updated_at FROM medical_uploads WHERE user_id = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &uploads, query, userID); err != nil {
		return nil, err
	}

	for i := range uploads {
		meta, err := r.GetFileMetadataByUploadID(ctx, uploads[i].ID)
		if err == nil {
			uploads[i].Metadata = meta
		} else {
			uploads[i].Metadata = []MedicalFileMetadata{}
		}
	}

	return uploads, nil
}

func (r *medicalRepository) GetUploadByID(ctx context.Context, uploadID string) (*MedicalUpload, error) {
	var upload MedicalUpload
	query := `SELECT id, user_id, upload_type, title, description, status, parsed_results, parsed_at, created_at, updated_at FROM medical_uploads WHERE id = $1`
	if err := r.db.GetContext(ctx, &upload, query, uploadID); err != nil {
		return nil, err
	}
	
	meta, err := r.GetFileMetadataByUploadID(ctx, uploadID)
	if err == nil {
		upload.Metadata = meta
	} else {
		upload.Metadata = []MedicalFileMetadata{}
	}

	return &upload, nil
}

func (r *medicalRepository) DeleteUpload(ctx context.Context, uploadID string) error {
	query := `DELETE FROM medical_uploads WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, uploadID)
	return err
}

func (r *medicalRepository) CreateFileMetadata(ctx context.Context, meta *MedicalFileMetadata) error {
	query := `INSERT INTO medical_file_metadata (upload_id, file_name, file_size, content_type, storage_key) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, meta.UploadID, meta.FileName, meta.FileSize, meta.ContentType, meta.StorageKey).Scan(&meta.ID, &meta.CreatedAt)
}

func (r *medicalRepository) GetFileMetadataByUploadID(ctx context.Context, uploadID string) ([]MedicalFileMetadata, error) {
	var metadata []MedicalFileMetadata
	query := `SELECT id, upload_id, file_name, file_size, content_type, storage_key, created_at FROM medical_file_metadata WHERE upload_id = $1 ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &metadata, query, uploadID)
	return metadata, err
}

func (r *medicalRepository) UpdateParsedResults(ctx context.Context, uploadID string, results json.RawMessage) error {
	query := `UPDATE medical_uploads SET parsed_results = $1, parsed_at = NOW(), status = 'completed' WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, results, uploadID)
	return err
}

func (r *medicalRepository) UpdateStatus(ctx context.Context, uploadID string, status string) error {
	query := `UPDATE medical_uploads SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, uploadID)
	return err
}
