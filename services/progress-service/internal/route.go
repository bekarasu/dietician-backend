package internal

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/services/progress-service/internal/progress/handler"
	dailylogHandler "dietician.local/services/progress-service/internal/dailylog/handler"
	trackingHandler "dietician.local/services/progress-service/internal/tracking/handler"
)

type RouteContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(rc *RouteContext)
}

type route struct {
	progressHandler handler.IProgressHandler
	dailylogHandler dailylogHandler.IDailyLogHandler
	trackingHandler trackingHandler.ITrackingHandler
}

func NewRoute(progressHandler handler.IProgressHandler, dailylogHandler dailylogHandler.IDailyLogHandler, trackingHandler trackingHandler.ITrackingHandler) IRoute {
	return &route{
		progressHandler: progressHandler,
		dailylogHandler: dailylogHandler,
		trackingHandler: trackingHandler,
	}
}

func (r *route) SetupRoutes(rc *RouteContext) {
	r.progressRoutes(rc.App)
	r.dailylogRoutes(rc.App)
	r.trackingRoutes(rc.App)
}

func (r *route) progressRoutes(app *fiber.App) {
	progressGroup := app.Group("/api/v1/progress/:userId")

	progressGroup.Get("/", r.progressHandler.GetProgress)
	progressGroup.Post("/weight", r.progressHandler.AddWeight)
	progressGroup.Get("/weekly-summary", r.progressHandler.GetWeeklySummary)
	progressGroup.Post("/habits", r.progressHandler.AddHabit)
}

func (r *route) dailylogRoutes(app *fiber.App) {
	dlGroup := app.Group("/api/v1/daily-logs")
	dlGroup.Post("/", r.dailylogHandler.CreateDailyLog)
	dlGroup.Get("/:user_id/:date", r.dailylogHandler.GetDailyLog)
}

func (r *route) trackingRoutes(app *fiber.App) {
	trGroup := app.Group("/api/v1/tracking")
	trGroup.Post("/", r.trackingHandler.CreateTrackingMetric)
	trGroup.Get("/:user_id", r.trackingHandler.GetTrackingMetrics)
}
