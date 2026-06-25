package service

import (
	"context"
	"errors"

	"dietician.local/packages/constants"
	"dietician.local/packages/utils"
	"dietician.local/services/progress-service/internal/progress/dto/request"
	"dietician.local/services/progress-service/internal/progress/dto/response"
	"dietician.local/services/progress-service/internal/progress/repository"
	
	dailylog_service "dietician.local/services/progress-service/internal/dailylog/service"
	tracking_service "dietician.local/services/progress-service/internal/tracking/service"
)

type IProgressService interface {
	GetProgress(ctx context.Context, userID string) (*response.ProgressResponse, error)
	AddWeight(ctx context.Context, userID string, req request.AddWeightRequest) (*repository.WeightLog, error)
	GetWeeklySummary(ctx context.Context, userID string) (*repository.WeeklyProgressSummary, error)
	AddHabit(ctx context.Context, userID string, req request.AddHabitRequest) (*repository.HabitLog, error)
}

type progressService struct {
	repo         repository.IProgressRepository
	dlService    dailylog_service.IDailyLogService
	trackService tracking_service.ITrackingService
}

func NewProgressService(repo repository.IProgressRepository, dlService dailylog_service.IDailyLogService, trackService tracking_service.ITrackingService) IProgressService {
	return &progressService{
		repo:         repo,
		dlService:    dlService,
		trackService: trackService,
	}
}

func (s *progressService) GetProgress(ctx context.Context, userID string) (*response.ProgressResponse, error) {
	weights, err := s.repo.GetWeightLogs(ctx, userID)
	if err != nil {
		return nil, err
	}
	habits, err := s.repo.GetHabitLogs(ctx, userID)
	if err != nil {
		return nil, err
	}

	if weights == nil {
		weights = []repository.WeightLog{}
	}
	if habits == nil {
		habits = []repository.HabitLog{}
	}

	dailyLogs, err := s.dlService.GetDailyLogs(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	trackingMetrics, err := s.trackService.GetTrackingMetrics(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &response.ProgressResponse{
		WeightLogs:      weights,
		HabitLogs:       habits,
		DailyLogs:       dailyLogs,
		TrackingMetrics: trackingMetrics,
	}, nil
}

func (s *progressService) AddWeight(ctx context.Context, userID string, req request.AddWeightRequest) (*repository.WeightLog, error) {
	if req.WeightKg <= 0 {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.WeightMustBePositive))
	}

	log := &repository.WeightLog{
		UserID:   userID,
		WeightKg: req.WeightKg,
		Notes:    req.Notes,
	}
	if err := s.repo.AddWeightLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *progressService) GetWeeklySummary(ctx context.Context, userID string) (*repository.WeeklyProgressSummary, error) {
	return s.repo.GetWeeklySummary(ctx, userID)
}

func (s *progressService) AddHabit(ctx context.Context, userID string, req request.AddHabitRequest) (*repository.HabitLog, error) {
	if req.HabitName == "" {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.HabitNameRequired))
	}

	log := &repository.HabitLog{
		UserID:    userID,
		HabitName: req.HabitName,
		Completed: req.Completed,
		Notes:     req.Notes,
	}
	if err := s.repo.AddHabitLog(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}
