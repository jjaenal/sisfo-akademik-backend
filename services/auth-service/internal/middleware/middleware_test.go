package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func makeCfg() config.Config {
	return config.Config{
		JWTIssuer:   "sisfo-akademik",
		JWTAudience: "api",
	}
}

type fakeAuthz struct {
	allow bool
	err   error
}

func (f *fakeAuthz) Allow(subjectID uuid.UUID, tenantID string, permission string) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.allow, nil
}

func TestAuthorization_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authz := &fakeAuthz{allow: false}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t1", UserID: uuid.New()})
	})
	r.Use(Authorization(authz, "user:read"))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

func TestAuthorization_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authz := &fakeAuthz{allow: true}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t1", UserID: uuid.New()})
	})
	r.Use(Authorization(authz, "user:read"))
	r.GET("/", func(c *gin.Context) { c.Status(204) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestAuthorization_Error_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authz := &fakeAuthz{err: errors.New("backend error")}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t1", UserID: uuid.New()})
	})
	r.Use(Authorization(authz, "user:read"))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestAuthorization_Unauthorized_NoClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authz := &fakeAuthz{allow: true}
	r := gin.New()
	r.Use(Authorization(authz, "user:read"))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(rr, req)
	if rr.Header().Get("X-Content-Type-Options") == "" || rr.Header().Get("X-Frame-Options") == "" || rr.Header().Get("X-XSS-Protection") == "" || rr.Header().Get("Content-Security-Policy") == "" {
		t.Fatalf("security headers missing")
	}
}
func TestAuth_Unauthorized_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := makeCfg()
	r := gin.New()
	r.Use(Auth("secret", cfg))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d want 401", rr.Code)
	}
}

func TestAuth_Authorized_WithValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := makeCfg()
	secret := "s"
	claims := jwtutil.Claims{UserID: uuid.New(), TenantID: "t1"}
	token, _ := jwtutil.GenerateAccessWith(secret, time.Minute, claims, cfg.JWTIssuer, cfg.JWTAudience)
	r := gin.New()
	r.Use(Auth(secret, cfg))
	r.GET("/", func(c *gin.Context) { c.Status(204) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("code=%d want 204", rr.Code)
	}
}

func TestAuth_InvalidIssuer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := makeCfg()
	secret := "s"
	claims := jwtutil.Claims{UserID: uuid.New(), TenantID: "t1"}
	token, _ := jwtutil.GenerateAccessWith(secret, time.Minute, claims, "wrong", cfg.JWTAudience)
	r := gin.New()
	r.Use(Auth(secret, cfg))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d want 401", rr.Code)
	}
}

func TestCORS_AllowsOriginAndOptions(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORS([]string{"*"}))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://x")
	r.ServeHTTP(rr, req)
	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("cors header missing")
	}
	rr2 := httptest.NewRecorder()
	reqOpt := httptest.NewRequest("OPTIONS", "/", nil)
	reqOpt.Header.Set("Origin", "http://x")
	r.ServeHTTP(rr2, reqOpt)
	if rr2.Code != http.StatusNoContent {
		t.Fatalf("options should be 204, got %d", rr2.Code)
	}
}

type fakeLimiter struct {
	count map[string]int64
	fail  bool
}

func (fl *fakeLimiter) Incr(_ context.Context, key string) (int64, error) {
	if fl.fail {
		return 0, errors.New("fail")
	}
	fl.count[key]++
	return fl.count[key], nil
}
func (fl *fakeLimiter) Expire(_ context.Context, _ string, _ time.Duration) error { return nil }

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fl := &fakeLimiter{count: map[string]int64{}}
	r := gin.New()
	r.Use(RateLimit(fl, 1))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	req := httptest.NewRequest("GET", "/", nil)
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req)
	if rr1.Code != 200 {
		t.Fatalf("first should pass, got %d", rr1.Code)
	}
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("second should be limited, got %d", rr2.Code)
	}
}

func TestRateLimit_ErrorPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fl := &fakeLimiter{count: map[string]int64{}, fail: true}
	r := gin.New()
	r.Use(RateLimit(fl, 1))
	r.GET("/", func(c *gin.Context) { c.Status(200) })
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("should be 500 on limiter error, got %d", rr.Code)
	}
}

func TestRateLimitByPolicy_Gin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fl := &fakeLimiter{count: map[string]int64{}}
	r := gin.New()
	// Default read=3, write=2, override auth prefix=1
	r.Use(RateLimitByPolicy(fl, 3, 2, map[string]int{"/api/v1/auth/": 1}))
	r.POST("/api/v1/auth/login", func(c *gin.Context) { c.Status(200) })
	// First POST to auth/login should pass
	req1 := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)
	if rr1.Code != 200 {
		t.Fatalf("first auth POST should pass, got %d", rr1.Code)
	}
	// Second POST to auth/login should be limited due to override=1
	req2 := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("second auth POST should be limited, got %d", rr2.Code)
	}
	// GET to non-auth path should use defaultRead=3
	r.GET("/api/v1/data", func(c *gin.Context) { c.Status(200) })
	reqG := httptest.NewRequest("GET", "/api/v1/data", nil)
	for i := 0; i < 3; i++ {
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, reqG)
		if rr.Code != 200 {
			t.Fatalf("get #%d should pass, got %d", i+1, rr.Code)
		}
	}
	rr4 := httptest.NewRecorder()
	r.ServeHTTP(rr4, reqG)
	if rr4.Code != http.StatusTooManyRequests {
		t.Fatalf("get fourth should be limited, got %d", rr4.Code)
	}
}
