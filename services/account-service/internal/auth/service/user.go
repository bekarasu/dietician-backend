package service

import (
	"context"
	"errors"
	"time"

	"dietician.local/packages/tokenizer"
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

// TODO locale

func (s *UserService) Register(ctx context.Context, req model.RegisterRequest) (*model.Register, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

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
	if req.OTPToken == "" || req.OTP == "" {
		return nil, errors.New("otp token and otp are required")
	}

	tok, err := s.tokenizer.VerifyJWT(req.OTPToken)
	if err != nil {
		return nil, errors.New("invalid otp token")
	}

	claims, err := s.tokenizer.ExtractClaims(tok)
	if err != nil {
		return nil, errors.New("invalid otp token claims")
	}

	email, ok := claims["email"].(string)
	if !ok || claims["type"] != "otp" {
		return nil, errors.New("invalid otp token type or missing email")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if user.IsVerified {
		return nil, errors.New("user already verified")
	}

	otpRecord, err := s.otpSvc.GetLatestOTPByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if otpRecord == nil || otpRecord.OTPCode != req.OTP {
		return nil, errors.New("invalid OTP")
	}
	if time.Now().After(otpRecord.ExpiresAt) {
		return nil, errors.New("OTP expired")
	}

	if err := s.userRepo.MarkVerified(ctx, user.ID); err != nil {
		return nil, err
	}
	s.otpSvc.DeleteOTPByEmail(ctx, email)

	return s.generateTokens(ctx, user)
}

func (s *UserService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokens(ctx, user)
}

func (s *UserService) Refresh(ctx context.Context, req model.RefreshRequest) (*model.AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	rt, err := s.refreshTokenSvc.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if rt == nil {
		return nil, errors.New("invalid refresh token")
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
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
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) generateTokens(ctx context.Context, user *model.User) (*model.AuthResponse, error) {
	accessStr, err := s.refreshTokenSvc.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
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
