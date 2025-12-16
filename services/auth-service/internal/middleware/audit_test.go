package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func testDBAuditMW(t *testing.T) *pgxpool.Pool {
	t.Helper()
	cfg, err := pgxpool.ParseConfig("postgres://dev:dev@localhost:55432/devdb?sslmode=disable")
	if err != nil {
		t.Skip("no db:", err)
	}
	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Skip("db connect failed:", err)
	}
	return db
}

func ensureMigrationsAuditMW(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	_, _ = db.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(context.Background(), string(b1))
	}
	b2, err := os.ReadFile("../../migrations/002_audit_logs.up.sql")
	if err == nil {
		_, _ = db.Exec(context.Background(), string(b2))
	}
}

func TestAuditMiddleware_LogsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuditMW(t)
	ensureMigrationsAuditMW(t, db)
	repo := repository.NewAuditRepo(db)
	r := gin.New()
	r.Use(Audit(repo))
	tenant := "t-" + uuid.NewString()
	users := repository.NewUsersRepo(db)
	u, err := users.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "auditmw@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	r.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	r.GET("/test", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", w.Code)
	}
	var cnt int
	for i := 0; i < 20; i++ {
		_ = db.QueryRow(context.Background(), `SELECT COUNT(1) FROM audit_logs WHERE action='http.request'`).Scan(&cnt)
		if cnt > 0 {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if cnt < 1 {
		t.Fatalf("expected http.request audit log")
	}
}
