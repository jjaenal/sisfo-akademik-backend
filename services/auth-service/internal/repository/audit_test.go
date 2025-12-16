package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDBAudit(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsAudit(t *testing.T, db *pgxpool.Pool) {
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

func TestAuditRepo_Log_WritesRecord(t *testing.T) {
	db := testDBAudit(t)
	ensureMigrationsAudit(t, db)
	repo := NewAuditRepo(db)
	tenant := "t-" + uuid.NewString()
	users := NewUsersRepo(db)
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "audit1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	usr := u.ID
	res := uuid.New()
	if err := repo.Log(context.Background(), tenant, &usr, "auth.test", "user", &res, map[string]any{"ok": true}); err != nil {
		t.Fatalf("log err: %v", err)
	}
	var cnt int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE tenant_id=$1 AND action='auth.test' AND resource_id=$2`, tenant, res).Scan(&cnt); err != nil || cnt != 1 {
		t.Fatalf("expected 1 log got %d err=%v", cnt, err)
	}
}

func TestAuditRepo_Log_NilIDs(t *testing.T) {
	db := testDBAudit(t)
	ensureMigrationsAudit(t, db)
	repo := NewAuditRepo(db)
	tenant := "t-" + uuid.NewString()
	if err := repo.Log(context.Background(), tenant, nil, "audit.nil", "system", nil, map[string]any{"x": 1}); err != nil {
		t.Fatalf("log err: %v", err)
	}
	var cnt int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM audit_logs WHERE tenant_id=$1 AND action='audit.nil' AND user_id IS NULL AND resource_id IS NULL`, tenant).Scan(&cnt); err != nil || cnt != 1 {
		t.Fatalf("expected 1 log with nil IDs got %d err=%v", cnt, err)
	}
}

func TestAuditRepo_List_And_Search(t *testing.T) {
	db := testDBAudit(t)
	ensureMigrationsAudit(t, db)
	repo := NewAuditRepo(db)
	tenant := "t-" + uuid.NewString()
	uRepo := NewUsersRepo(db)
	u, err := uRepo.Create(context.Background(), CreateUserParams{
		TenantID: tenant, Email: "auditlist@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	userID := u.ID
	rid := uuid.New()
	if err := repo.Log(context.Background(), tenant, &userID, "auth.login", "user", &rid, map[string]any{"success": true}); err != nil {
		t.Fatalf("log1 err: %v", err)
	}
	if err := repo.Log(context.Background(), tenant, &userID, "auth.refresh", "user", &rid, map[string]any{"success": true}); err != nil {
		t.Fatalf("log2 err: %v", err)
	}
	items, total, err := repo.List(context.Background(), tenant, ListParams{Action: "auth.login"}, 10, 0)
	if err != nil || total < 1 || len(items) < 1 {
		t.Fatalf("list err=%v total=%d len=%d", err, total, len(items))
	}
	sitems, stotal, err := repo.Search(context.Background(), tenant, "refresh", 10, 0)
	if err != nil || stotal < 1 || len(sitems) < 1 {
		t.Fatalf("search err=%v total=%d len=%d", err, stotal, len(sitems))
	}
}
