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

