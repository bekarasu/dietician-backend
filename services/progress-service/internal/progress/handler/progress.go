package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
	"dietician.local/packages/response"
	"dietician.local/packages/utils"
	requestdto "dietician.local/services/progress-service/internal/progress/dto/request"
	_ "dietician.local/services/progress-service/internal/progress/dto/response"
	"dietician.local/services/progress-service/internal/progress/orchestration"
)

type IProgressHandler interface {
	GetProgress(c *fiber.Ctx) error
	AddWeight(c *fiber.Ctx) error
	GetWeeklySummary(c *fiber.Ctx) error
	AddHabit(c *fiber.Ctx) error
}

type progressHandler struct {
	progressOrchestrator orchestration.IProgressOrchestrator
}

func NewProgressHandler(progressOrchestrator orchestration.IProgressOrchestrator) IProgressHandler {
	return &progressHandler{
		progressOrchestrator: progressOrchestrator,
	}
}

// GetProgress retrieves the progress history for a user.
// @Summary Get Progress History
// @Description Retrieve weight logs and habit logs for a specific user
// @Tags Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {object} response.SuccessResponse{data=response.ProgressResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/progress/{userId} [get]
func (h *progressHandler) GetProgress(c *fiber.Ctx) error {
	userID := c.Params("userId")

	progress, err := h.progressOrchestrator.GetProgress(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, utils.TranslateByIDWithContext(c.UserContext(), constants.FailedToGetProgress))
	}

	return response.Success(c, "progress retrieved successfully", progress)
}

// AddWeight creates a new weight log entry.
// @Summary Add Weight Log
// @Description Add a new weight log entry for a specific user
// @Tags Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param request body requestdto.AddWeightRequest true "Add Weight Request"
// @Success 200 {object} response.SuccessResponse{data=repository.WeightLog}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/progress/{userId}/weight [post]
func (h *progressHandler) AddWeight(c *fiber.Ctx) error {
	userID := c.Params("userId")
	var req requestdto.AddWeightRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	log, err := h.progressOrchestrator.AddWeight(c.Context(), userID, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, "weight log added successfully", log)
}

// GetWeeklySummary retrieves the weekly progress summary for a user.
// @Summary Get Weekly Summary
// @Description Retrieve the weekly progress summary including weight changes and habit completion
// @Tags Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Success 200 {object} response.SuccessResponse{data=repository.WeeklyProgressSummary}
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/progress/{userId}/weekly-summary [get]
func (h *progressHandler) GetWeeklySummary(c *fiber.Ctx) error {
	userID := c.Params("userId")

	summary, err := h.progressOrchestrator.GetWeeklySummary(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, utils.TranslateByIDWithContext(c.UserContext(), constants.FailedToGetWeeklySummary))
	}

	return response.Success(c, "weekly summary retrieved successfully", summary)
}

// AddHabit creates a new habit log entry.
// @Summary Add Habit Log
// @Description Add a new habit log entry for a specific user
// @Tags Progress
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Param request body requestdto.AddHabitRequest true "Add Habit Request"
// @Success 200 {object} response.SuccessResponse{data=repository.HabitLog}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/progress/{userId}/habits [post]
func (h *progressHandler) AddHabit(c *fiber.Ctx) error {
	userID := c.Params("userId")
	var req requestdto.AddHabitRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	log, err := h.progressOrchestrator.AddHabit(c.Context(), userID, req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, "habit log added successfully", log)
}
