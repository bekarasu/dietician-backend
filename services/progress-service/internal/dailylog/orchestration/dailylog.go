package orchestration

import (
	"context"

	"dietician.local/services/progress-service/internal/dailylog/dto"
	"dietician.local/services/progress-service/internal/dailylog/service"
)

type IDailyLogOrchestrator interface {
	CreateDailyLog(ctx context.Context, req dto.CreateDailyLogRequest) (*dto.DailyLogResponse, error)
	GetDailyLog(ctx context.Context, userID string, date string) (*dto.DailyLogResponse, error)
}

type dailyLogOrchestrator struct {
	svc service.IDailyLogService
}

func NewDailyLogOrchestrator(svc service.IDailyLogService) IDailyLogOrchestrator {
	return &dailyLogOrchestrator{
		svc: svc,
	}
}

func (o *dailyLogOrchestrator) CreateDailyLog(ctx context.Context, req dto.CreateDailyLogRequest) (*dto.DailyLogResponse, error) {
	return o.svc.CreateDailyLog(ctx, req)
}

func (o *dailyLogOrchestrator) GetDailyLog(ctx context.Context, userID string, date string) (*dto.DailyLogResponse, error) {
	return o.svc.GetDailyLog(ctx, userID, date)
}
