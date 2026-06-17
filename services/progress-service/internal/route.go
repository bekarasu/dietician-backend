package internal

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/services/progress-service/internal/progress/handler"
)

type RouteContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(rc *RouteContext)
}

type route struct {
	progressHandler handler.IProgressHandler
}

func NewRoute(progressHandler handler.IProgressHandler) IRoute {
	return &route{
		progressHandler: progressHandler,
	}
}

func (r *route) SetupRoutes(rc *RouteContext) {
	r.progressRoutes(rc.App)
}

func (r *route) progressRoutes(app *fiber.App) {
	progressGroup := app.Group("/api/v1/progress/:userId")

	progressGroup.Get("/", r.progressHandler.GetProgress)
	progressGroup.Post("/weight", r.progressHandler.AddWeight)
	progressGroup.Get("/weekly-summary", r.progressHandler.GetWeeklySummary)
	progressGroup.Post("/habits", r.progressHandler.AddHabit)
}
