package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"shared/pkg/httputil"
	redisutil "shared/pkg/redis"
)

func clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func RateLimit(lim redisutil.Limiter, limitPerMin int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := fmt.Sprintf("ratelimit:%s:%s", clientIP(r), r.URL.Path)
		n, err := lim.Incr(context.Background(), key)
		if err != nil {
			httputil.Error(w, http.StatusInternalServerError, "1001", "Internal error", nil)
			return
		}
		if n == 1 {
			_ = lim.Expire(context.Background(), key, time.Minute)
		}
		if n > int64(limitPerMin) {
			httputil.Error(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
