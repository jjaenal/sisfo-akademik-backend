package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
)

func testDBAuth(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		url = "postgres://dev:dev@localhost:55432/devdb?sslmode=disable"
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Skip("no db available:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Skip("db connect failed:", err)
	}
	return db
}

func ensureMigrationsAuth(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b1))
	}
	b2, err := os.ReadFile("../../migrations/002_audit_logs.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b2))
	}
}

func makeCfg() config.Config {
	return config.Config{
		Env:               "test",
		ServiceName:       "auth-service",
		HTTPPort:          0,
		PostgresURL:       "postgres://dev:dev@localhost:55432/devdb?sslmode=disable",
		RedisAddr:         "localhost:6379",
		RabbitURL:         "amqp://guest:guest@localhost:5672/",
		JWTAccessSecret:   "access-secret",
		JWTRefreshSecret:  "refresh-secret",
		JWTAccessTTL:      15 * time.Minute,
		JWTRefreshTTL:     7 * 24 * time.Hour,
		JWTIssuer:         "sisfo-akademik",
		JWTAudience:       "api",
		CORSAllowedOrigins: []string{"*"},
		RateLimitPerMinute: 60,
	}
}

func TestAuthHandler_Login_Success_And_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo)

	r := gin.New()
	h.Register(r)

	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "login1@test.local",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}

	body := map[string]string{"tenant_id": tenant, "email": "login1@test.local", "password": "password123"}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", w.Code, w.Body.String())
	}
	var resp struct {
		Data map[string]string `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	access := resp.Data["access_token"]
	refresh := resp.Data["refresh_token"]
	if access == "" || refresh == "" {
		t.Fatalf("missing tokens in response")
	}

	// Use access token to call protected /me
	var claims jwtutil.Claims
	if err := jwtutil.ValidateWith(cfg.JWTAccessSecret, access, &claims, cfg.JWTIssuer, cfg.JWTAudience); err != nil {
		t.Fatalf("validate access failed: %v", err)
	}
	usersRepo.Update(context.Background(), claims.UserID, repository.UpdateUserParams{
		Now: time.Now().UTC(),
	})
	r2 := gin.New()
	protected := r2.Group("/")
	protected.Use(func(c *gin.Context) { c.Set("claims", claims) })
	h.RegisterProtected(protected)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	r2.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("me status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Login_InvalidJSON_WrongPassword_Inactive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo)
	r := gin.New()
	h.Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}

	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "login2@test.local", Password: "password123",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	bodyWrong := map[string]string{"tenant_id": tenant, "email": "login2@test.local", "password": "wrongpass"}
	bWrong, _ := json.Marshal(bodyWrong)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bWrong))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}

	falseVal := false
	_, err = usersRepo.Update(context.Background(), u.ID, repository.UpdateUserParams{
		IsActive: &falseVal,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("deactivate err: %v", err)
	}
	bodyInactive := map[string]string{"tenant_id": tenant, "email": "login2@test.local", "password": "password123"}
	bInactive, _ := json.Marshal(bodyInactive)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bInactive))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", w.Code)
	}
}

func TestAuthHandler_Refresh_InvalidCases_And_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo)
	r := gin.New()
	h.Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader([]byte(`{"refresh_token":"not-a-jwt"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}

	wrongIss, _ := jwtutil.GenerateRefreshWith(cfg.JWTRefreshSecret, time.Minute, uuid.NewString(), "wrong-issuer", cfg.JWTAudience)
	bWrongIss, _ := json.Marshal(map[string]string{"refresh_token": wrongIss})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(bWrongIss))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 invalid issuer got %d", w.Code)
	}

	wrongAud, _ := jwtutil.GenerateRefreshWith(cfg.JWTRefreshSecret, time.Minute, uuid.NewString(), cfg.JWTIssuer, "other-aud")
	bWrongAud, _ := json.Marshal(map[string]string{"refresh_token": wrongAud})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(bWrongAud))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 invalid audience got %d", w.Code)
	}

	notUUID, _ := jwtutil.GenerateRefreshWith(cfg.JWTRefreshSecret, time.Minute, "not-a-uuid", cfg.JWTIssuer, cfg.JWTAudience)
	bNotUUID, _ := json.Marshal(map[string]string{"refresh_token": notUUID})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(bNotUUID))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 invalid subject got %d", w.Code)
	}

	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "refresh1@test.local", Password: "password123",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	okTkn, _ := jwtutil.GenerateRefreshWith(cfg.JWTRefreshSecret, time.Minute, u.ID.String(), cfg.JWTIssuer, cfg.JWTAudience)
	bOK, _ := json.Marshal(map[string]string{"refresh_token": okTkn})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(bOK))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("refresh status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Logout_Always200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo)
	r := gin.New()
	h.Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader([]byte(`{"refresh_token":"not-a-jwt"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("logout status=%d body=%s", w.Code, w.Body.String())
	}
}
