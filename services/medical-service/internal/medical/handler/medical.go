package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/response"
	"dietician.local/services/medical-service/internal/medical/service"
)

type IMedicalHandler interface {
	CreateUpload(c *fiber.Ctx) error
	ListUploads(c *fiber.Ctx) error
	GetUploadDetail(c *fiber.Ctx) error
	DeleteUpload(c *fiber.Ctx) error
}

type medicalHandler struct {
	medService service.IMedicalService
}

func NewMedicalHandler(medService service.IMedicalService) IMedicalHandler {
	return &medicalHandler{
		medService: medService,
	}
}

func (h *medicalHandler) CreateUpload(c *fiber.Ctx) error {
	userID := c.Params("userId")
	var req service.CreateUploadRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	res, err := h.medService.CreateUpload(c.Context(), userID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create upload")
	}

	return response.Success(c, "upload created successfully", res)
}

func (h *medicalHandler) ListUploads(c *fiber.Ctx) error {
	userID := c.Params("userId")

	res, err := h.medService.ListUploads(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to list uploads")
	}

	return response.Success(c, "uploads retrieved successfully", res)
}

func (h *medicalHandler) GetUploadDetail(c *fiber.Ctx) error {
	userID := c.Params("userId")
	uploadID, err := strconv.ParseInt(c.Params("uploadId"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid upload id")
	}

	res, err := h.medService.GetUploadDetail(c.Context(), userID, uploadID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get upload details")
	}

	return response.Success(c, "upload details retrieved successfully", res)
}

func (h *medicalHandler) DeleteUpload(c *fiber.Ctx) error {
	userID := c.Params("userId")
	uploadID, err := strconv.ParseInt(c.Params("uploadId"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid upload id")
	}

	err = h.medService.DeleteUpload(c.Context(), userID, uploadID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete upload")
	}

	return response.Success(c, "upload deleted successfully", nil)
}
