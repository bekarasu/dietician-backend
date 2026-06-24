package orchestration

import (
	"context"

	"dietician.local/packages/openai"
	"dietician.local/services/recommendation-service/internal/recommendation/service"
)

type IRecommendationOrchestrator interface {
	CreateDietRecommendations(ctx context.Context, req openai.DietRecommendationRequestDto) (*openai.DietRecommendationResponse, error)
}

type recommendationOrchestrator struct {
	recoService service.IRecommendationService
}

func NewRecommendationOrchestrator(recoService service.IRecommendationService) IRecommendationOrchestrator {
	return &recommendationOrchestrator{recoService: recoService}
}

func (o *recommendationOrchestrator) CreateDietRecommendations(ctx context.Context, req openai.DietRecommendationRequestDto) (*openai.DietRecommendationResponse, error) {
	return o.recoService.CreateDietRecommendations(ctx, req)
}
