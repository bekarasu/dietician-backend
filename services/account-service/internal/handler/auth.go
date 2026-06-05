package handler

import (
	"encoding/json"
	"net/http"

	pkgresponse "dietician.local/packages/response"
	"dietician.local/services/account-service/internal/auth/model"
	"dietician.local/services/account-service/internal/auth/service"
	"dietician.local/services/account-service/internal/dto/request"
	"dietician.local/services/account-service/internal/dto/response"
	"dietician.local/packages/validation"
)

type IAuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	VerifyOTP(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Me(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	authService service.IUserService
}

func NewAuthHandler(authService service.IUserService) IAuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, validation.FormatValidationError(r.Context(), err))
		return
	}

	svcReq := model.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	resp, err := h.authService.Register(r.Context(), svcReq)
	if err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	dtoResp := response.RegisterResponse{
		OTPToken: resp.OTPToken,
	}

	pkgresponse.JSON(w, http.StatusAccepted, pkgresponse.SuccessResponse{Message: "OTP sent to email", Data: dtoResp})
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req request.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, validation.FormatValidationError(r.Context(), err))
		return
	}

	svcReq := model.VerifyOTPRequest{
		OTPToken: req.OTPToken,
		OTP:      req.OTP,
	}

	resp, err := h.authService.VerifyOTP(r.Context(), svcReq)
	if err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	pkgresponse.Success(w, "OTP verified successfully", dtoResp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, validation.FormatValidationError(r.Context(), err))
		return
	}

	svcReq := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.authService.Login(r.Context(), svcReq)
	if err != nil {
		pkgresponse.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	pkgresponse.Success(w, "login successful", dtoResp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req request.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		pkgresponse.Error(w, http.StatusBadRequest, validation.FormatValidationError(r.Context(), err))
		return
	}

	svcReq := model.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}

	resp, err := h.authService.Refresh(r.Context(), svcReq)
	if err != nil {
		pkgresponse.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	pkgresponse.Success(w, "token refreshed", dtoResp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		pkgresponse.Error(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	user, err := h.authService.GetMe(r.Context(), userID)
	if err != nil {
		pkgresponse.Error(w, http.StatusNotFound, err.Error())
		return
	}

	dtoUser := &response.User{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	pkgresponse.Success(w, "user fetched", dtoUser)
}
