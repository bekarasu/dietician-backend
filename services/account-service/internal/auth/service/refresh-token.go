package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"dietician.local/packages/tokenizer"
	"dietician.local/services/account-service/config"
	"dietician.local/services/account-service/internal/auth/model"
	"dietician.local/services/account-service/internal/auth/repository"
	"github.com/golang-jwt/jwt/v5"
)

// IRefreshTokenService combines JWT token generation/validation with refresh-token persistence.
type IRefreshTokenService interface {
	GenerateAccessToken(userID, email string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(tokenString string) (string, error)
	RefreshTTL() time.Duration

	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteRefreshTokensByUserID(ctx context.Context, userID string) error
}

// RefreshTokenService is the concrete implementation of IRefreshTokenService.
type RefreshTokenService struct {
	tokenizer  tokenizer.ITokenizer
	accessTTL  time.Duration
	refreshTTL time.Duration
	repo       repository.IRefreshTokenRepo
}

func NewRefreshTokenService(
	cfg *config.AccountAppScheme,
	repo repository.IRefreshTokenRepo,
	tok tokenizer.ITokenizer,
) IRefreshTokenService {
	return &RefreshTokenService{
		tokenizer:  tok,
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
		repo:       repo,
	}
}

func (s *RefreshTokenService) AccessTTL() time.Duration  { return s.accessTTL }
func (s *RefreshTokenService) RefreshTTL() time.Duration { return s.refreshTTL }

// GenerateAccessToken creates a signed HS256 JWT with sub and email claims.
func (s *RefreshTokenService) GenerateAccessToken(userID, email string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iat":   now.Unix(),
		"exp":   now.Add(s.accessTTL).Unix(),
	}
	return s.tokenizer.GenerateJWT(claims)
}

// GenerateRefreshToken creates a cryptographically random opaque token.
func (s *RefreshTokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ValidateAccessToken parses the JWT and returns the subject (userID) on success.
func (s *RefreshTokenService) ValidateAccessToken(tokenString string) (string, error) {
	t, err := s.tokenizer.VerifyJWT(tokenString)
	if err != nil {
		return "", err
	}

	claims, err := s.tokenizer.ExtractClaims(t)
	if err != nil {
		return "", err
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	return userID, nil
}

func (s *RefreshTokenService) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	return s.repo.Create(ctx, token)
}

func (s *RefreshTokenService) GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	return s.repo.GetByToken(ctx, token)
}

func (s *RefreshTokenService) DeleteRefreshTokensByUserID(ctx context.Context, userID string) error {
	return s.repo.DeleteByUserID(ctx, userID)
}
