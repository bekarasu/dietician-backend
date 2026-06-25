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

// @Summary Create Diet Recommendations
// @Description Generate a list of diet recommendations using OpenAI
// @Tags Recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body openai.DietRecommendationRequestDto true "Recommendation Request"
// @Success 200 {object} response.SuccessResponse{data=openai.DietRecommendationResponse} "OK"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/recommendations [post]
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
