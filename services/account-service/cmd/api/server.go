package main

import (
	"context"
	"net/http"
	"time"

	"dietician.local/packages/middleware"
	"dietician.local/packages/tokenizer"
	accountservice "dietician.local/services/account-service"
	"dietician.local/services/account-service/config"
	"dietician.local/services/account-service/internal"
	"dietician.local/services/account-service/internal/auth/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	app         *application
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func initApplication(a *application) *Server {
	srv := &Server{
		app: a,
		mux: http.NewServeMux(),
	}

	srv.addHealthCheckRoutes()
	srv.addCommonMiddleware()

	route := accountservice.InitRoute(a.db, a.cfg, a.smtpSender, a.rdb, a.tokenizer)
	route.SetupRoutes(&internal.RouteContext{Mux: srv.mux})

	return srv
}

// Start builds the final handler (mux + middlewares) and begins listening.
func (s *Server) Start() error {
	addr := ":" + s.app.cfg.Web.Port
	s.app.logger.WithField("addr", addr).Info("starting server")

	srv := &http.Server{
		Addr:         addr,
		Handler:      s.buildHandler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return srv.ListenAndServe()
}

// buildHandler wraps the mux with all registered middlewares (outermost first).
func (s *Server) buildHandler() http.Handler {
	var h http.Handler = s.mux
	for i := len(s.middlewares) - 1; i >= 0; i-- {
		h = s.middlewares[i](h)
	}
	return h
}

// addHealthCheckRoutes registers lightweight probes that bypass middleware.
func (s *Server) addHealthCheckRoutes() {
	s.mux.HandleFunc("GET /health/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s.mux.HandleFunc("GET /health/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) addCommonMiddleware() {
	s.middlewares = append(s.middlewares, middleware.LoggerMiddleware(s.app.logger))
	s.middlewares = append(s.middlewares, middleware.LocalizerMiddleware())
	s.middlewares = append(s.middlewares, middleware.TokenizerContextMiddleware(s.app.tokenizer))
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.app.db != nil {
		s.app.db.Close()
	}
	if s.app.rdb != nil {
		s.app.rdb.Close()
	}
	return nil
}
