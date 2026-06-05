package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"dietician.local/packages/constants"
	"github.com/sirupsen/logrus"
)

var (
	excludedHeaderKeys = []string{"X-Api-Key"}
	excludedPaths      = []string{
		"/v1/app/onboarding",
	}
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rw *responseWriterWrapper) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func LoggerMiddleware(l *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			requestFields := getRequestLogFields(r)

			ctx := context.WithValue(r.Context(), constants.RequestLogFieldsKey, requestFields)
			r = r.WithContext(ctx)

			if isExcludedPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			rw := &responseWriterWrapper{
				ResponseWriter: w,
				status:         http.StatusOK, // Default to 200
				body:           &bytes.Buffer{},
			}

			next.ServeHTTP(rw, r)

			respBody := rw.body.Bytes()
			respStatus := rw.status

			fields := logrus.Fields{
				"request":  requestFields,
				"response": getResponseLogFields(respBody, respStatus, t),
			}

			if nrCtx, ok := r.Context().Value(constants.NewrelicContextKey).(context.Context); ok {
				l.WithContext(nrCtx).WithFields(fields).Info("weblogger")
			} else {
				l.WithFields(fields).Info("weblogger")
			}
		})
	}
}

func getRequestLogFields(r *http.Request) logrus.Fields {
	var bodyBytes []byte

	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	fields := logrus.Fields{
		"id":      r.Context().Value(constants.RequestIDKey),
		"method":  r.Method,
		"path":    r.URL.Path,
		"url":     r.RequestURI,
		"headers": parseRequestHeaders(r),
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

func parseRequestHeaders(r *http.Request) map[string]interface{} {
	headers := make(map[string]interface{})

	for k, v := range r.Header {
		if !isExcludedHeaderKey(k) {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
	}

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
