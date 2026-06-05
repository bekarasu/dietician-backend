package config

import (
	"golang.org/x/text/language"
)

type Credentials struct {
	AppID  string `mapstructure:"APP_ID"`
	Secret string `mapstructure:"APP_SECRET"`
}

type WebConfig struct {
	AppName  string `mapstructure:"APP_NAME"`
	Port     string `mapstructure:"AUTH_SERVICE_PORT"`
	Env      string `mapstructure:"APP_ENV"`
	LogDebug bool   `mapstructure:"LOG_DEBUG"`
}

type LanguageConfig struct {
	Default   language.Tag `mapstructure:"LANGUAGE_DEFAULT"`
	Languages []language.Tag
}

type SMTPConfig struct {
	Host     string `mapstructure:"SMTP_HOST"`
	Port     string `mapstructure:"SMTP_PORT"`
	Username string `mapstructure:"SMTP_USERNAME"`
	Password string `mapstructure:"SMTP_PASSWORD"`
}

type JWTConfig struct {
	Secret       string `mapstructure:"JWT_SECRET"`
	OTPSecret    string `mapstructure:"OTP_JWT_SECRET"`
	OTPExpireSec int    `mapstructure:"OTP_EXPIRE_TIME_IN_SECONDS"`
}
