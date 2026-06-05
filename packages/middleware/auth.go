package middleware

import (
	"context"
	"net/http"
	"strings"

	"dietician.local/packages/constants"
	"dietician.local/packages/response"
	"dietician.local/packages/tokenizer"
)

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.UserIDKey).(string); ok {
		return v
	}
	return ""
}

func ExtractUserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID != "" {
			ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func ExtractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}

func UserAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tok, ok := r.Context().Value(constants.TokenizerKey).(tokenizer.ITokenizer)
		if !ok || tok == nil {
			response.Error(w, http.StatusInternalServerError, "tokenizer not configured")
			return
		}

		tokenStr := ExtractBearerToken(r)
			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			valid, err := tok.IsJWTInRedis(r.Context(), tokenStr)
			if err != nil || !valid {
				response.Error(w, http.StatusUnauthorized, "token revoked or not found in redis")
				return
			}

			jwtToken, err := tok.VerifyJWT(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid token signature")
				return
			}

			claims, err := tok.ExtractClaims(jwtToken)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid token claims")
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				response.Error(w, http.StatusUnauthorized, "invalid user id in token")
				return
			}

			ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
}
