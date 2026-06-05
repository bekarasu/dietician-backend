package tokenizer

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Secret string
}

type BasicAuthentication struct {
	Username string
	Password string
}

type ITokenizer interface {
	GenerateJWT(claims jwt.MapClaims) (string, error)
	VerifyJWT(token string) (*jwt.Token, error)
	ExtractClaims(j *jwt.Token) (jwt.MapClaims, error)
	StoreJWTInRedis(ctx context.Context, token string, expiration time.Duration) error
	IsJWTInRedis(ctx context.Context, token string) (bool, error)
	RemoveJWTFromRedis(ctx context.Context, token string) error
	IsBasicAuthorized(token string, username, password string) bool
}

type tokenizer struct {
	config Config
	redis  *redis.Client
}

func NewTokenizer(c Config, r *redis.Client) ITokenizer {
	return &tokenizer{
		config: c,
		redis:  r,
	}
}

func (t *tokenizer) GenerateJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.config.Secret))
}

func (t *tokenizer) VerifyJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.config.Secret), nil
	})
}

func (t *tokenizer) ExtractClaims(j *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := j.Claims.(jwt.MapClaims); ok && j.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid jwt claims")
}

func (t *tokenizer) getRedisKey(token string) string {
	return fmt.Sprintf("jwt:%s", token)
}

func (t *tokenizer) StoreJWTInRedis(ctx context.Context, token string, expiration time.Duration) error {
	if t.redis == nil {
		return errors.New("redis client is not initialized")
	}
	return t.redis.Set(ctx, t.getRedisKey(token), "valid", expiration).Err()
}

func (t *tokenizer) IsJWTInRedis(ctx context.Context, token string) (bool, error) {
	if t.redis == nil {
		return false, errors.New("redis client is not initialized")
	}
	val, err := t.redis.Get(ctx, t.getRedisKey(token)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "valid", nil
}

func (t *tokenizer) RemoveJWTFromRedis(ctx context.Context, token string) error {
	if t.redis == nil {
		return errors.New("redis client is not initialized")
	}
	return t.redis.Del(ctx, t.getRedisKey(token)).Err()
}

func (t *tokenizer) IsBasicAuthorized(token string, username, password string) bool {
	decode, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	u, p, ok := strings.Cut(string(decode), ":")
	if !ok {
		return false
	}

	return u == username && p == password
}
