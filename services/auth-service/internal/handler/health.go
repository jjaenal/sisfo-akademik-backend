package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
)

type HealthHandler struct {
	db    *pgxpool.Pool
	redis *redisutil.Client
}

func NewHealthHandler(db *pgxpool.Pool, r *redisutil.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: r}
}

func (h *HealthHandler) Register(r *gin.Engine) {
	r.GET("/api/v1/health", h.health)
}

func (h *HealthHandler) health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	dbOK := true
	var dbErr string
	if _, err := h.db.Exec(ctx, "SELECT 1"); err != nil {
		dbOK = false
		dbErr = err.Error()
	}
	redisOK := true
	var redisErr string
	if err := h.redis.Raw().Ping(ctx).Err(); err != nil {
		redisOK = false
		redisErr = err.Error()
	}
	status := http.StatusOK
	if !dbOK || !redisOK {
		status = http.StatusServiceUnavailable
	}
	if status == http.StatusOK {
		httputil.Success(c.Writer, map[string]any{
			"db":          dbOK,
			"redis":       redisOK,
			"db_error":    dbErr,
			"redis_error": redisErr,
		})
		return
	}
	httputil.Error(c.Writer, http.StatusServiceUnavailable, "6002", "Infrastructure unavailable", map[string]any{
		"db":          dbOK,
		"redis":       redisOK,
		"db_error":    dbErr,
		"redis_error": redisErr,
	})
}
