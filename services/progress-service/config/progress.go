package config

import pkgconfig "dietician.local/packages/config"

var ProgressApp *ProgressAppScheme

type ProgressAppScheme struct {
	Web      WebConfig                `mapstructure:",squash"`
	Language LanguageConfig           `mapstructure:",squash"`
	Postgres pkgconfig.PostgresConfig `mapstructure:",squash"`
}
