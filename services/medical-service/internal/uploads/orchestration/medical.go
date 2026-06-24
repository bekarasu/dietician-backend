package orchestration

import (
	"context"

	"dietician.local/services/medical-service/internal/uploads/dto/request"
	"dietician.local/services/medical-service/internal/uploads/dto/response"
	"dietician.local/services/medical-service/internal/uploads/repository"
	"dietician.local/services/medical-service/internal/uploads/service"
)

type IMedicalOrchestrator interface {
	CreateUpload(ctx context.Context, userID string, req request.CreateUploadRequest) (*repository.MedicalUpload, error)
	ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error)
	GetUploadDetail(ctx context.Context, userID string, uploadID string) (*response.UploadDetailResponse, error)
	DeleteUpload(ctx context.Context, userID string, uploadID string) error
}

type medicalOrchestrator struct {
	medService service.IMedicalService
}

func NewMedicalOrchestrator(medService service.IMedicalService) IMedicalOrchestrator {
	return &medicalOrchestrator{medService: medService}
}

func (o *medicalOrchestrator) CreateUpload(ctx context.Context, userID string, req request.CreateUploadRequest) (*repository.MedicalUpload, error) {
	return o.medService.CreateUpload(ctx, userID, req)
}

func (o *medicalOrchestrator) ListUploads(ctx context.Context, userID string) ([]repository.MedicalUpload, error) {
	return o.medService.ListUploads(ctx, userID)
}

func (o *medicalOrchestrator) GetUploadDetail(ctx context.Context, userID string, uploadID string) (*response.UploadDetailResponse, error) {
	return o.medService.GetUploadDetail(ctx, userID, uploadID)
}

func (o *medicalOrchestrator) DeleteUpload(ctx context.Context, userID string, uploadID string) error {
	return o.medService.DeleteUpload(ctx, userID, uploadID)
}
