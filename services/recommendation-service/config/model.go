package config

import (
	"golang.org/x/text/language"
)

type WebConfig struct {
	AppName  string `mapstructure:"APP_NAME"`
	Port     string `mapstructure:"RECOMMENDATION_SERVICE_PORT"`
	Env      string `mapstructure:"APP_ENV"`
	LogDebug bool   `mapstructure:"LOG_DEBUG"`
}

type LanguageConfig struct {
	Default   language.Tag `mapstructure:"LANGUAGE_DEFAULT"`
	Languages []language.Tag
}

type OpenAIConfig struct {
	APIKey string `mapstructure:"OPENAI_API_KEY"`
}
