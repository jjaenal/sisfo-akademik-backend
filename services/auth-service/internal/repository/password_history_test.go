package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDBPH(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsPH(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b1))
	}
	b4, err := os.ReadFile("../../migrations/004_password_history.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b4))
	}
}

func TestPasswordHistoryRepo_Add_Recent_Prune(t *testing.T) {
	db := testDBPH(t)
	ensureMigrationsPH(t, db)
	users := NewUsersRepo(db)
	ph := NewPasswordHistoryRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "ph1@test.local",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	now := time.Now().UTC()
	for i := 0; i < 7; i++ {
		if err := ph.Add(context.Background(), u.ID, "hash-"+uuid.NewString(), now.Add(time.Duration(i)*time.Second)); err != nil {
			t.Fatalf("add history err: %v", err)
		}
	}
	items, err := ph.Recent(context.Background(), u.ID, 2)
	if err != nil || len(items) != 2 {
		t.Fatalf("recent err=%v len=%d", err, len(items))
	}
	if !items[0].CreatedAt.After(items[1].CreatedAt) && !items[0].CreatedAt.Equal(items[1].CreatedAt) {
		t.Fatalf("recent should be DESC ordered")
	}
	if err := ph.Prune(context.Background(), u.ID, 5); err != nil {
		t.Fatalf("prune err: %v", err)
	}
	var cnt int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.QueryRow(ctx, `SELECT COUNT(1) FROM password_history WHERE user_id=$1`, u.ID).Scan(&cnt); err != nil {
		t.Fatalf("count err: %v", err)
	}
	if cnt != 5 {
		t.Fatalf("expected 5 kept, got %d", cnt)
	}
}

