package middleware

import (
	"context"
	"net/http"
	"strings"

	"shared/pkg/httputil"
	jwtutil "shared/pkg/jwt"
)

type ctxKey string

const ClaimsKey ctxKey = "claims"

func Auth(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			httputil.Error(w, http.StatusUnauthorized, "2001", "Unauthorized", nil)
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		var c jwtutil.Claims
		if err := jwtutil.Validate(secret, token, &c); err != nil {
			httputil.Error(w, http.StatusUnauthorized, "2001", "Unauthorized", nil)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
