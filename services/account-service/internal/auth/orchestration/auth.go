package orchestration

import (
	"context"

	"dietician.local/services/account-service/internal/auth/model"
	"dietician.local/services/account-service/internal/auth/service"
)

type IAuthOrchestrator interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.Register, error)
	VerifyOTP(ctx context.Context, req model.VerifyOTPRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error)
	Refresh(ctx context.Context, req model.RefreshRequest) (*model.AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*model.User, error)
	Logout(ctx context.Context, userID string, accessToken string) error
}

type authOrchestrator struct {
	authService service.IUserService
}

func NewAuthOrchestrator(authService service.IUserService) IAuthOrchestrator {
	return &authOrchestrator{authService: authService}
}

func (o *authOrchestrator) Register(ctx context.Context, req model.RegisterRequest) (*model.Register, error) {
	return o.authService.Register(ctx, req)
}

func (o *authOrchestrator) VerifyOTP(ctx context.Context, req model.VerifyOTPRequest) (*model.AuthResponse, error) {
	return o.authService.VerifyOTP(ctx, req)
}

func (o *authOrchestrator) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	return o.authService.Login(ctx, req)
}

func (o *authOrchestrator) Refresh(ctx context.Context, req model.RefreshRequest) (*model.AuthResponse, error) {
	return o.authService.Refresh(ctx, req)
}

func (o *authOrchestrator) GetMe(ctx context.Context, userID string) (*model.User, error) {
	return o.authService.GetMe(ctx, userID)
}

func (o *authOrchestrator) Logout(ctx context.Context, userID string, accessToken string) error {
	return o.authService.Logout(ctx, userID, accessToken)
}
