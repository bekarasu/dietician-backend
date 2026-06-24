package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/middleware"
	pkgresponse "dietician.local/packages/response"
	"dietician.local/packages/validation"
	"dietician.local/services/account-service/internal/profile/dto/request"
	profileservice "dietician.local/services/account-service/internal/profile/service"
)

type IProfileHandler interface {
	GetProfile(c *fiber.Ctx) error
	UpsertProfile(c *fiber.Ctx) error
	GetPreferences(c *fiber.Ctx) error
	UpdatePreferences(c *fiber.Ctx) error
}

type ProfileHandler struct {
	profileService *profileservice.ProfileService
}

func NewProfileHandler(profileService *profileservice.ProfileService) IProfileHandler {
	return &ProfileHandler{profileService: profileService}
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

	profile, err := h.profileService.GetProfile(c.UserContext(), userID)
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
		return pkgresponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	profile, err := h.profileService.UpsertProfile(c.UserContext(), userID, req.ToDomain())
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "profile updated", profile.ToResource())
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

	prefs, err := h.profileService.GetPreferences(c.UserContext(), userID)
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
		return pkgresponse.Error(c, fiber.StatusBadRequest, "user ID is required")
	}

	var req request.UpdatePreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	prefs, err := h.profileService.UpdatePreferences(c.UserContext(), userID, req.ToDomain())
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "preferences updated", prefs.ToResource())
}
