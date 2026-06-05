package repository

import (
	"context"
	"database/sql"

	"dietician.local/services/account-service/internal/auth/model"
	"github.com/jmoiron/sqlx"
)

type IOTPRepository interface {
	Create(ctx context.Context, otp *model.OTP) error
	GetLatestByEmail(ctx context.Context, email string) (*model.OTP, error)
	DeleteByEmail(ctx context.Context, email string) error
}

type OTPRepository struct {
	db *sqlx.DB
}

func NewOTPRepository(db *sqlx.DB) IOTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(ctx context.Context, otp *model.OTP) error {
	query := `INSERT INTO otps (email, otp_code, expires_at)
		VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		otp.Email, otp.OTPCode, otp.ExpiresAt,
	).Scan(&otp.ID, &otp.CreatedAt)
}

func (r *OTPRepository) GetLatestByEmail(ctx context.Context, email string) (*model.OTP, error) {
	var otp model.OTP
	query := `SELECT id, email, otp_code, expires_at, created_at
		FROM otps WHERE email = $1 ORDER BY created_at DESC LIMIT 1`
	err := r.db.GetContext(ctx, &otp, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &otp, err
}

func (r *OTPRepository) DeleteByEmail(ctx context.Context, email string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM otps WHERE email = $1`, email)
	return err
}
