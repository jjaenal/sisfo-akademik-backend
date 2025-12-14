package middleware

import (
	"net/http"

	"shared/pkg/httputil"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				httputil.Error(w, http.StatusInternalServerError, "1001", "Internal error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
