package validation

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"dietician.local/packages/utils"
	govalidator "github.com/go-playground/validator/v10"
)

var Validate *govalidator.Validate

func init() {
	Validate = govalidator.New(govalidator.WithRequiredStructEnabled())

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func FormatValidationError(ctx context.Context, err error) string {
	if err == nil {
		return ""
	}

	var validationErrs govalidator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var errMsgs []string
		for _, e := range validationErrs {
			field := e.Field()

			// Map tag to i18n message ID (e.g., "validation.required")
			msgID := "validation." + e.Tag()

			// Try to translate
			template := map[string]interface{}{
				"Field": field,
				"Param": e.Param(),
			}
			translated := utils.TranslateByTemplateWithContext(ctx, msgID, template)

			if translated != "" {
				errMsgs = append(errMsgs, translated)
			} else {
				// Fallback if no translation is found
				switch e.Tag() {
				case "required":
					errMsgs = append(errMsgs, fmt.Sprintf("%s is required", field))
				case "email":
					errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid email address", field))
				case "min":
					errMsgs = append(errMsgs, fmt.Sprintf("%s must be at least %s characters", field, e.Param()))
				case "len":
					errMsgs = append(errMsgs, fmt.Sprintf("%s must be exactly %s characters", field, e.Param()))
				case "gt":
					errMsgs = append(errMsgs, fmt.Sprintf("%s must be greater than %s", field, e.Param()))
				case "datetime":
					errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid date/time format", field))
				default:
					errMsgs = append(errMsgs, fmt.Sprintf("%s is invalid", field))
				}
			}
		}
		return strings.Join(errMsgs, ", ")
	}

	return err.Error()
}
