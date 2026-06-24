package config

import pkgconfig "dietician.local/packages/config"

var MedicalApp *MedicalAppScheme

type MedicalAppScheme struct {
	Web      WebConfig                `mapstructure:",squash"`
	Language LanguageConfig           `mapstructure:",squash"`
	Postgres pkgconfig.PostgresConfig `mapstructure:",squash"`
	Redis    pkgconfig.RedisConfig    `mapstructure:",squash"`
	JWT      JWTConfig                `mapstructure:",squash"`
}
