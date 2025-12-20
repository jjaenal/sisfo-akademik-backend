package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	smiddleware "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/middleware"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
)

func RateLimitByPolicy(rl redisutil.Limiter, defaultRead int, defaultWrite int, perPrefixOverride map[string]int) gin.HandlerFunc {
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { c.Next() })
		h := smiddleware.RateLimitByPolicy(rl, defaultRead, defaultWrite, perPrefixOverride, next)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
