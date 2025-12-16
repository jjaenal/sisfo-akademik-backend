package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
)

func makeCfg() config.Config {
	return config.Config{
		Env:                "test",
		ServiceName:        "api-gateway",
		HTTPPort:           0,
		PostgresURL:        "postgres://dev:dev@localhost:55432/devdb?sslmode=disable",
		RedisAddr:          "localhost:6379",
		RabbitURL:          "amqp://guest:guest@localhost:5672/",
		JWTAccessSecret:    "access",
		JWTRefreshSecret:   "refresh",
		JWTAccessTTL:       15 * time.Minute,
		JWTRefreshTTL:      7 * 24 * time.Hour,
		JWTIssuer:          "sisfo-akademik",
		JWTAudience:        "api",
		CORSAllowedOrigins: []string{"*"},
		RateLimitPerMinute: 60,
	}
}

func TestParseUpstreams_ListAndSingle(t *testing.T) {
	_ = os.Setenv("APP_UPSTREAM_TEST_URLS", "http://a.local,http://b.local")
	got := parseUpstreams("APP_UPSTREAM_TEST")
	if len(got) != 2 {
		t.Fatalf("want 2 got %d", len(got))
	}
	_ = os.Unsetenv("APP_UPSTREAM_TEST_URLS")
	_ = os.Setenv("APP_UPSTREAM_TEST_URL", "http://c.local")
	got = parseUpstreams("APP_UPSTREAM_TEST")
	if len(got) != 1 || got[0] != "http://c.local" {
		t.Fatalf("fallback single url failed")
	}
	_ = os.Unsetenv("APP_UPSTREAM_TEST_URL")
}

func TestGatewayHealth_NoUpstreams(t *testing.T) {
	_ = os.Unsetenv("APP_UPSTREAM_AUTH_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_AUTH_URL")
	_ = os.Unsetenv("APP_UPSTREAM_ACADEMIC_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_ACADEMIC_URL")
	_ = os.Unsetenv("APP_UPSTREAM_ATTENDANCE_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_ATTENDANCE_URL")
	_ = os.Unsetenv("APP_UPSTREAM_ASSESSMENT_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_ASSESSMENT_URL")
	_ = os.Unsetenv("APP_UPSTREAM_ADMISSION_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_ADMISSION_URL")
	_ = os.Unsetenv("APP_UPSTREAM_FINANCE_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_FINANCE_URL")
	_ = os.Unsetenv("APP_UPSTREAM_NOTIFICATION_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_NOTIFICATION_URL")
	_ = os.Unsetenv("APP_UPSTREAM_FILE_URLS")
	_ = os.Unsetenv("APP_UPSTREAM_FILE_URL")
	cfg := makeCfg()
	h := gatewayHealthHandler(cfg)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gateway/health", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d", rr.Code)
	}
	var resp struct {
		Services map[string]struct {
			Up     bool   `json:"up"`
			Status int    `json:"status"`
			Error  string `json:"error"`
		} `json:"services"`
	}
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp.Services["auth"].Up {
		t.Fatalf("expected auth up=false")
	}
}

func TestRegisterRoutes_AuthProxy(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/auth/test" {
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(404)
	}))
	defer up.Close()
	_ = os.Setenv("APP_UPSTREAM_AUTH_URLS", up.URL)
	cfg := makeCfg()
	mux := http.NewServeMux()
	registerRoutes(mux, cfg)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/test", nil)
	mux.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("proxy code=%d want 200", rr.Code)
	}
	_ = os.Unsetenv("APP_UPSTREAM_AUTH_URLS")
}
