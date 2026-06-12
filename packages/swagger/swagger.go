package swagger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Setup registers the Swagger UI routes on the provided Fiber app.
// It serves the swagger documentation under the given basePath (e.g., "/swagger").
func Setup(app *fiber.App, basePath string) {
	if basePath == "" {
		basePath = "/swagger"
	}
	app.Get(basePath+"/*", swagger.HandlerDefault)
}
