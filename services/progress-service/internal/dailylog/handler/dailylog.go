package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
	"dietician.local/packages/response"
	"dietician.local/packages/utils"
	"dietician.local/services/progress-service/internal/dailylog/dto"
	"dietician.local/services/progress-service/internal/dailylog/orchestration"
)

type IDailyLogHandler interface {
	CreateDailyLog(c *fiber.Ctx) error
	GetDailyLog(c *fiber.Ctx) error
}

type dailyLogHandler struct {
	orchestrator orchestration.IDailyLogOrchestrator
}

func NewDailyLogHandler(orchestrator orchestration.IDailyLogOrchestrator) IDailyLogHandler {
	return &dailyLogHandler{
		orchestrator: orchestrator,
	}
}

func (h *dailyLogHandler) CreateDailyLog(c *fiber.Ctx) error {
	var req dto.CreateDailyLogRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	res, err := h.orchestrator.CreateDailyLog(c.Context(), req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, utils.TranslateByIDWithContext(c.UserContext(), constants.FailedToCreateDailyLog))
	}
	return response.Success(c, "daily log created", res)
}

func (h *dailyLogHandler) GetDailyLog(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	date := c.Params("date")

	res, err := h.orchestrator.GetDailyLog(c.Context(), userID, date)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, utils.TranslateByIDWithContext(c.UserContext(), constants.DailyLogNotFound))
	}
	return response.Success(c, "daily log retrieved", res)
}
