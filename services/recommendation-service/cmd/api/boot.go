package main

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"

	"dietician.local/packages/localizer"
	"dietician.local/packages/logging"
	"dietician.local/packages/openai"
	"dietician.local/services/recommendation-service/config"
)

func boot(logger *logrus.Logger, cfg *config.RecommendationAppScheme) (*application, error) {
	bundle := initLocalizer(cfg)

	openAIService := openai.NewService(openai.Config{
		APIKey: cfg.OpenAI.APIKey,
	}, logger)

	return &application{
		logger:         logger,
		cfg:            cfg,
		languageBundle: bundle,
		openaiService:  openAIService,
	}, nil
}

func initConfig() (*config.RecommendationAppScheme, error) {
	return config.Load()
}

func initLogger(cfg *config.RecommendationAppScheme) *logrus.Logger {
	return logging.NewLogger(logging.Config{
		Service: logging.ServiceConfig{
			Env:     cfg.Web.Env,
			AppName: "recommendation-service",
		},
	})
}

func initLocalizer(cfg *config.RecommendationAppScheme) *i18n.Bundle {
	return localizer.InitLocalizer(localizer.Config{
		Default: language.Turkish,
		Languages: []language.Tag{
			language.English,
			language.Turkish,
		},
	})
}
