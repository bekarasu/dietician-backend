package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/response"
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

func (h *dietPlanHandler) CreateDietPlan(c *fiber.Ctx) error {
	var req dto.CreateDietPlanRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.orchestrator.CreateDietPlan(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create diet plan")
	}

	return response.Success(c, "diet plan created successfully", res)
}
