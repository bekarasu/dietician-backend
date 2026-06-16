package config

import (
	"os"

	"golang.org/x/text/language"

	"dietician.local/packages/viperconfig"
	"github.com/spf13/viper"
)

func Load() (*RecommendationAppScheme, error) {
	path := "."
	if envPath := os.Getenv("CONFIG_FILE_PATH"); envPath != "" {
		path = envPath
	}

	file := ".env"
	if envFile := os.Getenv("CONFIG_FILE_NAME"); envFile != "" {
		file = envFile
	}

	viper.Set("LANGUAGES", []language.Tag{language.English, language.Turkish})
	viper.Set("LANGUAGE_DEFAULT", language.Turkish)

	vc := viperconfig.Config{
		Path:     path,
		FileName: file,
		Env:      os.Getenv("ENV"),
	}

	c, err := viperconfig.Load(vc, RecommendationAppScheme{})
	if err != nil {
		return nil, err
	}

	cfg := c.(RecommendationAppScheme)
	RecoApp = &cfg

	return RecoApp, nil
}
