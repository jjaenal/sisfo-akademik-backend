package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

type Authorizer interface {
	Allow(subjectID uuid.UUID, tenantID string, permission string) (bool, error)
}

func Authorization(authz Authorizer, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("claims")
		if !ok {
			httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", nil)
			c.Abort()
			return
		}
		claims, _ := val.(jwtutil.Claims)
		okay, err := authz.Allow(claims.UserID, claims.TenantID, permission)
		if err != nil {
			httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal error", nil)
			c.Abort()
			return
		}
		if !okay {
			httputil.Error(c.Writer, http.StatusForbidden, "3001", "Forbidden", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

