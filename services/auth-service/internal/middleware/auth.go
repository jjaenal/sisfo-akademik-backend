package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func Auth(secret string, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Request.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", nil)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		var claims jwtutil.Claims
		if err := jwtutil.ValidateWith(secret, token, &claims, cfg.JWTIssuer, cfg.JWTAudience); err != nil {
			httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", nil)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
