package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
	"dietician.local/packages/response"
	"dietician.local/packages/utils"
	"dietician.local/services/progress-service/internal/tracking/dto"
	"dietician.local/services/progress-service/internal/tracking/orchestration"
)

type ITrackingHandler interface {
	CreateTrackingMetric(c *fiber.Ctx) error
	GetTrackingMetrics(c *fiber.Ctx) error
}

type trackingHandler struct {
	orchestrator orchestration.ITrackingOrchestrator
}

func NewTrackingHandler(orchestrator orchestration.ITrackingOrchestrator) ITrackingHandler {
	return &trackingHandler{orchestrator: orchestrator}
}

func (h *trackingHandler) CreateTrackingMetric(c *fiber.Ctx) error {
	var req dto.CreateTrackingMetricRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	res, err := h.orchestrator.CreateTrackingMetric(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, utils.TranslateByIDWithContext(c.UserContext(), constants.FailedToCreateTrackingMetric))
	}

	return response.Success(c, "tracking metric created", res)
}

func (h *trackingHandler) GetTrackingMetrics(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	res, err := h.orchestrator.GetTrackingMetrics(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, utils.TranslateByIDWithContext(c.UserContext(), constants.FailedToGetTrackingMetrics))
	}

	return response.Success(c, "tracking metrics retrieved", res)
}
