package handler

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/middleware"
	"dietician.local/packages/response"
	"dietician.local/services/medical-service/internal/uploads/dto/request"
	_ "dietician.local/services/medical-service/internal/uploads/dto/response"
	"dietician.local/services/medical-service/internal/uploads/orchestration"
)

type IMedicalHandler interface {
	CreateUpload(c *fiber.Ctx) error
	ListUploads(c *fiber.Ctx) error
	GetUploadDetail(c *fiber.Ctx) error
	DeleteUpload(c *fiber.Ctx) error
}

type medicalHandler struct {
	medOrchestrator orchestration.IMedicalOrchestrator
}

func NewMedicalHandler(medOrchestrator orchestration.IMedicalOrchestrator) IMedicalHandler {
	return &medicalHandler{
		medOrchestrator: medOrchestrator,
	}
}

// CreateUpload creates a new medical upload record.
// @Summary Create Medical Upload
// @Description Create a new medical upload record for a user
// @Tags Medical Uploads
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param uploadType formData string true "Upload Type"
// @Param file formData file true "File to upload"
// @Success 200 {object} response.SuccessResponse{data=repository.MedicalUpload}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/uploads [post]
func (h *medicalHandler) CreateUpload(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "file is required")
	}

	var req request.CreateUploadRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to open file")
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to read file")
	}

	req.File = &request.FileData{
		FileName:    file.Filename,
		FileSize:    file.Size,
		ContentType: file.Header.Get("Content-Type"),
		Data:        buf,
	}

	res, err := h.medOrchestrator.CreateUpload(c.Context(), userID, req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to create upload")
	}

	return response.Success(c, "upload created successfully", res)
}

// ListUploads retrieves all medical uploads for a user.
// @Summary List Medical Uploads
// @Description Retrieve a list of medical uploads for a specific user
// @Tags Medical Uploads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]repository.MedicalUpload}
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/uploads [get]
func (h *medicalHandler) ListUploads(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())

	res, err := h.medOrchestrator.ListUploads(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to list uploads")
	}

	return response.Success(c, "uploads retrieved successfully", res)
}

// GetUploadDetail retrieves the details of a specific medical upload.
// @Summary Get Medical Upload Details
// @Description Retrieve the details of a specific medical upload, including file metadata
// @Tags Medical Uploads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uploadId path int true "Upload ID"
// @Success 200 {object} response.SuccessResponse{data=response.UploadDetailResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/uploads/{uploadId} [get]
func (h *medicalHandler) GetUploadDetail(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())
	uploadID := c.Params("uploadId")
	if uploadID == "" {
		return response.Error(c, fiber.StatusBadRequest, "invalid upload id")
	}

	res, err := h.medOrchestrator.GetUploadDetail(c.Context(), userID, uploadID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to get upload details")
	}

	return response.Success(c, "upload details retrieved successfully", res)
}

// DeleteUpload deletes a specific medical upload and its associated files.
// @Summary Delete Medical Upload
// @Description Delete a specific medical upload and its associated files
// @Tags Medical Uploads
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param uploadId path int true "Upload ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/uploads/{uploadId} [delete]
func (h *medicalHandler) DeleteUpload(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c.UserContext())
	uploadID := c.Params("uploadId")
	if uploadID == "" {
		return response.Error(c, fiber.StatusBadRequest, "invalid upload id")
	}

	err := h.medOrchestrator.DeleteUpload(c.Context(), userID, uploadID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to delete upload")
	}

	return response.Success(c, "upload deleted successfully", nil)
}
