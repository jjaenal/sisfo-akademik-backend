package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
}

func TestUsersRepo_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	tenant := "test-" + uuid.NewString()
	u, err := repo.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "user1@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	got, err := repo.FindByID(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("find err: %v", err)
	}
	if got.Email != "user1@example.com" || got.TenantID != tenant {
		t.Fatalf("unexpected user: %+v", got)
	}
	items, total, err := repo.List(context.Background(), tenant, 10, 0)
	if err != nil || total < 1 || len(items) < 1 {
		t.Fatalf("list err: %v total=%d len=%d", err, total, len(items))
	}
	newEmail := "user1b@example.com"
	upd, err := repo.Update(context.Background(), u.ID, UpdateUserParams{
		Email: &newEmail,
		Now:   time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if upd.Email != newEmail {
		t.Fatalf("email not updated: %+v", upd)
	}
	if err := repo.SoftDelete(context.Background(), u.ID); err != nil {
		t.Fatalf("soft delete err: %v", err)
	}
}
