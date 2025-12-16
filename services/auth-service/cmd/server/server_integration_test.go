package main

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
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/middleware"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
)

func testDB(t *testing.T) *pgxpool.Pool {
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

func ensureMigrations(t *testing.T, db *pgxpool.Pool) {
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

func TestServer_Integration_RateLimit_Login_429(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)
	cfg := makeCfg()
	cfg.RateLimitPerMinute = 1
	redis := redisutil.New(cfg.RedisAddr)
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	r.Use(middleware.RateLimit(limiter, cfg.RateLimitPerMinute))
	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	authHandler.Register(r)
	tenant := "t-" + uuid.NewString()
	_, _ = authRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "rl1@test.local", Password: "Password123!",
	})
	body := map[string]string{"tenant_id": tenant, "email": "rl1@test.local", "password": "Password123!"}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first login status=%d body=%s", w.Code, w.Body.String())
	}
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 got %d body=%s", w2.Code, w2.Body.String())
	}
}

func makeCfg() config.Config {
	return config.Config{
		Env:                "test",
		ServiceName:        "auth-service",
		HTTPPort:           0,
		PostgresURL:        "postgres://dev:dev@localhost:55432/devdb?sslmode=disable",
		RedisAddr:          "localhost:6379",
		RabbitURL:          "amqp://guest:guest@localhost:5672/",
		JWTAccessSecret:    "access-secret",
		JWTRefreshSecret:   "refresh-secret",
		JWTAccessTTL:       15 * time.Minute,
		JWTRefreshTTL:      7 * 24 * time.Hour,
		JWTIssuer:          "sisfo-akademik",
		JWTAudience:        "api",
		CORSAllowedOrigins: []string{"*"},
		RateLimitPerMinute: 1000,
	}
}

