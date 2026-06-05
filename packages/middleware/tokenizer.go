package middleware

import (
	"context"
	"net/http"

	"dietician.local/packages/constants"
	"dietician.local/packages/tokenizer"
)

// TokenizerContextMiddleware injects the Tokenizer into the request context.
func TokenizerContextMiddleware(tok tokenizer.ITokenizer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), constants.TokenizerKey, tok)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
