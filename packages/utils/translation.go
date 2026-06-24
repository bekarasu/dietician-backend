package utils

import (
	"context"
	"fmt"

	"dietician.local/packages/localizer"
	"dietician.local/packages/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	TranslateByIDWithLanguageFunc       = TranslateByIDWithLanguage
	TranslateByIDWithContextFunc        = TranslateByIDWithContext
	TranslateByTemplateWithLanguageFunc = TranslateByTemplateWithLanguage
)

func TranslateByIDWithContext(ctx context.Context, msgID string) string {
	l := middleware.GetLocalizerFromContext(ctx)
	if l != nil {
		fmt.Println(msgID)
		msg, _ := l.LocalizeMessage(&i18n.Message{
			ID: msgID,
		})
		fmt.Println(msg)
		return msg
	}
	return ""
}

func GetLanguageWithContext(ctx context.Context) string {
	return middleware.GetLanguageFromContext(ctx)
}

func GetLocalizerWithContext(ctx context.Context) *i18n.Localizer {
	return middleware.GetLocalizerFromContext(ctx)
}

func TranslateByIDWithLanguage(msgID, lang string) string {
	if bundle := localizer.GetBundle(); bundle != nil {
		l := i18n.NewLocalizer(bundle, lang)
		msg, _ := l.LocalizeMessage(&i18n.Message{
			ID: msgID,
		})

		return msg
	}

	return ""
}

func TranslateByTemplateWithLanguage(msgID, lang string, template map[string]interface{}) string {
	if bundle := localizer.GetBundle(); bundle != nil {
		l := i18n.NewLocalizer(bundle, lang)
		msg, _ := l.Localize(&i18n.LocalizeConfig{
			MessageID:    msgID,
			TemplateData: template,
		})
		return msg
	}

	return ""
}

func TranslateByTemplateWithContext(ctx context.Context, msgID string, template map[string]interface{}) string {
	l := middleware.GetLocalizerFromContext(ctx)
	if l != nil {
		msg, _ := l.Localize(&i18n.LocalizeConfig{
			MessageID:    msgID,
			TemplateData: template,
		})
		return msg
	}
	return ""
}
