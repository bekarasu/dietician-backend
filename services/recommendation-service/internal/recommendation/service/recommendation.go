package service

import (
	"context"

	"dietician.local/packages/openai"
)

type IRecommendationService interface {
	CreateDietRecommendations(ctx context.Context, req openai.DietRecommendationRequestDto) (*openai.DietRecommendationResponse, error)
}

type recommendationService struct {
	openaiService openai.Service
}

func NewRecommendationService(openaiService openai.Service) IRecommendationService {
	return &recommendationService{
		openaiService: openaiService,
	}
}

func (s *recommendationService) CreateDietRecommendations(ctx context.Context, req openai.DietRecommendationRequestDto) (*openai.DietRecommendationResponse, error) {
	// Directly returning from openai service since we're not saving to DB
	return s.openaiService.CreateDietRecommendations(ctx, req)
}
