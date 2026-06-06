package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"dietician.local/packages/middleware"
	"dietician.local/packages/tokenizer"
	accountservice "dietician.local/services/account-service"
	"dietician.local/services/account-service/config"
	"dietician.local/services/account-service/internal"
	"dietician.local/services/account-service/internal/auth/service"
)

type application struct {
	logger         *logrus.Logger
	languageBundle *i18n.Bundle
	cfg            *config.AccountAppScheme
	db             *sqlx.DB
	rdb            *redis.Client
	smtpSender     service.EmailSender
	tokenizer      tokenizer.ITokenizer
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

	route := accountservice.InitRoute(a.db, a.cfg, a.smtpSender, a.rdb, a.tokenizer)
	route.SetupRoutes(&internal.RouteContext{App: srv.srv})

	return srv
}

// Start begins listening.
func (s *Server) Start() error {
	addr := ":" + s.app.cfg.Web.Port
	s.app.logger.WithField("addr", addr).Info("starting server")

	return s.srv.Listen(addr)
}

// addHealthCheckRoutes registers lightweight probes that bypass middleware.
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
	s.srv.Use(middleware.TokenizerContextMiddleware(s.app.tokenizer))
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.app.db != nil {
		s.app.db.Close()
	}
	if s.app.rdb != nil {
		s.app.rdb.Close()
	}
	return s.srv.ShutdownWithContext(ctx)
}
