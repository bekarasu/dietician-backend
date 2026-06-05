package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"dietician.local/services/account-service/config"
	"dietician.local/services/account-service/internal/auth/model"
	"dietician.local/services/account-service/internal/auth/repository"
)

// ─── OTP email-sending service ───────────────────────────────────────────────

type IOTPService interface {
	CreateOTP(ctx context.Context, otp *model.OTP) error
	GetLatestOTPByEmail(ctx context.Context, email string) (*model.OTP, error)
	DeleteOTPByEmail(ctx context.Context, email string) error
	ExpireSeconds() int
	SendEmail(to, code string) error
	Generate() string
}

type EmailSender interface {
	Send(to, subject, body string) error
}

type service struct {
	sender        EmailSender
	expireSeconds int
	repo          repository.IOTPRepository
}

func NewOTPService(sender EmailSender, cfg *config.AccountAppScheme, repo repository.IOTPRepository) IOTPService {
	return &service{sender: sender, expireSeconds: cfg.JWT.OTPExpireSec, repo: repo}
}

func (s *service) CreateOTP(ctx context.Context, otp *model.OTP) error {
	return s.repo.Create(ctx, otp)
}

func (s *service) GetLatestOTPByEmail(ctx context.Context, email string) (*model.OTP, error) {
	return s.repo.GetLatestByEmail(ctx, email)
}

func (s *service) DeleteOTPByEmail(ctx context.Context, email string) error {
	return s.repo.DeleteByEmail(ctx, email)
}

func (s *service) ExpireSeconds() int { return s.expireSeconds }

func (s *service) Generate() string {
	max := big.NewInt(1_000_000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "123456"
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *service) SendEmail(to, code string) error {
	// TODO locale
	subject := "Registration OTP"
	// TODO locale
	body := fmt.Sprintf("Your verification code is: %s", code)
	return s.sender.Send(to, subject, body)
}
