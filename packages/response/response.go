package response

import (
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

func Error(c *fiber.Ctx, status int, errMsg string) error {
	reqID := c.GetRespHeader("X-Request-Id")
	return JSON(c, status, ErrorResponse{
		Error:     errMsg,
		Code:      status,
		RequestID: reqID,
	})
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusOK, SuccessResponse{Message: message, Data: data})
}
