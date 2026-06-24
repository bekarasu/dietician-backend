package orchestration

import (
	"context"

	"dietician.local/services/progress-service/internal/progress/dto/request"
	"dietician.local/services/progress-service/internal/progress/dto/response"
	"dietician.local/services/progress-service/internal/progress/repository"
	"dietician.local/services/progress-service/internal/progress/service"
)

type IProgressOrchestrator interface {
	GetProgress(ctx context.Context, userID string) (*response.ProgressResponse, error)
	AddWeight(ctx context.Context, userID string, req request.AddWeightRequest) (*repository.WeightLog, error)
	GetWeeklySummary(ctx context.Context, userID string) (*repository.WeeklyProgressSummary, error)
	AddHabit(ctx context.Context, userID string, req request.AddHabitRequest) (*repository.HabitLog, error)
}

type progressOrchestrator struct {
	progressService service.IProgressService
}

func NewProgressOrchestrator(progressService service.IProgressService) IProgressOrchestrator {
	return &progressOrchestrator{progressService: progressService}
}

func (o *progressOrchestrator) GetProgress(ctx context.Context, userID string) (*response.ProgressResponse, error) {
	return o.progressService.GetProgress(ctx, userID)
}

func (o *progressOrchestrator) AddWeight(ctx context.Context, userID string, req request.AddWeightRequest) (*repository.WeightLog, error) {
	return o.progressService.AddWeight(ctx, userID, req)
}

func (o *progressOrchestrator) GetWeeklySummary(ctx context.Context, userID string) (*repository.WeeklyProgressSummary, error) {
	return o.progressService.GetWeeklySummary(ctx, userID)
}

func (o *progressOrchestrator) AddHabit(ctx context.Context, userID string, req request.AddHabitRequest) (*repository.HabitLog, error) {
	return o.progressService.AddHabit(ctx, userID, req)
}
