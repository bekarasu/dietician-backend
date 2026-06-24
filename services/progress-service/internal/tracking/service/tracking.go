package service

import (
	"context"

	"dietician.local/services/progress-service/internal/tracking/dto"
	"dietician.local/services/progress-service/internal/tracking/model"
	"dietician.local/services/progress-service/internal/tracking/repository"
)

type ITrackingService interface {
	CreateTrackingMetric(ctx context.Context, req dto.CreateTrackingMetricRequest) (*dto.TrackingMetricResponse, error)
	GetTrackingMetrics(ctx context.Context, userID string) ([]dto.TrackingMetricResponse, error)
}

type trackingService struct {
	repo repository.ITrackingRepository
}

func NewTrackingService(repo repository.ITrackingRepository) ITrackingService {
	return &trackingService{repo: repo}
}

func (s *trackingService) CreateTrackingMetric(ctx context.Context, req dto.CreateTrackingMetricRequest) (*dto.TrackingMetricResponse, error) {
	metric := &model.TrackingMetric{
		UserID:       req.UserID,
		MetricType:   req.MetricType,
		TargetValue:  req.TargetValue,
		CurrentValue: req.CurrentValue,
		Deadline:     req.Deadline,
		Status:       "active",
	}

	err := s.repo.CreateTrackingMetric(ctx, metric)
	if err != nil {
		return nil, err
	}

	return s.mapToResponse(metric), nil
}

func (s *trackingService) GetTrackingMetrics(ctx context.Context, userID string) ([]dto.TrackingMetricResponse, error) {
	metrics, err := s.repo.GetTrackingMetricsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []dto.TrackingMetricResponse
	for _, m := range metrics {
		res = append(res, *s.mapToResponse(&m))
	}
	return res, nil
}

func (s *trackingService) mapToResponse(m *model.TrackingMetric) *dto.TrackingMetricResponse {
	return &dto.TrackingMetricResponse{
		ID:           m.ID,
		UserID:       m.UserID,
		MetricType:   m.MetricType,
		TargetValue:  m.TargetValue,
		CurrentValue: m.CurrentValue,
		Deadline:     m.Deadline,
		Status:       m.Status,
	}
}
