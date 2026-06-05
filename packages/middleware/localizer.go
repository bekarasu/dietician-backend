package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"dietician.local/packages/localizer"
)

type contextKey string

const (
	LocalizerKey contextKey = "localizer"
	LanguageKey  contextKey = "language"
)

func LocalizerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := r.Header.Get("Accept-Language")

			if lang != "" {
				lang = strings.Split(lang, ",")[0]
				lang = strings.Split(lang, "-")[0]
				lang = strings.TrimSpace(lang)
			}

			var l *i18n.Localizer
			var selectedLang string

			if lang == localizer.EnglishLanguage || lang == "en" {
				l = localizer.GetEnglishLocalizer()
				selectedLang = localizer.EnglishLanguage
			} else {
				l = localizer.GetTurkishLocalizer()
				selectedLang = localizer.TurkishLanguage
			}

			ctx := context.WithValue(r.Context(), LocalizerKey, l)
			ctx = context.WithValue(ctx, LanguageKey, selectedLang)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLocalizerFromContext(ctx context.Context) *i18n.Localizer {
	if l, ok := ctx.Value(LocalizerKey).(*i18n.Localizer); ok {
		return l
	}
	return localizer.GetTurkishLocalizer()
}

func GetLanguageFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(LanguageKey).(string); ok {
		return lang
	}
	return localizer.TurkishLanguage
}
