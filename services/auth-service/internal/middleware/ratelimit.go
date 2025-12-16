package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
)

func RateLimit(rl redisutil.Limiter, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "ratelimit:" + ip + ":" + c.FullPath()
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		n, err := rl.Incr(ctx, key)
		if err != nil {
			httputil.Error(c.Writer, 500, "6001", "Rate limiter error", err.Error())
			c.Abort()
			return
		}
		if n == 1 {
			_ = rl.Expire(ctx, key, time.Minute)
		}
		if n > int64(limit) {
			httputil.Error(c.Writer, 429, "1001", "Too Many Requests", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func RateLimitByPolicy(rl redisutil.Limiter, defaultRead int, defaultWrite int, perPrefixOverride map[string]int) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := defaultWrite
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			limit = defaultRead
		}
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		for p, v := range perPrefixOverride {
			if len(p) > 0 && len(path) >= len(p) && path[:len(p)] == p {
				limit = v
				break
			}
		}
		ip := c.ClientIP()
		key := "ratelimit:" + ip + ":" + c.Request.Method + ":" + path
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		n, err := rl.Incr(ctx, key)
		if err != nil {
			httputil.Error(c.Writer, 500, "6001", "Rate limiter error", err.Error())
			c.Abort()
			return
		}
		if n == 1 {
			_ = rl.Expire(ctx, key, time.Minute)
		}
		if n > int64(limit) {
			httputil.Error(c.Writer, 429, "1001", "Too Many Requests", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
