package internal

import (
	"net/http"

	"dietician.local/packages/middleware"
	"dietician.local/services/account-service/internal/handler"
)

// RouteContext holds the dependencies passed into SetupRoutes.
type RouteContext struct {
	Mux *http.ServeMux
}

// IRoute is the contract every route group must satisfy.
type IRoute interface {
	SetupRoutes(rc *RouteContext)
}

type route struct {
	authHandler    handler.IAuthHandler
	profileHandler handler.IProfileHandler
}

// NewRoute wires all handlers into the route implementation.
func NewRoute(authHandler handler.IAuthHandler, profileHandler handler.IProfileHandler) IRoute {
	return &route{
		authHandler:    authHandler,
		profileHandler: profileHandler,
	}
}

// SetupRoutes registers all route groups on the mux.
func (r *route) SetupRoutes(rc *RouteContext) {
	r.authRoutes(rc.Mux)
	r.profileRoutes(rc.Mux)
}

func (r *route) authRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", r.authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/verify-otp", r.authHandler.VerifyOTP)
	mux.HandleFunc("POST /api/v1/auth/login", r.authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", r.authHandler.Refresh)
	mux.HandleFunc("GET /api/v1/auth/me", r.authHandler.Me)
}

func (r *route) profileRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/profiles/{userId}", middleware.UserAuthMiddleware(r.profileHandler.GetProfile))
	mux.HandleFunc("PUT /api/v1/profiles/{userId}", r.profileHandler.UpsertProfile)
	mux.HandleFunc("GET /api/v1/profiles/{userId}/preferences", r.profileHandler.GetPreferences)
	mux.HandleFunc("PUT /api/v1/profiles/{userId}/preferences", r.profileHandler.UpdatePreferences)
}
