package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockLimiter struct {
	count int64
	err   error
}

func (m *mockLimiter) Incr(ctx context.Context, key string) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	m.count++
	return m.count, nil
}

func (m *mockLimiter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return m.err
}

func TestRateLimitByPolicy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("allowed request", func(t *testing.T) {
		limiter := &mockLimiter{count: 0}
		r := gin.New()
		r.Use(RateLimitByPolicy(limiter, 10, 10, nil))
		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("blocked request", func(t *testing.T) {
		limiter := &mockLimiter{count: 11} // Exceeds default limit of 10
		r := gin.New()
		r.Use(RateLimitByPolicy(limiter, 10, 10, nil))
		r.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})
}
