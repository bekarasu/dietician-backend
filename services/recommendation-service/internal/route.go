package internal

import (
	"github.com/gofiber/fiber/v2"

	dietplanHandler "dietician.local/services/recommendation-service/internal/dietplan/handler"
	"dietician.local/services/recommendation-service/internal/recommendation/handler"
)

type RouteContext struct {
	App *fiber.App
}

type IRoute interface {
	SetupRoutes(rc *RouteContext)
}

type route struct {
	recoHandler     handler.IRecommendationHandler
	dietplanHandler dietplanHandler.IDietPlanHandler
}

func NewRoute(recoHandler handler.IRecommendationHandler, dietplanHandler dietplanHandler.IDietPlanHandler) IRoute {
	return &route{
		recoHandler:     recoHandler,
		dietplanHandler: dietplanHandler,
	}
}

func (r *route) SetupRoutes(rc *RouteContext) {
	r.recoRoutes(rc.App)
	r.dietplanRoutes(rc.App)
}

func (r *route) recoRoutes(app *fiber.App) {
	recoGroup := app.Group("/api/v1/recommendations")
	recoGroup.Post("/", r.recoHandler.CreateDietRecommendations)
}

func (r *route) dietplanRoutes(app *fiber.App) {
	dpGroup := app.Group("/api/v1/diet-plans")
	dpGroup.Post("/", r.dietplanHandler.CreateDietPlan)
}
