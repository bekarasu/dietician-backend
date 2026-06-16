package service

import (
	"context"

	"dietician.local/services/medical-service/internal/medical/repository"
	"dietician.local/services/medical-service/internal/storage"
)

type IMedicalService interface {
	CreateUpload(ctx context.Context, userID string, req CreateUploadRequest) (*repository.MedicalUpload, error)
	ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error)
	GetUploadDetail(ctx context.Context, userID string, uploadID int64) (*UploadDetailResponse, error)
	DeleteUpload(ctx context.Context, userID string, uploadID int64) error
}

type CreateUploadRequest struct {
	UploadType  string `json:"uploadType"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UploadDetailResponse struct {
	Upload   repository.MedicalUpload         `json:"upload"`
	Metadata []repository.MedicalFileMetadata `json:"metadata"`
}

type medicalService struct {
	repo     repository.IMedicalRepository
	provider storage.Provider
}

func NewMedicalService(repo repository.IMedicalRepository, provider storage.Provider) IMedicalService {
	return &medicalService{
		repo:     repo,
		provider: provider,
	}
}

func (s *medicalService) CreateUpload(ctx context.Context, userID string, req CreateUploadRequest) (*repository.MedicalUpload, error) {
	upload := &repository.MedicalUpload{
		UserID:      userID,
		UploadType:  req.UploadType,
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
	}

	err := s.repo.CreateUpload(ctx, upload)
	if err != nil {
		return nil, err
	}

	return upload, nil
}

func (s *medicalService) ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error) {
	return s.repo.GetUploadsByUserID(ctx, userID)
}

func (s *medicalService) GetUploadDetail(ctx context.Context, userID string, uploadID int64) (*UploadDetailResponse, error) {
	upload, err := s.repo.GetUploadByID(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	metadata, err := s.repo.GetFileMetadataByUploadID(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	return &UploadDetailResponse{
		Upload:   *upload,
		Metadata: metadata,
	}, nil
}

func (s *medicalService) DeleteUpload(ctx context.Context, userID string, uploadID int64) error {
	metadata, err := s.repo.GetFileMetadataByUploadID(ctx, uploadID)
	if err == nil {
		for _, m := range metadata {
			_ = s.provider.Delete(ctx, m.StorageKey)
		}
	}

	return s.repo.DeleteUpload(ctx, uploadID)
}
