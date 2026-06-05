package localizer

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle *i18n.Bundle

	trLocalizer     *i18n.Localizer
	trLocalizerOnce sync.Once

	enLocalizer     *i18n.Localizer
	enLocalizerOnce sync.Once
)

type Config struct {
	Default   language.Tag
	Languages []language.Tag
}

func InitLocalizer(config Config) *i18n.Bundle {
	b := i18n.NewBundle(config.Default)
	b.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, lang := range config.Languages {
		b.MustLoadMessageFile(fmt.Sprintf("locale/system.%s.toml", lang.String()))
		b.MustLoadMessageFile(fmt.Sprintf("locale/active.%s.toml", lang.String()))
		b.LoadMessageFile(fmt.Sprintf("locale/validation.%s.toml", lang.String()))
	}

	bundle = b

	return bundle
}

func GetBundle() *i18n.Bundle {
	return bundle
}

func GetTurkishLocalizer() *i18n.Localizer {
	trLocalizerOnce.Do(func() {
		trLocalizer = i18n.NewLocalizer(GetBundle(), TurkishLanguage)
	})

	return trLocalizer
}

func GetEnglishLocalizer() *i18n.Localizer {
	enLocalizerOnce.Do(func() {
		enLocalizer = i18n.NewLocalizer(GetBundle(), EnglishLanguage)
	})

	return enLocalizer
}
