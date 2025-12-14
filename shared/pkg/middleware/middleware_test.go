package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwtutil "shared/pkg/jwt"

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
	h := Logging(l, fn)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != 204 {
		t.Fatalf("code=%d want 204", rr.Code)
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
