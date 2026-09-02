package repository

import (
	"context"
	"encoding/json"
	"time"

	"dietician.local/services/account-service/internal/auth/model"
	"github.com/redis/go-redis/v9"
)

type IOTPRepository interface {
	Create(ctx context.Context, otp *model.OTP) error
	GetLatestByEmail(ctx context.Context, email string) (*model.OTP, error)
	DeleteByEmail(ctx context.Context, email string) error
}

type OTPRepository struct {
	rdb *redis.Client
}

func NewOTPRepository(rdb *redis.Client) IOTPRepository {
	return &OTPRepository{rdb: rdb}
}

func (r *OTPRepository) Create(ctx context.Context, otp *model.OTP) error {
	data, err := json.Marshal(otp)
	if err != nil {
		return err
	}
	expiration := time.Until(otp.ExpiresAt)
	if expiration <= 0 {
		expiration = 5 * time.Minute
	}
	return r.rdb.Set(ctx, "otp:"+otp.Email, data, expiration).Err()
}

func (r *OTPRepository) GetLatestByEmail(ctx context.Context, email string) (*model.OTP, error) {
	data, err := r.rdb.Get(ctx, "otp:"+email).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var otp model.OTP
	if err := json.Unmarshal(data, &otp); err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepository) DeleteByEmail(ctx context.Context, email string) error {
	return r.rdb.Del(ctx, "otp:"+email).Err()
}
