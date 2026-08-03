package response

import (
	"fmt"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	Message   string `json:"message,omitempty"`
	Code      int    `json:"code"`
	RequestID string `json:"requestId,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}

func Error(c *fiber.Ctx, status int, errMsg string, errs ...error) error {
	reqID := c.GetRespHeader("X-Request-Id")

	if status >= fiber.StatusInternalServerError && len(errs) > 0 && errs[0] != nil {
		buf := make([]byte, 4096)
		n := runtime.Stack(buf, false)
		stackTrace := string(buf[:n])
		c.Locals("handler_error", fmt.Sprintf("Error: %v\nStackTrace: %s", errs[0], stackTrace))
	}

	return JSON(c, status, ErrorResponse{
		Error:     errMsg,
		Code:      status,
		RequestID: reqID,
	})
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusOK, SuccessResponse{Message: message, Data: data})
}
