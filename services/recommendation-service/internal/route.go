package internal

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/services/recommendation-service/internal/recommendation/handler"
)

type RouteContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(rc *RouteContext)
}

type route struct {
	recoHandler handler.IRecommendationHandler
}

func NewRoute(recoHandler handler.IRecommendationHandler) IRoute {
	return &route{
		recoHandler: recoHandler,
	}
}

func (r *route) SetupRoutes(rc *RouteContext) {
	r.recoRoutes(rc.App)
}

func (r *route) recoRoutes(app *fiber.App) {
	recoGroup := app.Group("/api/v1/recommendations")
	recoGroup.Post("/", r.recoHandler.CreateDietRecommendations)
}
