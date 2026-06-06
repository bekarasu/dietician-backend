package handler

import (
	"github.com/gofiber/fiber/v2"

	pkgresponse "dietician.local/packages/response"
	"dietician.local/packages/validation"
	"dietician.local/services/account-service/internal/auth/model"
	"dietician.local/services/account-service/internal/auth/service"
	"dietician.local/services/account-service/internal/dto/request"
	"dietician.local/services/account-service/internal/dto/response"
)

type IAuthHandler interface {
	Register(c *fiber.Ctx) error
	VerifyOTP(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type AuthHandler struct {
	authService service.IUserService
}

func NewAuthHandler(authService service.IUserService) IAuthHandler {
	return &AuthHandler{authService: authService}
}

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

	resp, err := h.authService.Register(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dtoResp := response.RegisterResponse{
		OTPToken: resp.OTPToken,
	}

	return pkgresponse.JSON(c, fiber.StatusAccepted, pkgresponse.SuccessResponse{Message: "OTP sent to email", Data: dtoResp})
}

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

	resp, err := h.authService.VerifyOTP(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	return pkgresponse.Success(c, "OTP verified successfully", dtoResp)
}

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

	resp, err := h.authService.Login(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	return pkgresponse.Success(c, "login successful", dtoResp)
}

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

	resp, err := h.authService.Refresh(c.UserContext(), svcReq)
	if err != nil {
		return pkgresponse.Error(c, fiber.StatusUnauthorized, err.Error())
	}

	dtoResp := response.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}

	return pkgresponse.Success(c, "token refreshed", dtoResp)
}
