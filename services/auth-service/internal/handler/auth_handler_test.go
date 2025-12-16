package handler

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	b3, err := os.ReadFile("../../migrations/003_password_resets.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b3))
	}
	b4, err := os.ReadFile("../../migrations/004_password_history.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b4))
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
		LockoutThreshold:   5,
		LockoutTTL:         15 * time.Minute,
		FailWindowTTL:      15 * time.Minute,
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
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)

	r := gin.New()
	h.Register(r)

	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "login1@test.local",
		Password: "password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}

	body := map[string]string{"tenant_id": tenant, "email": "login1@test.local", "password": "password123!"}
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

func TestAuthHandler_AccountLockout_Threshold_And_Success_Reset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "lock1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		body := map[string]string{"tenant_id": tenant, "email": "lock1@test.local", "password": "wrong"}
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if i < 4 {
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("attempt %d expected 401 got %d", i+1, w.Code)
			}
		} else {
			if w.Code != http.StatusForbidden {
				t.Fatalf("attempt %d expected 403 got %d", i+1, w.Code)
			}
		}
	}
	w := httptest.NewRecorder()
	bodyOK := map[string]string{"tenant_id": tenant, "email": "lock1@test.local", "password": "Password123!"}
	bOK, _ := json.Marshal(bodyOK)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bOK))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when locked got %d body=%s", w.Code, w.Body.String())
	}
	_ = redisutil.Del(context.Background(), redis.Raw(), "lockout:"+tenant+":lock1@test.local")
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bOK))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 after unlock got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_AccountLockout_ManualLock_PreventsLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "lock2@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	_ = redisutil.Set(context.Background(), redis.Raw(), "lockout:"+tenant+":lock2@test.local", "1", time.Minute)
	w := httptest.NewRecorder()
	body := map[string]string{"tenant_id": tenant, "email": "lock2@test.local", "password": "Password123!"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d body=%s", w.Code, w.Body.String())
	}
}
func TestAuthHandler_Me_Unauthorized_When_NoClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)

	r := gin.New()
	protected := r.Group("/")
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", w.Code, w.Body.String())
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
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
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
		TenantID: tenant, Email: "login2@test.local", Password: "password123!",
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
	bodyInactive := map[string]string{"tenant_id": tenant, "email": "login2@test.local", "password": "password123!"}
	bInactive, _ := json.Marshal(bodyInactive)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bInactive))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", w.Code)
	}

	// login with unknown user -> expect 401
	w = httptest.NewRecorder()
	bodyUnknown := map[string]string{"tenant_id": tenant, "email": "unknown@test.local", "password": "password123"}
	bUnknown, _ := json.Marshal(bodyUnknown)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(bUnknown))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
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
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
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
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
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

func TestAuthHandler_Logout_EmptyBody_200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Fatalf("logout empty body unexpected status=%d body=%s", w.Code, w.Body.String())
	}
}
func TestAuthHandler_ChangePassword_Reuse_Prevented(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	protected := r.Group("/")
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "reuse1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"Password123!","new_password":"Password123!"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ResetPassword_Reuse_Prevented(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "reuse2@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	w := httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "reuse2@test.local"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var fresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &fresp)
	token := fresp.Data["reset_token"]
	if token == "" {
		t.Fatalf("missing reset_token")
	}
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": token, "password": "Password123!"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}
func TestAuthHandler_Forgot_Reset_And_Login_With_New_Password(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)

	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "forgot1@test.local",
		Password: "password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	// forgot
	w := httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "forgot1@test.local"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var fresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &fresp)
	token := fresp.Data["reset_token"]
	if token == "" {
		t.Fatalf("missing reset_token in test env")
	}
	// reset
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": token, "password": "NewPass123!"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("reset status=%d body=%s", w.Code, w.Body.String())
	}
	// login with new password
	w = httptest.NewRecorder()
	lb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "forgot1@test.local", "password": "NewPass123!"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login new pass status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ChangePassword_Fails_When_Wrong_Old_Then_Succeeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	protected := r.Group("/")
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "change1@test.local",
		Password: "password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	// wrong old
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"wrong","new_password":"NewPass123!"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", w.Code, w.Body.String())
	}
	// correct old
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"password123!","new_password":"NewPass123!"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("change status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Refresh_Blacklists_Old_Token(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "refreshbl@test.local", Password: "password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	oldRefresh, _ := jwtutil.GenerateRefreshWith(cfg.JWTRefreshSecret, time.Minute, u.ID.String(), cfg.JWTIssuer, cfg.JWTAudience)
	w := httptest.NewRecorder()
	bOK, _ := json.Marshal(map[string]string{"refresh_token": oldRefresh})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(bOK))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first refresh status=%d body=%s", w.Code, w.Body.String())
	}
	// Try using the old refresh again -> should be unauthorized (revoked)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(bOK))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 revoked got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_Logout_Blacklists_Provided_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "logoutbl@test.local", Password: "password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	refresh, _ := jwtutil.GenerateRefreshWith(cfg.JWTRefreshSecret, time.Minute, u.ID.String(), cfg.JWTIssuer, cfg.JWTAudience)
	// Logout to blacklist the provided refresh
	w := httptest.NewRecorder()
	lb, _ := json.Marshal(map[string]string{"refresh_token": refresh})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("logout status=%d body=%s", w.Code, w.Body.String())
	}
	// Try refreshing with the same token -> should be unauthorized
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 revoked got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ResetPassword_InvalidToken_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	// Call reset without previous forgot (no redis key set)
	w := httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": "nonexistent-token", "password": "NewPass123!"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_AuditLogs_Actions_Recorded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "audit1@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	// login
	w := httptest.NewRecorder()
	lb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "audit1@test.local", "password": "Password123!"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", w.Code, w.Body.String())
	}
	var lresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &lresp)
	refresh := lresp.Data["refresh_token"]
	// refresh
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"refresh_token": refresh})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("refresh status=%d body=%s", w.Code, w.Body.String())
	}
	// logout
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("logout status=%d body=%s", w.Code, w.Body.String())
	}
	// change-password
	r2 := gin.New()
	protected := r2.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"Password123!","new_password":"NewPass123!"}`)))
	req.Header.Set("Content-Type", "application/json")
	r2.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("change status=%d body=%s", w.Code, w.Body.String())
	}
	// forgot + reset
	w = httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "audit1@test.local"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var fresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &fresp)
	token := fresp.Data["reset_token"]
	w = httptest.NewRecorder()
	rb2, _ := json.Marshal(map[string]string{"token": token, "password": "NewPass1234!"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb2))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("reset status=%d body=%s", w.Code, w.Body.String())
	}
	// Verify audit_logs entries exist for actions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	check := func(action string) {
		var cnt int
		err := db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action=$1`, action).Scan(&cnt)
		if err != nil || cnt < 1 {
			t.Fatalf("expected audit log for %s got cnt=%d err=%v", action, cnt, err)
		}
	}
	check("auth.login")
	check("auth.refresh")
	check("auth.logout")
	check("auth.change_password")
	check("auth.forgot_password")
	check("auth.reset_password")
}

