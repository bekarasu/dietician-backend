package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"dietician.local/packages/constants"
)

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-Id")
			if reqID == "" {
				reqID = generateRequestID()
			}

			ctx := context.WithValue(r.Context(), constants.RequestIDKey, reqID)
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-Id", reqID)

			next.ServeHTTP(w, r)
		})
	}
}
