package internal

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/middleware"
	"dietician.local/services/medical-service/internal/uploads/handler"
)

type RouteContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(rc *RouteContext)
}

type route struct {
	medicalHandler handler.IMedicalHandler
}

func NewRoute(medicalHandler handler.IMedicalHandler) IRoute {
	return &route{
		medicalHandler: medicalHandler,
	}
}

func (r *route) SetupRoutes(rc *RouteContext) {
	r.medicalRoutes(rc.App)
}

func (r *route) medicalRoutes(app *fiber.App) {
	medicalGroup := app.Group("/api/v1/uploads", middleware.UserAuthMiddleware)

	medicalGroup.Post("/", r.medicalHandler.CreateUpload)
	medicalGroup.Get("/", r.medicalHandler.ListUploads)
	medicalGroup.Get("/:uploadId", r.medicalHandler.GetUploadDetail)
	medicalGroup.Delete("/:uploadId", r.medicalHandler.DeleteUpload)
}
