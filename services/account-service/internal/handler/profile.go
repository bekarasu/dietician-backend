package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/middleware"
	pkgresponse "dietician.local/packages/response"
	"dietician.local/packages/validation"
	"dietician.local/services/account-service/internal/dto/request"
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

func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	profile, err := h.profileService.GetProfile(c.UserContext(), userID)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusNotFound, err.Error())
	}

	return pkgresponse.Success(c, "profile fetched", profile.ToResource())
}

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

func (h *ProfileHandler) GetPreferences(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	prefs, err := h.profileService.GetPreferences(c.UserContext(), userID)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "preferences fetched", prefs.ToResource())
}

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
