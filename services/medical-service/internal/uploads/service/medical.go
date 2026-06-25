package service

import (
	"context"
	"fmt"
	"time"

	"dietician.local/services/medical-service/internal/storage"
	"dietician.local/services/medical-service/internal/uploads/dto/request"
	"dietician.local/services/medical-service/internal/uploads/dto/response"
	"dietician.local/services/medical-service/internal/uploads/repository"
	"github.com/sirupsen/logrus"
)

type IMedicalService interface {
	CreateUpload(ctx context.Context, userID string, req request.CreateUploadRequest) (*repository.MedicalUpload, error)
	ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error)
	GetUploadDetail(ctx context.Context, userID string, uploadID string) (*response.UploadDetailResponse, error)
	DeleteUpload(ctx context.Context, userID string, uploadID string) error
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
		UploadType:  string(req.UploadType),
		Title:       fmt.Sprintf("%s-%s", string(req.UploadType), time.Now().Format(time.RFC3339)),
		Description: string(req.UploadType),
		Status:      "pending",
	}

	err := s.repo.CreateUpload(ctx, upload)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create upload")
		return nil, err
	}

	if req.File != nil {
		storageKey := fmt.Sprintf("medicaluploads/%s/%s/%s", userID, upload.ID, req.File.FileName)

		_, err := s.provider.Upload(ctx, storageKey, req.File.Data, req.File.ContentType)
		if err != nil {
			s.logger.WithError(err).Error("Failed to upload file to storage")
			return nil, err
		}

		meta := &repository.MedicalFileMetadata{
			UploadID:    upload.ID,
			FileName:    req.File.FileName,
			FileSize:    req.File.FileSize,
			ContentType: req.File.ContentType,
			StorageKey:  storageKey,
		}
		if err := s.repo.CreateFileMetadata(ctx, meta); err != nil {
			s.logger.WithError(err).Error("Failed to create file metadata")
			return nil, err
		}
	}

	return upload, nil
}

func (s *medicalService) ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error) {
	return s.repo.GetUploadsByUserID(ctx, userID)
}

func (s *medicalService) GetUploadDetail(ctx context.Context, userID string, uploadID string) (*response.UploadDetailResponse, error) {
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

func (s *medicalService) DeleteUpload(ctx context.Context, userID string, uploadID string) error {
	metadata, err := s.repo.GetFileMetadataByUploadID(ctx, uploadID)
	if err == nil {
		for _, m := range metadata {
			_ = s.provider.Delete(ctx, m.StorageKey)
		}
	}

	return s.repo.DeleteUpload(ctx, uploadID)
}
