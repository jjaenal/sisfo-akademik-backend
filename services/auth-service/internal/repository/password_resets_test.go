package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDBPR(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsPR(t *testing.T, db *pgxpool.Pool) {
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
}

func TestPasswordResetRepo_Create_Find_MarkUsed_Expired(t *testing.T) {
	db := testDBPR(t)
	ensureMigrationsPR(t, db)
	users := NewUsersRepo(db)
	pr := NewPasswordResetRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "pr1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	now := time.Now().UTC()
	rec, err := pr.Create(context.Background(), tenant, u.ID, "hash-"+uuid.NewString(), now.Add(10*time.Minute))
	if err != nil {
		t.Fatalf("create reset err: %v", err)
	}
	got, err := pr.FindValidByTokenHash(context.Background(), rec.TokenHash, now.Add(1*time.Minute))
	if err != nil || got == nil || got.UserID != u.ID {
		t.Fatalf("find valid err=%v got=%+v", err, got)
	}
	if err := pr.MarkUsed(context.Background(), rec.ID, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("mark used err: %v", err)
	}
	if _, err := pr.FindValidByTokenHash(context.Background(), rec.TokenHash, now.Add(3*time.Minute)); err == nil {
		t.Fatalf("expected not found after used")
	}
	rec2, err := pr.Create(context.Background(), tenant, u.ID, "hash-"+uuid.NewString(), now.Add(1*time.Second))
	if err != nil {
		t.Fatalf("create reset2 err: %v", err)
	}
	time.Sleep(1200 * time.Millisecond)
	if _, err := pr.FindValidByTokenHash(context.Background(), rec2.TokenHash, time.Now().UTC()); err == nil {
		t.Fatalf("expected expired token to be invalid")
	}
}

func TestPasswordResetRepo_ExpiredWithoutSleep(t *testing.T) {
	db := testDBPR(t)
	ensureMigrationsPR(t, db)
	users := NewUsersRepo(db)
	pr := NewPasswordResetRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "pr2@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	now := time.Now().UTC()
	rec, err := pr.Create(context.Background(), tenant, u.ID, "hash-"+uuid.NewString(), now.Add(-1*time.Minute))
	if err != nil {
		t.Fatalf("create expired reset err: %v", err)
	}
	if _, err := pr.FindValidByTokenHash(context.Background(), rec.TokenHash, now); err == nil {
		t.Fatalf("expected expired token to be invalid without sleep")
	}
}

func TestPasswordResetRepo_MarkUsed_Idempotent(t *testing.T) {
	db := testDBPR(t)
	ensureMigrationsPR(t, db)
	users := NewUsersRepo(db)
	pr := NewPasswordResetRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "pr3@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	now := time.Now().UTC()
	rec, err := pr.Create(context.Background(), tenant, u.ID, "hash-"+uuid.NewString(), now.Add(10*time.Minute))
	if err != nil {
		t.Fatalf("create reset err: %v", err)
	}
	if err := pr.MarkUsed(context.Background(), rec.ID, now.Add(1*time.Minute)); err != nil {
		t.Fatalf("mark used err: %v", err)
	}
	if err := pr.MarkUsed(context.Background(), rec.ID, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("mark used second err: %v", err)
	}
	if _, err := pr.FindValidByTokenHash(context.Background(), rec.TokenHash, now.Add(3*time.Minute)); err == nil {
		t.Fatalf("expected not found after used twice")
	}
}
