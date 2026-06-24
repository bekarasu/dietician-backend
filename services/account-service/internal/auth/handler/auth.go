package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/middleware"
	pkgresponse "dietician.local/packages/response"
	"dietician.local/packages/validation"
	"dietician.local/services/account-service/internal/auth/model"
	"dietician.local/services/account-service/internal/auth/orchestration"
	"dietician.local/services/account-service/internal/auth/dto/request"
	"dietician.local/services/account-service/internal/auth/dto/response"
)

type IAuthHandler interface {
	Register(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type AuthHandler struct {
	authOrchestrator orchestration.IAuthOrchestrator
}

func NewAuthHandler(authOrchestrator orchestration.IAuthOrchestrator) IAuthHandler {
	return &AuthHandler{authOrchestrator: authOrchestrator}
}

// Register creates a new user account and sends an OTP.
// @Summary Register a new user
// @Description Register a new user and send an OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Register Request"
// @Success 202 {object} pkgresponse.SuccessResponse{data=response.RegisterResponse}
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req request.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	svcReq := model.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	resp, err := h.authOrchestrator.Register(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dtoResp := response.RegisterResponse{
		OTPToken: resp.OTPToken,
	}

	return pkgresponse.JSON(c, fiber.StatusAccepted, pkgresponse.SuccessResponse{Message: "OTP sent to email", Data: dtoResp})
}

// VerifyOTP verifies the given OTP for a user.
// @Summary Verify OTP
// @Description Verify OTP for a registered user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.VerifyOTPRequest true "Verify OTP Request"
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.AuthResponse}
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Router /api/v1/auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req request.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	svcReq := model.VerifyOTPRequest{
		OTPToken: req.OTPToken,
		OTP:      req.OTP,
	}

	resp, err := h.authOrchestrator.VerifyOTP(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	return pkgresponse.Success(c, "OTP verified successfully", dtoResp)
}

// Login authenticates a user and returns tokens.
// @Summary Login
// @Description Authenticate user and return tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Login Request"
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.AuthResponse}
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	svcReq := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.authOrchestrator.Login(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	return pkgresponse.Success(c, "login successful", dtoResp)
}

// Refresh generates new tokens from a refresh token.
// @Summary Refresh Token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.RefreshRequest true "Refresh Request"
// @Success 200 {object} pkgresponse.SuccessResponse{data=response.AuthResponse}
// @Failure 400 {object} pkgresponse.ErrorResponse
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var req request.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	if err := validation.Validate.Struct(req); err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, validation.FormatValidationError(c.UserContext(), err))
	}

	svcReq := model.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}

	resp, err := h.authOrchestrator.Refresh(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	return pkgresponse.Success(c, "token refreshed", dtoResp)
}

// Logout authenticates a user and invalidates their tokens.
// @Summary Logout
// @Description Logout user and invalidate access/refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} pkgresponse.SuccessResponse
// @Failure 401 {object} pkgresponse.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())
	if userID == "" {
		return pkgresponse.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	token := middleware.ExtractBearerToken(c)

	err := h.authOrchestrator.Logout(c.UserContext(), userID, token)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return pkgresponse.Success(c, "logout successful 3", nil)
}
