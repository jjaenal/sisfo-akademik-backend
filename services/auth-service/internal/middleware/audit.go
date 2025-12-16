package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func Audit(repo *repository.AuditRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		lat := time.Since(start)
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := c.Writer.Status()
		tenant := ""
		var userID *uuid.UUID
		if v, ok := c.Get("claims"); ok {
			if cl, ok2 := v.(jwtutil.Claims); ok2 {
				tenant = cl.TenantID
				userID = &cl.UserID
			}
		}
		m := method
		p := path
		s := status
		t := tenant
		uid := userID
		dur := lat.Milliseconds()
		go func() {
			_ = repo.Log(context.Background(), t, uid, "http.request", "route", nil, map[string]any{
				"method":   m,
				"path":     p,
				"status":   s,
				"duration": dur,
			})
		}()
	}
}
