package middleware

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/constants"
	"dietician.local/packages/tokenizer"
)

// TokenizerContextMiddleware injects the Tokenizer into the request context.
func TokenizerContextMiddleware(tok tokenizer.ITokenizer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(constants.TokenizerKey, tok)
		return c.Next()
	}
}
