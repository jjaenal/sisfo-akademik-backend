package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
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

func RateLimitByPrefix(lim redisutil.Limiter, defaultLimit int, perPrefix map[string]int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := defaultLimit
		path := r.URL.Path
		for p, v := range perPrefix {
			if len(p) > 0 && len(path) >= len(p) && path[:len(p)] == p {
				limit = v
			}
		}
		key := fmt.Sprintf("ratelimit:%s:%s", clientIP(r), path)
		n, err := lim.Incr(context.Background(), key)
		if err != nil {
			httputil.Error(w, http.StatusInternalServerError, "1001", "Internal error", nil)
			return
		}
		if n == 1 {
			_ = lim.Expire(context.Background(), key, time.Minute)
		}
		if n > int64(limit) {
			httputil.Error(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RateLimitByPolicy(lim redisutil.Limiter, defaultRead int, defaultWrite int, perPrefixOverride map[string]int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := defaultWrite
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			limit = defaultRead
		}
		for p, v := range perPrefixOverride {
			if len(p) > 0 && len(r.URL.Path) >= len(p) && r.URL.Path[:len(p)] == p {
				limit = v
				break
			}
		}
		key := fmt.Sprintf("ratelimit:%s:%s:%s", clientIP(r), r.Method, r.URL.Path)
		n, err := lim.Incr(context.Background(), key)
		if err != nil {
			httputil.Error(w, http.StatusInternalServerError, "1001", "Internal error", nil)
			return
		}
		if n == 1 {
			_ = lim.Expire(context.Background(), key, time.Minute)
		}
		if n > int64(limit) {
			httputil.Error(w, http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
