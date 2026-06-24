package orchestration

import (
	"context"

	"dietician.local/services/progress-service/internal/tracking/dto"
	"dietician.local/services/progress-service/internal/tracking/service"
)

type ITrackingOrchestrator interface {
	CreateTrackingMetric(ctx context.Context, req dto.CreateTrackingMetricRequest) (*dto.TrackingMetricResponse, error)
	GetTrackingMetrics(ctx context.Context, userID string) ([]dto.TrackingMetricResponse, error)
}

type trackingOrchestrator struct {
	svc service.ITrackingService
}

func NewTrackingOrchestrator(svc service.ITrackingService) ITrackingOrchestrator {
	return &trackingOrchestrator{svc: svc}
}

func (o *trackingOrchestrator) CreateTrackingMetric(ctx context.Context, req dto.CreateTrackingMetricRequest) (*dto.TrackingMetricResponse, error) {
	return o.svc.CreateTrackingMetric(ctx, req)
}

func (o *trackingOrchestrator) GetTrackingMetrics(ctx context.Context, userID string) ([]dto.TrackingMetricResponse, error) {
	return o.svc.GetTrackingMetrics(ctx, userID)
}
