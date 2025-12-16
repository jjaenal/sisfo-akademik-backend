package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDBAuditRetention(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		url = "postgres://dev:dev@localhost:55432/devdb?sslmode=disable"
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Skip("no db available:", err)
	}
	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Skip("db connect failed:", err)
	}
	return db
}

func ensureMigrationsRetention(t *testing.T, db *pgxpool.Pool) {
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

func TestAuditRepo_CleanupOlderThan(t *testing.T) {
	db := testDBAuditRetention(t)
	ensureMigrationsRetention(t, db)
	repo := NewAuditRepo(db)
	users := NewUsersRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant, Email: "retention@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Skip("cannot create user:", err)
	}
	old := time.Now().UTC().Add(-91 * 24 * time.Hour)
	newer := time.Now().UTC().Add(-30 * 24 * time.Hour)
	uid := u.ID
	rid := uuid.New()
	if _, err := db.Exec(context.Background(), `
		INSERT INTO audit_logs (id, tenant_id, user_id, action, resource_type, resource_id, new_values, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, '{}'::jsonb, $7)
	`, uuid.New(), tenant, uid, "old.action", "user", rid, old); err != nil {
		t.Fatalf("insert old err: %v", err)
	}
	if _, err := db.Exec(context.Background(), `
		INSERT INTO audit_logs (id, tenant_id, user_id, action, resource_type, resource_id, new_values, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, '{}'::jsonb, $7)
	`, uuid.New(), tenant, uid, "new.action", "user", rid, newer); err != nil {
		t.Fatalf("insert new err: %v", err)
	}
	cutoff := time.Now().UTC().Add(-90 * 24 * time.Hour)
	err = nil
	_, err = repo.CleanupOlderThan(context.Background(), cutoff)
	if err != nil {
		t.Fatalf("cleanup err: %v", err)
	}
	var cntOld, cntNew int
	_ = db.QueryRow(context.Background(), `SELECT COUNT(1) FROM audit_logs WHERE tenant_id=$1 AND action='old.action'`, tenant).Scan(&cntOld)
	_ = db.QueryRow(context.Background(), `SELECT COUNT(1) FROM audit_logs WHERE tenant_id=$1 AND action='new.action'`, tenant).Scan(&cntNew)
	if cntOld != 0 || cntNew != 1 {
		t.Fatalf("unexpected counts old=%d new=%d", cntOld, cntNew)
	}
}
