package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/openai"
	"dietician.local/packages/response"
	"dietician.local/services/recommendation-service/internal/recommendation/orchestration"
)

type IRecommendationHandler interface {
	CreateDietRecommendations(c *fiber.Ctx) error
}

type recommendationHandler struct {
	recoOrchestrator orchestration.IRecommendationOrchestrator
}

func NewRecommendationHandler(recoOrchestrator orchestration.IRecommendationOrchestrator) IRecommendationHandler {
	return &recommendationHandler{
		recoOrchestrator: recoOrchestrator,
	}
}

func (h *recommendationHandler) CreateDietRecommendations(c *fiber.Ctx) error {
	var req openai.DietRecommendationRequestDto

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.recoOrchestrator.CreateDietRecommendations(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to generate recommendations")
	}

	return response.Success(c, "recommendations generated successfully", res)
}
