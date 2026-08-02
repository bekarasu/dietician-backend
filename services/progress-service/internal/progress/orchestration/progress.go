package orchestration

import (

	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"dietician.local/packages/httpclient"
	"dietician.local/services/progress-service/config"
	"dietician.local/services/progress-service/internal/progress/dto/request"
	"dietician.local/services/progress-service/internal/progress/dto/response"
	"dietician.local/services/progress-service/internal/progress/repository"
	"dietician.local/services/progress-service/internal/progress/service"
)

type IProgressOrchestrator interface {
	GetProgress(ctx context.Context, userID string) (*response.ProgressResponse, error)
	AddWeight(ctx context.Context, userID string, req request.AddWeightRequest, token string) (*repository.WeightLog, error)
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

func (o *progressOrchestrator) AddWeight(ctx context.Context, userID string, req request.AddWeightRequest, token string) (*repository.WeightLog, error) {
	log, err := o.progressService.AddWeight(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// Sync weight with account-service
	go func() {
		client := httpclient.Default()
		url := fmt.Sprintf("%s/api/v1/profiles/weight", config.ProgressApp.Web.AccountServiceURL)
		
		payload, _ := json.Marshal(map[string]float64{"weightKg": req.WeightKg})
		httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/json")
			if token != "" {
				httpReq.Header.Set("Authorization", "Bearer "+token)
			}
			client.Do(httpReq)
		}
	}()

	return log, nil
}

func (o *progressOrchestrator) GetWeeklySummary(ctx context.Context, userID string) (*repository.WeeklyProgressSummary, error) {
	return o.progressService.GetWeeklySummary(ctx, userID)
}

func (o *progressOrchestrator) AddHabit(ctx context.Context, userID string, req request.AddHabitRequest) (*repository.HabitLog, error) {
	return o.progressService.AddHabit(ctx, userID, req)
}
