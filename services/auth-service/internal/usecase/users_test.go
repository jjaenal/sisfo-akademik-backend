package usecase

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
)

func testDBUC(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsUC(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	_, _ = db.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(context.Background(), string(b1))
	}
}

func TestUsersUsecase_RegisterGetListUpdateDelete(t *testing.T) {
	db := testDBUC(t)
	ensureMigrationsUC(t, db)
	repo := repository.NewUsersRepo(db)
	uc := NewUsers(repo)
	tenant := "t-" + uuid.NewString()
	// invalid password
	if _, err := uc.Register(context.Background(), UserRegisterInput{
		TenantID: tenant, Email: "a@b.c", Password: "short",
	}); err == nil {
		t.Fatalf("expected error for short password")
	}
	// register
	u, err := uc.Register(context.Background(), UserRegisterInput{
		TenantID: tenant, Email: "uc1@test.local", Password: "password123",
	})
	if err != nil {
		t.Fatalf("register err: %v", err)
	}
	// get
	got, err := uc.Get(context.Background(), u.ID)
	if err != nil || got.Email != "uc1@test.local" {
		t.Fatalf("get err: %v got=%+v", err, got)
	}
	// list
	items, total, err := uc.List(context.Background(), tenant, 0, -1)
	if err != nil || total < 1 || len(items) < 1 {
		t.Fatalf("list err: %v total=%d len=%d", err, total, len(items))
	}
	// update
	newEmail := "uc1b@test.local"
	upd, err := uc.Update(context.Background(), u.ID, UserUpdateInput{
		Email: &newEmail,
		IsActive: func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if upd.Email != newEmail {
		t.Fatalf("email not updated: %s", upd.Email)
	}
	// delete
	if err := uc.Delete(context.Background(), u.ID); err != nil {
		t.Fatalf("delete err: %v", err)
	}
}
