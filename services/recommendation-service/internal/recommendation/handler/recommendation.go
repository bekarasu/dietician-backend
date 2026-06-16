package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/openai"
	"dietician.local/packages/response"
	"dietician.local/services/recommendation-service/internal/recommendation/service"
)

type IRecommendationHandler interface {
	CreateDietRecommendations(c *fiber.Ctx) error
}

type recommendationHandler struct {
	recoService service.IRecommendationService
}

func NewRecommendationHandler(recoService service.IRecommendationService) IRecommendationHandler {
	return &recommendationHandler{
		recoService: recoService,
	}
}

func (h *recommendationHandler) CreateDietRecommendations(c *fiber.Ctx) error {
	var req openai.DietRecommendationRequestDto

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.recoService.CreateDietRecommendations(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to generate recommendations")
	}

	return response.Success(c, "recommendations generated successfully", res)
}
