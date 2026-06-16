package config

var RecoApp *RecommendationAppScheme

type RecommendationAppScheme struct {
	Web      WebConfig      `mapstructure:",squash"`
	Language LanguageConfig `mapstructure:",squash"`
	OpenAI   OpenAIConfig   `mapstructure:",squash"`
}
