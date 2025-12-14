package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type fakeRedis struct {
	count map[string]int64
}

func (f *fakeRedis) Incr(_ context.Context, key string) *fakeIntCmd {
	f.count[key]++
	return &fakeIntCmd{val: f.count[key]}
}
func (f *fakeRedis) Expire(_ context.Context, _ string, _ time.Duration) *fakeBoolCmd {
	return &fakeBoolCmd{val: true}
}

type fakeIntCmd struct{ val int64 }
func (c *fakeIntCmd) Result() (int64, error) { return c.val, nil }

type fakeBoolCmd struct{ val bool }
func (c *fakeBoolCmd) Result() (bool, error) { return c.val, nil }
func (c *fakeBoolCmd) Err() error            { return nil }

type fakeLimiter struct{ r *fakeRedis }
func (fl *fakeLimiter) Incr(ctx context.Context, key string) (int64, error) {
	return fl.r.Incr(ctx, key).Result()
}
func (fl *fakeLimiter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return fl.r.Expire(ctx, key, ttl).Err()
}

type recordingLimiter struct {
	lastKey string
	fail    bool
	r       *fakeRedis
}
func (rl *recordingLimiter) Incr(ctx context.Context, key string) (int64, error) {
	rl.lastKey = key
	if rl.fail {
		return 0, fmt.Errorf("fail")
	}
	return rl.r.Incr(ctx, key).Result()
}
func (rl *recordingLimiter) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return rl.r.Expire(ctx, key, ttl).Err()
}

func TestCORS(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := CORS([]string{"*"}, fn)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://x")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("cors header missing")
	}
}

func TestCORSPreflight(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := CORS([]string{"http://x"}, fn)
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://x")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("preflight should return 204")
	}
}

func TestLogging(t *testing.T) {
	l, _ := zap.NewProduction()
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	h := RequestID(Logging(l, fn))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != 204 {
		t.Fatalf("code=%d want 204", rr.Code)
	}
	if rr.Header().Get("X-Request-ID") == "" {
		t.Fatalf("request id should be set")
	}
}

func TestAuthUnauthorized(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := Auth("secret", fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("should be unauthorized without header")
	}
}

func TestAuthWrongScheme(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := Auth("secret", fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Basic abc")
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("should be unauthorized with wrong scheme")
	}
}

func TestRecover(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic("x") })
	h := Recover(fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("should recover panic and return 500")
	}
}
func TestAuth(t *testing.T) {
	secret := "s"
	token, _ := jwtutil.GenerateAccess(secret, time.Minute, jwtutil.Claims{UserID: uuid.New(), TenantID: "t"})
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := Auth(secret, fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d want 200", rr.Code)
	}
}

func TestAuthInvalidToken(t *testing.T) {
	validSecret := "s"
	token, _ := jwtutil.GenerateAccess(validSecret, time.Minute, jwtutil.Claims{UserID: uuid.New(), TenantID: "t"})
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := Auth("other", fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("should be unauthorized with invalid token")
	}
}

func TestAuthEmptyBearerToken(t *testing.T) {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := Auth("secret", fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer ")
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("should be unauthorized with empty bearer token")
	}
}

func TestAuthWithIssuerAudience(t *testing.T) {
	secret := "s"
	issuer := "sisfo-akademik"
	audience := "api"
	claims := jwtutil.Claims{UserID: uuid.New(), TenantID: "t"}
	token, _ := jwtutil.GenerateAccessWith(secret, time.Minute, claims, issuer, audience)
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := AuthWith(secret, issuer, audience, fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("code=%d want 200", rr.Code)
	}
}

func TestAuthWithInvalidIssuer(t *testing.T) {
	secret := "s"
	issuer := "sisfo-akademik"
	audience := "api"
	claims := jwtutil.Claims{UserID: uuid.New(), TenantID: "t"}
	token, _ := jwtutil.GenerateAccessWith(secret, time.Minute, claims, issuer, audience)
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := AuthWith(secret, "wrong", audience, fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("should be unauthorized with wrong issuer")
	}
}

func TestRateLimit(t *testing.T) {
	f := &fakeRedis{count: map[string]int64{}}
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := RateLimit(&fakeLimiter{r: f}, 1, fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("first request should pass")
	}
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req)
	if rr2.Code != 429 {
		t.Fatalf("second request should be limited, got %d", rr2.Code)
	}
}

func TestRateLimitErrorPath(t *testing.T) {
	f := &fakeRedis{count: map[string]int64{}}
	rl := &recordingLimiter{fail: true, r: f}
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := RateLimit(rl, 1, fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("should be 500 on limiter error")
	}
}

func TestClientIPForwardedHeader(t *testing.T) {
	f := &fakeRedis{count: map[string]int64{}}
	rl := &recordingLimiter{r: f}
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := RateLimit(rl, 10, fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/path", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	h.ServeHTTP(rr, req)
	if !strings.Contains(rl.lastKey, "ratelimit:1.2.3.4:/path") {
		t.Fatalf("key should include forwarded IP and path")
	}
}

func TestRateLimitByPrefix(t *testing.T) {
	f := &fakeRedis{count: map[string]int64{}}
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := RateLimitByPrefix(&fakeLimiter{r: f}, 10, map[string]int{"/api/v1/auth/": 2}, fn)
	req := httptest.NewRequest("GET", "/api/v1/auth/login", nil)
	req.Header.Set("X-Forwarded-For", "9.9.9.9")
	rr1 := httptest.NewRecorder()
	h.ServeHTTP(rr1, req)
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req)
	rr3 := httptest.NewRecorder()
	h.ServeHTTP(rr3, req)
	if rr3.Code != http.StatusTooManyRequests {
 	t.Fatalf("code=%d want 429", rr3.Code)
 }
}

func TestRateLimitByPolicy(t *testing.T) {
	f := &fakeRedis{count: map[string]int64{}}
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := RateLimitByPolicy(&fakeLimiter{r: f}, 3, 2, map[string]int{"/api/v1/auth/": 1}, fn)
	// GET should use defaultRead=3 unless prefix override
	reqGet := httptest.NewRequest("GET", "/api/v1/data", nil)
	reqGet.Header.Set("X-Forwarded-For", "7.7.7.7")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, reqGet)
	h.ServeHTTP(rr, reqGet)
	h.ServeHTTP(rr, reqGet)
	if rr.Code != 200 {
		t.Fatalf("get third should pass, got %d", rr.Code)
	}
	rr4 := httptest.NewRecorder()
	h.ServeHTTP(rr4, reqGet)
	if rr4.Code != http.StatusTooManyRequests {
		t.Fatalf("get fourth should be limited, got %d", rr4.Code)
	}
	// POST to /api/v1/auth/* should be override=1
	reqPost := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	reqPost.Header.Set("X-Forwarded-For", "7.7.7.7")
	rrp1 := httptest.NewRecorder()
	h.ServeHTTP(rrp1, reqPost)
	rrp2 := httptest.NewRecorder()
	h.ServeHTTP(rrp2, reqPost)
	if rrp2.Code != http.StatusTooManyRequests {
		t.Fatalf("post second should be limited, got %d", rrp2.Code)
	}
}
