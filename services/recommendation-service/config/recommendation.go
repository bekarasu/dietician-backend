package config

import pkgconfig "dietician.local/packages/config"

var RecoApp *RecommendationAppScheme

type RecommendationAppScheme struct {
	Web      WebConfig                `mapstructure:",squash"`
	Language LanguageConfig           `mapstructure:",squash"`
	OpenAI   OpenAIConfig             `mapstructure:",squash"`
	Postgres pkgconfig.PostgresConfig `mapstructure:",squash"`
}
