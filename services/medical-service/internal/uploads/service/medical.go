package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dietician.local/packages/pdfparser"
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
	UpdateVisibility(ctx context.Context, uploadID string, isHidden bool) error
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
		upload.Metadata = append(upload.Metadata, *meta)

		// Parse blood test results from PDF.
		if req.UploadType == request.BloodTest {
			s.parseAndStoreResults(ctx, upload, req.File.Data)
		}
	}

	return upload, nil
}

// parseAndStoreResults extracts blood test data from the PDF and saves it.
// On failure it marks the upload as "parse_failed" but does not return an
// error — the upload itself is still considered successful.
func (s *medicalService) parseAndStoreResults(ctx context.Context, upload *repository.MedicalUpload, fileData []byte) {
	report, err := pdfparser.ParseBloodTestPDF(fileData)
	if err != nil {
		s.logger.WithError(err).Error("Failed to parse blood test PDF")
		_ = s.repo.UpdateStatus(ctx, upload.ID, "parse_failed")
		upload.Status = "parse_failed"
		return
	}

	resultsJSON, err := json.Marshal(report)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal parsed results")
		_ = s.repo.UpdateStatus(ctx, upload.ID, "parse_failed")
		upload.Status = "parse_failed"
		return
	}

	if err := s.repo.UpdateParsedResults(ctx, upload.ID, resultsJSON); err != nil {
		s.logger.WithError(err).Error("Failed to store parsed results")
		_ = s.repo.UpdateStatus(ctx, upload.ID, "parse_failed")
		upload.Status = "parse_failed"
		return
	}

	rawJSON := json.RawMessage(resultsJSON)
	upload.ParsedResults = &rawJSON
	upload.Status = "completed"
	s.logger.Infof("Successfully parsed blood test with %d results for upload %s", report.TotalTestCount, upload.ID)
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

func (s *medicalService) UpdateVisibility(ctx context.Context, uploadID string, isHidden bool) error {
	return s.repo.UpdateVisibility(ctx, uploadID, isHidden)
}

