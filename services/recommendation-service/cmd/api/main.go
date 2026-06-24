package main

// @title Recommendation Service API
// @version 1.0
// @description This is the recommendation service for the dietician backend.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@dietician.local

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("config initialization: %v", err)
	}

	logger := initLogger(cfg)

	app, err := boot(logger, cfg)
	if err != nil {
		logger.Fatalf("initialization: %v", err)
	}

	srv := initApplication(app)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	logger.Info("application successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		logger.Fatalf("server error: %v", err)
	case sig := <-quit:
		logger.Infof("received signal %s, shutting down...", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("shutdown error: %v", err)
	}

	logger.Info("shutdown complete")
}
