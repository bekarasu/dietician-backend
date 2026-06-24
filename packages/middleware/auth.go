package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
	"dietician.local/packages/response"
	"dietician.local/packages/tokenizer"
)

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.UserIDKey).(string); ok {
		return v
	}
	return ""
}

func ExtractBearerToken(c *fiber.Ctx) string {
	auth := c.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}

func UserAuthMiddleware(c *fiber.Ctx) error {
	tok, ok := c.Locals(constants.TokenizerKey).(tokenizer.ITokenVerifier)
	if !ok || tok == nil {
		return response.Error(c, fiber.StatusInternalServerError, "tokenizer not configured")
	}

	tokenStr := ExtractBearerToken(c)
	if tokenStr == "" {
		return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}

	valid, err := tok.IsJWTInRedis(c.UserContext(), tokenStr)
	if err != nil || !valid {
		return response.Error(c, fiber.StatusUnauthorized, "token revoked or not found in redis")
	}

	jwtToken, err := tok.VerifyJWT(tokenStr)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "invalid token signature")
	}

	claims, err := tok.ExtractClaims(jwtToken)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "invalid user id in token")
	}

	ctx := context.WithValue(c.UserContext(), constants.UserIDKey, userID)
	c.SetUserContext(ctx)

	return c.Next()
}
