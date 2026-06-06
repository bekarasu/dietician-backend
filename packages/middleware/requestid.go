package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
)

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-Id")
		if reqID == "" {
			reqID = generateRequestID()
		}

		ctx := context.WithValue(c.UserContext(), constants.RequestIDKey, reqID)
		c.SetUserContext(ctx)
		c.Set("X-Request-Id", reqID)

		return c.Next()
	}
}
