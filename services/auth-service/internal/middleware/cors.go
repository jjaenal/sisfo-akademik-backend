package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigin := ""
		for _, ao := range allowed {
			if ao == "*" {
				allowedOrigin = "*"
				break
			}
			if ao == origin {
				allowedOrigin = origin
				break
			}
			// wildcard domain (*.example.com)
			if strings.HasPrefix(ao, "*.") && origin != "" {
				if strings.HasSuffix(origin, strings.TrimPrefix(ao, "*")) {
					allowedOrigin = origin
					break
				}
			}
		}
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	}
}

