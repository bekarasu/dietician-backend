package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"dietician.local/packages/localizer"
	"dietician.local/packages/logging"
	"dietician.local/packages/tokenizer"
	"dietician.local/services/progress-service/config"
)

func boot(logger *logrus.Logger, cfg *config.ProgressAppScheme) (*application, error) {
	db, err := initPostgres(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("postgres initialization: %v", err)
	}

	bundle := initLocalizer(cfg)

	rdb := initRedis(cfg, logger)

	tok := tokenizer.NewTokenVerifier(tokenizer.Config{
		Secret: cfg.JWT.Secret,
	}, rdb)

	return &application{
		logger:         logger,
		cfg:            cfg,
		languageBundle: bundle,
		db:             db,
		rdb:            rdb,
		tokenizer:      tok,
	}, nil
}

func initConfig() (*config.ProgressAppScheme, error) {
	return config.Load()
}

func initLogger(cfg *config.ProgressAppScheme) *logrus.Logger {
	return logging.NewLogger(logging.Config{
		Service: logging.ServiceConfig{
			Env:     cfg.Web.Env,
			AppName: "progress-service",
		},
	})
}

func initLocalizer(cfg *config.ProgressAppScheme) *i18n.Bundle {
	return localizer.InitLocalizer(localizer.Config{
		Default: language.Turkish,
		Languages: []language.Tag{
			language.English,
			language.Turkish,
		},
	})
}

func initPostgres(cfg *config.ProgressAppScheme, logger *logrus.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", cfg.Postgres.DSN())
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	logger.Info("connected to postgres")
	return db, nil
}

func initRedis(cfg *config.ProgressAppScheme, logger *logrus.Logger) *redis.Client {
	dbNum, _ := strconv.Atoi(cfg.Redis.DB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       dbNum,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.WithError(err).Warn("redis not available")
	} else {
		logger.Info("connected to redis")
	}

	return rdb
}
