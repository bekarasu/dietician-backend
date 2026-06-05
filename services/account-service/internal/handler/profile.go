package handler

import (
	"encoding/json"
	"net/http"

	pkgresponse "dietician.local/packages/response"
	"dietician.local/packages/validation"
	"dietician.local/services/account-service/internal/dto/request"
	profileservice "dietician.local/services/account-service/internal/profile/service"
)

type IProfileHandler interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
	UpsertProfile(w http.ResponseWriter, r *http.Request)
	GetPreferences(w http.ResponseWriter, r *http.Request)
	UpdatePreferences(w http.ResponseWriter, r *http.Request)
}

type ProfileHandler struct {
	profileService *profileservice.ProfileService
}

func NewProfileHandler(profileService *profileservice.ProfileService) IProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		pkgresponse.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		pkgresponse.Error(w, http.StatusNotFound, err.Error())
		return
	}

	pkgresponse.Success(w, "profile fetched", profile.ToResource())
}

func (h *ProfileHandler) UpsertProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		pkgresponse.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	var req request.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, validation.FormatValidationError(r.Context(), err))
		return
	}

	profile, err := h.profileService.UpsertProfile(r.Context(), userID, req.ToDomain())
	if err != nil {
		pkgresponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkgresponse.Success(w, "profile updated", profile.ToResource())
}

func (h *ProfileHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		pkgresponse.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	prefs, err := h.profileService.GetPreferences(r.Context(), userID)
	if err != nil {
		pkgresponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkgresponse.Success(w, "preferences fetched", prefs.ToResource())
}

func (h *ProfileHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		pkgresponse.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	var req request.UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, validation.FormatValidationError(r.Context(), err))
		return
	}

	prefs, err := h.profileService.UpdatePreferences(r.Context(), userID, req.ToDomain())
	if err != nil {
		pkgresponse.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkgresponse.Success(w, "preferences updated", prefs.ToResource())
}
