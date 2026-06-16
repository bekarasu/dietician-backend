package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"dietician.local/packages/middleware"
	"dietician.local/packages/swagger"
	medicalservice "dietician.local/services/medical-service"
	_ "dietician.local/services/medical-service/docs"
	"dietician.local/services/medical-service/internal"
)

type Server struct {
	app *application
	srv *fiber.App
}

func initApplication(a *application) *Server {
	srv := &Server{
		app: a,
		srv: fiber.New(fiber.Config{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}),
	}

	srv.addHealthCheckRoutes()
	srv.addCommonMiddleware()

	route := medicalservice.InitRoute(a.db, a.cfg, a.storageProvider)
	route.SetupRoutes(&internal.RouteContext{App: srv.srv})

	// Setup Swagger UI
	swagger.Setup(srv.srv, "/swagger")

	return srv
}

func (s *Server) Start() error {
	addr := ":" + s.app.cfg.Web.Port
	s.app.logger.WithField("addr", addr).Info("starting server")

	return s.srv.Listen(addr)
}

func (s *Server) addHealthCheckRoutes() {
	s.srv.Get("/health/liveness", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	s.srv.Get("/health/readiness", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	// Legacy endpoint from README
	s.srv.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

func (s *Server) addCommonMiddleware() {
	s.srv.Use(middleware.LoggerMiddleware(s.app.logger))
	s.srv.Use(middleware.LocalizerMiddleware())
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.app.db != nil {
		s.app.db.Close()
	}
	return s.srv.ShutdownWithContext(ctx)
}
