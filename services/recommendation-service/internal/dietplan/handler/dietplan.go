package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
	"dietician.local/packages/response"
	"dietician.local/packages/utils"
	"dietician.local/services/recommendation-service/internal/dietplan/dto"
	"dietician.local/services/recommendation-service/internal/dietplan/orchestration"
)

type IDietPlanHandler interface {
	CreateDietPlan(c *fiber.Ctx) error
}

type dietPlanHandler struct {
	orchestrator orchestration.IDietPlanOrchestrator
}

func NewDietPlanHandler(orchestrator orchestration.IDietPlanOrchestrator) IDietPlanHandler {
	return &dietPlanHandler{
		orchestrator: orchestrator,
	}
}

// @Summary Create Diet Plan
// @Description Create a new diet plan based on given parameters
// @Tags Diet Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateDietPlanRequest true "Diet Plan Request"
// @Success 200 {object} response.SuccessResponse{data=dto.DietPlanResponse} "OK"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/diet-plans [post]
func (h *dietPlanHandler) CreateDietPlan(c *fiber.Ctx) error {
	var req dto.CreateDietPlanRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	res, err := h.orchestrator.CreateDietPlan(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, utils.TranslateByIDWithContext(c.UserContext(), constants.FailedToCreateDietPlan))
	}

	return response.Success(c, "diet plan created successfully", res)
}
