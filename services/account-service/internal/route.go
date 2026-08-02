package internal

import (
	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/middleware"
	authhandler "dietician.local/services/account-service/internal/auth/handler"
	profilehandler "dietician.local/services/account-service/internal/profile/handler"
)

// RouteContext holds the dependencies passed into SetupRoutes.
type RouteContext struct {
	App *fiber.App
}

// IRoute is the contract every route group must satisfy.
type IRoute interface {
	SetupRoutes(rc *RouteContext)
}

type route struct {
	authHandler    authhandler.IAuthHandler
	profileHandler profilehandler.IProfileHandler
}

// NewRoute wires all handlers into the route implementation.
func NewRoute(authHandler authhandler.IAuthHandler, profileHandler profilehandler.IProfileHandler) IRoute {
	return &route{
		authHandler:    authHandler,
		profileHandler: profileHandler,
	}
}

// SetupRoutes registers all route groups on the app.
func (r *route) SetupRoutes(rc *RouteContext) {
	r.authRoutes(rc.App)
	r.profileRoutes(rc.App)
}

func (r *route) authRoutes(app *fiber.App) {
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/register", r.authHandler.Register)
	authGroup.Post("/verify-otp", r.authHandler.VerifyOTP)
	authGroup.Post("/login", r.authHandler.Login)
	authGroup.Post("/refresh", r.authHandler.Refresh)
	authGroup.Post("/logout", middleware.UserAuthMiddleware, r.authHandler.Logout)
}

func (r *route) profileRoutes(app *fiber.App) {
	profileGroup := app.Group("/api/v1/profiles", middleware.UserAuthMiddleware)

	profileGroup.Get("/", r.profileHandler.GetProfile)
	profileGroup.Put("/", r.profileHandler.UpsertProfile)
	profileGroup.Put("/weight", r.profileHandler.UpdateWeight)
	profileGroup.Post("/onboarding", r.profileHandler.Onboarding)
	profileGroup.Get("/preferences", r.profileHandler.GetPreferences)
	profileGroup.Put("/preferences", r.profileHandler.UpdatePreferences)
}
