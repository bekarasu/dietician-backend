package orchestration

import (
	"context"

	"dietician.local/services/recommendation-service/internal/dietplan/dto"
	"dietician.local/services/recommendation-service/internal/dietplan/service"
)

type IDietPlanOrchestrator interface {
	CreateDietPlan(ctx context.Context, req dto.CreateDietPlanRequest) (*dto.DietPlanResponse, error)
}

type dietPlanOrchestrator struct {
	dietPlanService service.IDietPlanService
}

func NewDietPlanOrchestrator(dietPlanService service.IDietPlanService) IDietPlanOrchestrator {
	return &dietPlanOrchestrator{
		dietPlanService: dietPlanService,
	}
}

func (o *dietPlanOrchestrator) CreateDietPlan(ctx context.Context, req dto.CreateDietPlanRequest) (*dto.DietPlanResponse, error) {
	// You can add validation, external API calls, or other complex workflows here
	return o.dietPlanService.CreateDietPlan(ctx, req)
}
