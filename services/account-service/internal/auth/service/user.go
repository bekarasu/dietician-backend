package service

import (
	"context"
	"errors"
	"time"

	"dietician.local/packages/constants"
	"dietician.local/packages/tokenizer"
	"dietician.local/packages/utils"
	"dietician.local/services/account-service/internal/auth/model"
	"dietician.local/services/account-service/internal/auth/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Register(ctx context.Context, req model.RegisterRequest) (*model.Register, error)
	VerifyOTP(ctx context.Context, req model.VerifyOTPRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error)
	Refresh(ctx context.Context, req model.RefreshRequest) (*model.AuthResponse, error)
	GetMe(ctx context.Context, userID string) (*model.User, error)
	Logout(ctx context.Context, userID string, accessToken string) error
}

type UserService struct {
	userRepo        repository.IUserRepository
	refreshTokenSvc IRefreshTokenService
	otpSvc          IOTPService
	tokenizer       tokenizer.ITokenizer
}

func NewUserService(
	userRepo repository.IUserRepository,
	refreshTokenSvc IRefreshTokenService,
	otpSvc IOTPService,
	tok tokenizer.ITokenizer,
) IUserService {
	return &UserService{
		userRepo:        userRepo,
		refreshTokenSvc: refreshTokenSvc,
		otpSvc:          otpSvc,
		tokenizer:       tok,
	}
}

func (s *UserService) Register(ctx context.Context, req model.RegisterRequest) (*model.Register, error) {
	if len(req.Password) < 8 {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.PasswordTooShort))
	}

	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.IsVerified {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.EmailAlreadyRegistered))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		existing.PasswordHash = string(hash)
		existing.FirstName = req.FirstName
		existing.LastName = req.LastName
		if err := s.userRepo.UpdateUnverified(ctx, existing); err != nil {
			return nil, err
		}
	} else {
		user := &model.User{
			Email:        req.Email,
			PasswordHash: string(hash),
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			IsVerified:   false,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	code := s.otpSvc.Generate()
	otpRecord := &model.OTP{
		Email:     req.Email,
		OTPCode:   code,
		ExpiresAt: time.Now().Add(time.Duration(s.otpSvc.ExpireSeconds()) * time.Second),
	}
	if err := s.otpSvc.CreateOTP(ctx, otpRecord); err != nil {
		return nil, err
	}

	go s.otpSvc.SendEmail(req.Email, code) //nolint:errcheck

	otpClaims := jwt.MapClaims{
		"email": req.Email,
		"type":  "otp",
		"exp":   time.Now().Add(time.Duration(s.otpSvc.ExpireSeconds()) * time.Second).Unix(),
	}

	otpToken, err := s.tokenizer.GenerateJWT(otpClaims)
	if err != nil {
		return nil, err
	}

	return &model.Register{
		OTPToken: otpToken,
	}, nil
}

func (s *UserService) VerifyOTP(ctx context.Context, req model.VerifyOTPRequest) (*model.AuthResponse, error) {
	tok, err := s.tokenizer.VerifyJWT(req.OTPToken)
	if err != nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidOTPToken))
	}

	claims, err := s.tokenizer.ExtractClaims(tok)
	if err != nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidOTPTokenClaims))
	}

	email, ok := claims["email"].(string)
	if !ok || claims["type"] != "otp" {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidOTPTokenClaims))
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.AccountNotFound))
	}
	if user.IsVerified {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.EmailAlreadyRegistered))
	}

	otpRecord, err := s.otpSvc.GetLatestOTPByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if otpRecord == nil || otpRecord.OTPCode != req.OTP {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidOTP))
	}
	if time.Now().After(otpRecord.ExpiresAt) {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidOTP))
	}

	if err := s.userRepo.MarkVerified(ctx, user.ID); err != nil {
		return nil, err
	}
	s.otpSvc.DeleteOTPByEmail(ctx, email)

	return s.generateTokens(ctx, user)
}

func (s *UserService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidCredentials))
	}

	if !user.IsVerified {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.EmailNotVerified))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidCredentials))
	}

	return s.generateTokens(ctx, user)
}

func (s *UserService) Refresh(ctx context.Context, req model.RefreshRequest) (*model.AuthResponse, error) {
	rt, err := s.refreshTokenSvc.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidRefreshToken))
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.InvalidRefreshToken))
	}

	user, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.UserNotFound))
	}

	if err := s.refreshTokenSvc.DeleteRefreshTokensByUserID(ctx, user.ID); err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, user)
}

func (s *UserService) GetMe(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(utils.TranslateByIDWithContext(ctx, constants.AccountNotFound))
	}
	return user, nil
}

func (s *UserService) Logout(ctx context.Context, userID string, accessToken string) error {
	if err := s.refreshTokenSvc.DeleteRefreshTokensByUserID(ctx, userID); err != nil {
		return err
	}
	if accessToken != "" {
		if err := s.tokenizer.RemoveJWTFromRedis(ctx, accessToken); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) generateTokens(ctx context.Context, user *model.User) (*model.AuthResponse, error) {
	accessStr, err := s.refreshTokenSvc.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	if err := s.tokenizer.StoreJWTInRedis(ctx, accessStr, s.refreshTokenSvc.AccessTTL()); err != nil {
		return nil, err
	}

	refreshStr, err := s.refreshTokenSvc.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	rt := &model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshStr,
		ExpiresAt: time.Now().Add(s.refreshTokenSvc.RefreshTTL()),
	}
	if err := s.refreshTokenSvc.CreateRefreshToken(ctx, rt); err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
		User:         user,
	}, nil
}
