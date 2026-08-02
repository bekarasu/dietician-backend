package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
	"dietician.local/packages/middleware"
	pkgresponse "dietician.local/packages/response"
	"dietician.local/packages/utils"
	"dietician.local/packages/validation"
	"dietician.local/services/account-service/internal/profile/dto/request"
	"dietician.local/services/account-service/internal/profile/orchestration"
)

type IProfileHandler interface {
	GetProfile(c *fiber.Ctx) error
	UpsertProfile(c *fiber.Ctx) error
	UpdateWeight(c *fiber.Ctx) error
	Onboarding(c *fiber.Ctx) error
	GetPreferences(c *fiber.Ctx) error
	UpdatePreferences(c *fiber.Ctx) error
}

type ProfileHandler struct {
	profileOrchestrator orchestration.IProfileOrchestrator
}

func NewProfileHandler(profileOrchestrator orchestration.IProfileOrchestrator) IProfileHandler {
	return &ProfileHandler{profileOrchestrator: profileOrchestrator}
}

// GetProfile retrieves the authenticated user's profile.
// @Summary Get User Profile
// @Description Retrieve the profile of the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.UserProfile}
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Failure 404 {object} pkgresponse.ErrorResponse
// @Router /api/v1/profiles [get]
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	profile, err := h.profileOrchestrator.GetProfile(c.UserContext(), userID)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	return pkgresponse.Success(c, "profile fetched", profile.ToResource())
}

// UpsertProfile updates or creates the user's profile.
// @Summary Update User Profile
// @Description Update or create the profile for the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.UserProfile}
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Failure 500 {object} pkgresponse.ErrorResponse
// @Router /api/v1/profiles [put]
func (h *ProfileHandler) UpsertProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	var req request.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	profile, err := h.profileOrchestrator.UpsertProfile(c.UserContext(), userID, req.ToDomain())
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "profile updated", profile.ToResource())
}

// UpdateWeight updates the user's weight (intended for internal/service-to-service use).
// @Summary Update User Weight
// @Description Update only the weight for the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateWeightRequest true "Update Weight Request"
// @Success 200 {object} pkgresponse.SuccessResponse
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Failure 500 {object} pkgresponse.ErrorResponse
// @Router /api/v1/profiles/weight [put]
func (h *ProfileHandler) UpdateWeight(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	var req request.UpdateWeightRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	err := h.profileOrchestrator.UpdateWeight(c.UserContext(), userID, &model.UpdateWeightRequest{WeightKg: req.WeightKg})
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "weight updated successfully", nil)
}

// Onboarding handles the comprehensive onboarding endpoint after registration.
// @Summary Onboarding Profile Setup
// @Description Sets up the user's profile and preferences in one step
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.OnboardingRequest true "Onboarding Request"
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.UserProfile}
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Failure 500 {object} pkgresponse.ErrorResponse
// @Router /api/v1/profiles/onboarding [post]
func (h *ProfileHandler) Onboarding(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())
	if userID == "" {
		return pkgresponse.Error(c, fiber.StatusUnauthorized, utils.TranslateByIDWithContext(c.UserContext(), constants.UserIDRequired))
	}

	var req request.OnboardingRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	profile, err := h.profileOrchestrator.Onboarding(c.UserContext(), userID, req.ToDomain())
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "onboarding completed", profile.ToResource())
}

// GetPreferences retrieves the user's preferences.
// @Summary Get User Preferences
// @Description Retrieve the preferences of the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.PreferencesResponse}
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Failure 500 {object} pkgresponse.ErrorResponse
// @Router /api/v1/profiles/preferences [get]
func (h *ProfileHandler) GetPreferences(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	prefs, err := h.profileOrchestrator.GetPreferences(c.UserContext(), userID)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "preferences fetched", prefs.ToResource())
}

// UpdatePreferences updates the user's preferences.
// @Summary Update User Preferences
// @Description Update the preferences for the authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdatePreferencesRequest true "Update Preferences Request"
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.PreferencesResponse}
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Failure 500 {object} pkgresponse.ErrorResponse
// @Router /api/v1/profiles/preferences [put]
func (h *ProfileHandler) UpdatePreferences(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())
	if userID == "" {
		return pkgresponse.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.UserIDRequired))
	}

	var req request.UpdatePreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, utils.TranslateByIDWithContext(c.UserContext(), constants.InvalidRequestBody))
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	prefs, err := h.profileOrchestrator.UpdatePreferences(c.UserContext(), userID, req.ToDomain())
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "preferences updated", prefs.ToResource())
}
