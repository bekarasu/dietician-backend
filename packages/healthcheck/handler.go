package healthcheck

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/response"
)

type IHealthCheckHandler interface {
	Liveness(c *fiber.Ctx) error
	Readiness(c *fiber.Ctx) error
}

type healthCheckHandler struct{}

func NewHealthCheckHandler() IHealthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) Liveness(c *fiber.Ctx) error {
	if !Liveness() {
		return response.Error(c, fiber.StatusInternalServerError, "not healthy")
	}

	return response.Success(c, "healthy", nil)
}

func (h *healthCheckHandler) Readiness(c *fiber.Ctx) error {
	readiness := Readiness()
	if !IsConnectionSuccessful(readiness) {
		return response.Error(c, fiber.StatusInternalServerError, "not healthy")
	}

	return response.Success(c, "healthy", readiness)
}