func TestServer_Integration_AuthFlows(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)

	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	r.Use(middleware.RateLimit(limiter, cfg.RateLimitPerMinute))

	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	authHandler.Register(r)

	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTAccessSecret, cfg))
	authHandler.RegisterProtected(protected)
	usersUC := usecase.NewUsers(authRepo)
	usersHandler := handler.NewUsersHandler(usersUC)
	usersHandler.RegisterProtected(protected)
	rolesRepo := repository.NewRolesRepo(db)
	rolesUC := usecase.NewRoles(authRepo, rolesRepo)
	rolesHandler := handler.NewRolesHandler(rolesUC)
	rolesHandler.RegisterProtected(protected)

	tenant := "t-" + uuid.NewString()
	initClaims := jwtutil.Claims{TenantID: tenant, UserID: uuid.New()}
	initAccess, _ := jwtutil.GenerateAccessWith(cfg.JWTAccessSecret, time.Minute, initClaims, cfg.JWTIssuer, cfg.JWTAudience)

	w := httptest.NewRecorder()
	reqBody := map[string]string{"email": "it1@test.local", "password": "Password123!"}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+initAccess)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("users create status=%d body=%s", w.Code, w.Body.String())
	}
	var created struct {
		Data map[string]any `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	uidStr, _ := created.Data["id"].(string)
	if uidStr == "" {
		t.Fatalf("missing id: %s", w.Body.String())
	}

	w = httptest.NewRecorder()
	loginBody := map[string]string{"tenant_id": tenant, "email": "it1@test.local", "password": "Password123!"}
	lb, _ := json.Marshal(loginBody)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", w.Code, w.Body.String())
	}
	var lresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &lresp)
	access := lresp.Data["access_token"]
	refresh := lresp.Data["refresh_token"]
	if access == "" || refresh == "" {
		t.Fatalf("missing tokens")
	}

	w = httptest.NewRecorder()
	newPass := `{"old_password":"Password123!","new_password":"NewPass123!"}`
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(newPass)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+access)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("change status=%d body=%s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "it1@test.local"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewReader(fb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("forgot status=%d body=%s", w.Code, w.Body.String())
	}
	var fresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &fresp)
	rst := fresp.Data["reset_token"]
	if rst == "" {
		t.Fatalf("missing reset_token")
	}

	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"token": rst, "password": "NewPass1234!"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("reset status=%d body=%s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	lb2, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "it1@test.local", "password": "NewPass1234!"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb2))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login new status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestServer_Integration_Refresh_Blacklist_Old(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	r.Use(middleware.RateLimit(limiter, cfg.RateLimitPerMinute))
	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	authHandler.Register(r)
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTAccessSecret, cfg))
	authHandler.RegisterProtected(protected)
	tenant := "t-" + uuid.NewString()
	claims := jwtutil.Claims{TenantID: tenant, UserID: uuid.New()}
	access, _ := jwtutil.GenerateAccessWith(cfg.JWTAccessSecret, time.Minute, claims, cfg.JWTIssuer, cfg.JWTAudience)
	_, _ = authRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "it2@test.local", Password: "Password123!",
	})
	w := httptest.NewRecorder()
	lb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "it2@test.local", "password": "Password123!"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", w.Code, w.Body.String())
	}
	var lresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &lresp)
	refresh := lresp.Data["refresh_token"]
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"refresh_token": refresh})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("refresh status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 revoked got %d body=%s", w.Code, w.Body.String())
	}
	_ = access
}

func TestServer_Integration_Logout_Blacklists_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	r.Use(middleware.RateLimit(limiter, cfg.RateLimitPerMinute))
	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	authHandler.Register(r)
	tenant := "t-" + uuid.NewString()
	_, _ = authRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "it3@test.local", Password: "Password123!",
	})
	w := httptest.NewRecorder()
	lb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "it3@test.local", "password": "Password123!"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", w.Code, w.Body.String())
	}
	var lresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &lresp)
	refresh := lresp.Data["refresh_token"]
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"refresh_token": refresh})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("logout status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 revoked got %d body=%s", w.Code, w.Body.String())
	}
}

func TestServer_Integration_AuditLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	r.Use(middleware.RateLimit(limiter, cfg.RateLimitPerMinute))
	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	authHandler.Register(r)
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTAccessSecret, cfg))
	authHandler.RegisterProtected(protected)
	usersUC := usecase.NewUsers(authRepo)
	usersHandler := handler.NewUsersHandler(usersUC)
	usersHandler.RegisterProtected(protected)
	tenant := "t-" + uuid.NewString()
	claims := jwtutil.Claims{TenantID: tenant, UserID: uuid.New()}
	access, _ := jwtutil.GenerateAccessWith(cfg.JWTAccessSecret, time.Minute, claims, cfg.JWTIssuer, cfg.JWTAudience)
	w := httptest.NewRecorder()
	cb, _ := json.Marshal(map[string]string{"email": "al1@test.local", "password": "Password123!"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(cb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+access)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("users create status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	lb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "al1@test.local", "password": "Password123!"})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", w.Code, w.Body.String())
	}
	var lresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &lresp)
	access2 := lresp.Data["access_token"]
	refresh := lresp.Data["refresh_token"]
	w = httptest.NewRecorder()
	rb, _ := json.Marshal(map[string]string{"refresh_token": refresh})
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("refresh status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader(rb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("logout status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/change-password", bytes.NewReader([]byte(`{"old_password":"Password123!","new_password":"NewPass123!"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+access2)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("change status=%d body=%s", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	fb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "al1@test.local"})
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cnt int
	err := db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action='auth.login'`).Scan(&cnt)
	if err != nil || cnt < 1 {
		t.Fatalf("audit auth.login missing cnt=%d err=%v", cnt, err)
	}
	err = db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action='auth.refresh'`).Scan(&cnt)
	if err != nil || cnt < 1 {
		t.Fatalf("audit auth.refresh missing cnt=%d err=%v", cnt, err)
	}
	err = db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action='auth.logout'`).Scan(&cnt)
	if err != nil || cnt < 1 {
		t.Fatalf("audit auth.logout missing cnt=%d err=%v", cnt, err)
	}
	err = db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action='auth.change_password'`).Scan(&cnt)
	if err != nil || cnt < 1 {
		t.Fatalf("audit auth.change_password missing cnt=%d err=%v", cnt, err)
	}
	err = db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action='auth.forgot_password'`).Scan(&cnt)
	if err != nil || cnt < 1 {
		t.Fatalf("audit auth.forgot_password missing cnt=%d err=%v", cnt, err)
	}
	err = db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE action='auth.reset_password'`).Scan(&cnt)
	if err != nil || cnt < 1 {
		t.Fatalf("audit auth.reset_password missing cnt=%d err=%v", cnt, err)
	}
}

func TestServer_Integration_Me_ReturnsUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)
	cfg := makeCfg()
	redis := redisutil.New(cfg.RedisAddr)
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	r.Use(middleware.RateLimit(limiter, cfg.RateLimitPerMinute))
	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	prRepo := repository.NewPasswordResetRepo(db)
	phRepo := repository.NewPasswordHistoryRepo(db)
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo, prRepo, nil, phRepo)
	authHandler.Register(r)
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTAccessSecret, cfg))
	authHandler.RegisterProtected(protected)
	tenant := "t-" + uuid.NewString()
	created, err := authRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "me1@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed err: %v", err)
	}
	w := httptest.NewRecorder()
	lb, _ := json.Marshal(map[string]string{"tenant_id": tenant, "email": "me1@test.local", "password": "Password123!"})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(lb))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", w.Code, w.Body.String())
	}
	var lresp struct{ Data map[string]string `json:"data"` }
	_ = json.Unmarshal(w.Body.Bytes(), &lresp)
	access := lresp.Data["access_token"]
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("me status=%d body=%s", w.Code, w.Body.String())
	}
	var mresp struct {
		Data map[string]any `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &mresp)
	idStr, _ := mresp.Data["id"].(string)
	if idStr != created.ID.String() {
		t.Fatalf("me id mismatch got=%s want=%s", idStr, created.ID.String())
	}
}

func TestServer_CORS_Preflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := makeCfg()
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	r.POST("/api/v1/auth/login", func(c *gin.Context) { c.Status(200) })
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "http://example.com")
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("preflight code=%d want 204", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("missing cors allow origin")
	}
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatalf("missing cors allow methods")
	}
	if rr.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Fatalf("missing cors allow headers")
	}
}
