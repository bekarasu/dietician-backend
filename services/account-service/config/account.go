package config

import pkgconfig "dietician.local/packages/config"

var AuthApp *AccountAppScheme

type AccountAppScheme struct {
	Web         WebConfig                `mapstructure:",squash"`
	Credentials Credentials              `mapstructure:",squash"`
	Language    LanguageConfig           `mapstructure:",squash"`
	Postgres    pkgconfig.PostgresConfig `mapstructure:",squash"`
	Redis       pkgconfig.RedisConfig    `mapstructure:",squash"`
	SMTP        SMTPConfig               `mapstructure:",squash"`
	JWT         JWTConfig                `mapstructure:",squash"`
}
