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

// ITokenVerifier provides read-only JWT operations: verify, extract claims,
// check Redis session, and basic auth. Used by services that only need to
// validate incoming tokens (e.g. medical-service, progress-service).
type ITokenVerifier interface {
	VerifyJWT(token string) (*jwt.Token, error)
	ExtractClaims(j *jwt.Token) (jwt.MapClaims, error)
	IsJWTInRedis(ctx context.Context, token string) (bool, error)
	IsBasicAuthorized(token string, username, password string) bool
}

// ITokenWriter provides write JWT operations: generate tokens and manage
// Redis sessions. Used only by account-service.
type ITokenWriter interface {
	GenerateJWT(claims jwt.MapClaims) (string, error)
	StoreJWTInRedis(ctx context.Context, token string, expiration time.Duration) error
	RemoveJWTFromRedis(ctx context.Context, token string) error
}

// ITokenizer combines both read and write JWT operations.
// Used by account-service which needs full access.
type ITokenizer interface {
	ITokenVerifier
	ITokenWriter
}

type tokenizer struct {
	config Config
	redis  *redis.Client
}

// NewTokenizer creates a full-access tokenizer (read + write).
// Use this in account-service.
func NewTokenizer(c Config, r *redis.Client) ITokenizer {
	return &tokenizer{
		config: c,
		redis:  r,
	}
}

// NewTokenVerifier creates a read-only tokenizer (verify + extract only).
// Use this in services that only need to validate incoming JWTs.
func NewTokenVerifier(c Config, r *redis.Client) ITokenVerifier {
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
