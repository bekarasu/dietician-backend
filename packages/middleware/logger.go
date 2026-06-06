package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"dietician.local/packages/constants"
)

var (
	excludedHeaderKeys = []string{"X-Api-Key"}
	excludedPaths      = []string{
		"/v1/app/onboarding",
	}
)

func LoggerMiddleware(l *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := time.Now()

		requestFields := getRequestLogFields(c)

		ctx := context.WithValue(c.UserContext(), constants.RequestLogFieldsKey, requestFields)
		c.SetUserContext(ctx)

		if isExcludedPath(c.Path()) {
			return c.Next()
		}

		err := c.Next()

		respBody := c.Response().Body()
		respStatus := c.Response().StatusCode()

		fields := logrus.Fields{
			"request":  requestFields,
			"response": getResponseLogFields(respBody, respStatus, t),
		}

		if nrCtx, ok := c.UserContext().Value(constants.NewrelicContextKey).(context.Context); ok {
			l.WithContext(nrCtx).WithFields(fields).Info("weblogger")
		} else {
			l.WithFields(fields).Info("weblogger")
		}

		return err
	}
}

func getRequestLogFields(c *fiber.Ctx) logrus.Fields {
	bodyBytes := c.Body()

	fields := logrus.Fields{
		"id":      c.UserContext().Value(constants.RequestIDKey),
		"method":  c.Method(),
		"path":    c.Path(),
		"url":     string(c.Request().RequestURI()),
		"headers": parseRequestHeaders(c),
		"body":    unmarshalBody(bodyBytes),
	}

	return fields
}

func getResponseLogFields(body []byte, status int, t time.Time) logrus.Fields {
	return logrus.Fields{
		"status":   status,
		"duration": fmt.Sprint(time.Since(t).Round(time.Millisecond)),
		"body":     getResponseBody(body),
	}
}

func parseRequestHeaders(c *fiber.Ctx) map[string]interface{} {
	headers := make(map[string]interface{})

	c.Request().Header.VisitAll(func(k, v []byte) {
		key := string(k)
		if !isExcludedHeaderKey(key) {
			headers[key] = string(v)
		}
	})

	return headers
}

func unmarshalBody(b []byte) interface{} {
	if len(b) == 0 {
		return nil
	}
	var i interface{}
	_ = json.Unmarshal(b, &i)

	return i
}

func getResponseBody(b []byte) interface{} {
	if len(b) == 0 {
		return map[string]interface{}{"data": map[string]interface{}{}}
	}

	return unmarshalBody(b)
}

func isExcludedHeaderKey(key string) bool {
	for _, ek := range excludedHeaderKeys {
		if key == ek {
			return true
		}
	}

	return false
}

func isExcludedPath(path string) bool {
	for _, p := range excludedPaths {
		if path == p {
			return true
		}
	}

	return false
}
