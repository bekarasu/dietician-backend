package repository

import (
	"context"
	"database/sql"
	"time"

	"dietician.local/services/account-service/internal/auth/model"
	"github.com/jmoiron/sqlx"
)

// RefreshTokenRepo abstracts refresh-token persistence.
type IRefreshTokenRepo interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID string) error
}

type RefreshTokenRepository struct {
	db *sqlx.DB
}

func NewRefreshTokenRepository(db *sqlx.DB) IRefreshTokenRepo {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query,
		token.UserID, token.Token, token.ExpiresAt,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	query := `SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens WHERE token = $1`
	err := r.db.GetContext(ctx, &rt, query, token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rt, err
}

func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < $1`, time.Now())
	return err
}
