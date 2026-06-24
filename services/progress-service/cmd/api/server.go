package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"

	"dietician.local/packages/middleware"
	"dietician.local/packages/swagger"
	progressservice "dietician.local/services/progress-service"
	"dietician.local/services/progress-service/config"
	_ "dietician.local/services/progress-service/docs"
	"dietician.local/services/progress-service/internal"
)

type application struct {
	logger         *logrus.Logger
	cfg            *config.ProgressAppScheme
	languageBundle *i18n.Bundle
	db             *sqlx.DB
}

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

	route := progressservice.InitRoute(a.db)
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
