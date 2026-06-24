package service

import (
	"context"

	"dietician.local/services/medical-service/internal/storage"
	"dietician.local/services/medical-service/internal/uploads/dto/request"
	"dietician.local/services/medical-service/internal/uploads/dto/response"
	"dietician.local/services/medical-service/internal/uploads/repository"
	"github.com/sirupsen/logrus"
)

type IMedicalService interface {
	CreateUpload(ctx context.Context, userID string, req request.CreateUploadRequest) (*repository.MedicalUpload, error)
	ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error)
	GetUploadDetail(ctx context.Context, userID string, uploadID int64) (*response.UploadDetailResponse, error)
	DeleteUpload(ctx context.Context, userID string, uploadID int64) error
}

type medicalService struct {
	repo     repository.IMedicalRepository
	provider storage.Provider
	logger   *logrus.Logger
}

func NewMedicalService(repo repository.IMedicalRepository, provider storage.Provider, l *logrus.Logger) IMedicalService {
	return &medicalService{
		repo:     repo,
		provider: provider,
		logger:   l,
	}
}

func (s *medicalService) CreateUpload(ctx context.Context, userID string, req request.CreateUploadRequest) (*repository.MedicalUpload, error) {
	upload := &repository.MedicalUpload{
		UserID:      userID,
		UploadType:  req.UploadType,
		Title:       req.Title,
		Description: req.Description,
		Status:      "pending",
	}

	s.logger.Info(upload)

	err := s.repo.CreateUpload(ctx, upload)
	if err != nil {
		s.logger.Error("Failed to create upload", err)
		return nil, err
	}

	return upload, nil
}

func (s *medicalService) ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error) {
	return s.repo.GetUploadsByUserID(ctx, userID)
}

func (s *medicalService) GetUploadDetail(ctx context.Context, userID string, uploadID int64) (*response.UploadDetailResponse, error) {
	upload, err := s.repo.GetUploadByID(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	metadata, err := s.repo.GetFileMetadataByUploadID(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	return &response.UploadDetailResponse{
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
