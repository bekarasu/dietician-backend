package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type MedicalUpload struct {
	ID          int64  `db:"id" json:"id"`
	UserID      string `db:"user_id" json:"userId"`
	UploadType  string `db:"upload_type" json:"uploadType"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Status      string `db:"status" json:"status"`
}

type MedicalFileMetadata struct {
	ID          int64  `db:"id" json:"id"`
	UploadID    int64  `db:"upload_id" json:"uploadId"`
	FileName    string `db:"file_name" json:"fileName"`
	FileSize    int64  `db:"file_size" json:"fileSize"`
	ContentType string `db:"content_type" json:"contentType"`
	StorageKey  string `db:"storage_key" json:"storageKey"`
}

type IMedicalRepository interface {
	CreateUpload(ctx context.Context, upload *MedicalUpload) error
	GetUploadsByUserID(ctx context.Context, userID string) ([]MedicalUpload, error)
	GetUploadByID(ctx context.Context, uploadID int64) (*MedicalUpload, error)
	DeleteUpload(ctx context.Context, uploadID int64) error

	CreateFileMetadata(ctx context.Context, meta *MedicalFileMetadata) error
	GetFileMetadataByUploadID(ctx context.Context, uploadID int64) ([]MedicalFileMetadata, error)
}

type medicalRepository struct {
	db *sqlx.DB
}

func NewMedicalRepository(db *sqlx.DB) IMedicalRepository {
	return &medicalRepository{db: db}
}

func (r *medicalRepository) CreateUpload(ctx context.Context, upload *MedicalUpload) error {
	query := `INSERT INTO medical_uploads (user_id, upload_type, title, description, status) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, upload.UserID, upload.UploadType, upload.Title, upload.Description, upload.Status).Scan(&upload.ID)
}

func (r *medicalRepository) GetUploadsByUserID(ctx context.Context, userID string) ([]MedicalUpload, error) {
	var uploads []MedicalUpload
	query := `SELECT id, user_id, upload_type, title, description, status FROM medical_uploads WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &uploads, query, userID)
	return uploads, err
}

func (r *medicalRepository) GetUploadByID(ctx context.Context, uploadID int64) (*MedicalUpload, error) {
	var upload MedicalUpload
	query := `SELECT id, user_id, upload_type, title, description, status FROM medical_uploads WHERE id = $1`
	err := r.db.GetContext(ctx, &upload, query, uploadID)
	return &upload, err
}

func (r *medicalRepository) DeleteUpload(ctx context.Context, uploadID int64) error {
	query := `DELETE FROM medical_uploads WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, uploadID)
	return err
}

func (r *medicalRepository) CreateFileMetadata(ctx context.Context, meta *MedicalFileMetadata) error {
	query := `INSERT INTO medical_file_metadata (upload_id, file_name, file_size, content_type, storage_key) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, meta.UploadID, meta.FileName, meta.FileSize, meta.ContentType, meta.StorageKey).Scan(&meta.ID)
}

func (r *medicalRepository) GetFileMetadataByUploadID(ctx context.Context, uploadID int64) ([]MedicalFileMetadata, error) {
	var metadata []MedicalFileMetadata
	query := `SELECT id, upload_id, file_name, file_size, content_type, storage_key FROM medical_file_metadata WHERE upload_id = $1`
	err := r.db.SelectContext(ctx, &metadata, query, uploadID)
	return metadata, err
}