func TestAuthHandler_ChangePassword_Fails_On_Weak_Password(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	protected := r.Group("/")
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "weak1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"Password123!","new_password":"weakpass"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ChangePassword_Fails_On_Common_Password(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	protected := r.Group("/")
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "common1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"Password123!","new_password":"p@ssw0rd"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ChangePassword_Fails_On_Short_Password(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	protected := r.Group("/")
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "short1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"Password123!","new_password":"short"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ResetPassword_Fails_On_Weak_Password(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "weak2@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	w := httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "weak2@test.local"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var fresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &fresp)
	token := fresp.Data["reset_token"]
	if token == "" {
		t.Fatalf("missing reset_token")
	}
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": token, "password": "weakpass"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ResetPassword_Fails_On_Short_Password(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "short2@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	w := httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "short2@test.local"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var fresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &fresp)
	token := fresp.Data["reset_token"]
	if token == "" {
		t.Fatalf("missing reset_token")
	}
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": token, "password": "short"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}
func TestAuthHandler_ResetPassword_Fails_On_Common_Password(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	_, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "common2@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	w := httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "common2@test.local"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var fresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &fresp)
	token := fresp.Data["reset_token"]
	if token == "" {
		t.Fatalf("missing reset_token")
	}
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": token, "password": "p@ssw0rd"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ResetPassword_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ResetPassword_FallbackRedis_Succeeds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "fallback1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	// generate token manually and set redis fallback key (no DB record)
	token := uuid.NewString()
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])
	if err := redisutil.Set(context.Background(), redis.Raw(), "password_reset:th:"+tokenHash, u.ID.String(), 30*time.Minute); err != nil {
		t.Fatalf("redis set err: %v", err)
	}
	// call reset -> should succeed via fallback
	w := httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": token, "password": "NewPass123!"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("reset via fallback status=%d body=%s", w.Code, w.Body.String())
	}
	// verify audit log exists
	var cnt int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action='auth.reset_password' AND resource_id=$1`, u.ID).Scan(&cnt); err != nil || cnt < 1 {
		t.Fatalf("expected audit log entry got cnt=%d err=%v", cnt, err)
	}
}
func TestAuthHandler_ChangePassword_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	protected := r.Group("/")
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "invalidjson1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ChangePassword_Reuse_Prevented_By_History(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	protected := r.Group("/")
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "reusehist1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"Password123!","new_password":"NewPass123!"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first change status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"NewPass123!","new_password":"Password123!"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 reuse via history got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthHandler_ForgotPassword_UnknownUser_GenericSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	usersRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	h := NewAuthHandler(usersRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	r := gin.New()
	h.Register(r)
	tenant := "t-" + uuid.NewString()
	w := httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "notfound@test.local"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var resp struct{ Data map[string]any `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if _, ok := resp.Data["message"]; !ok {
		t.Fatalf("expected generic success message")
	}
}
