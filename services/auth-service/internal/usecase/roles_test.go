package usecase

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
)

func testDBRolesUC(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsRolesUC(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	_, _ = db.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(context.Background(), string(b1))
	}
}

func TestRolesUsecase_AssignListUnassign(t *testing.T) {
	db := testDBRolesUC(t)
	ensureMigrationsRolesUC(t, db)
	users := repository.NewUsersRepo(db)
	roles := repository.NewRolesRepo(db)
	uc := NewRoles(users, roles)

	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant,
		Email:    "rolesuc@test.local",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	if _, err := uc.AssignByName(context.Background(), tenant, u.ID, "admin"); err != nil {
		t.Fatalf("assign err: %v", err)
	}
	items, err := uc.List(context.Background(), u.ID)
	if err != nil || len(items) == 0 {
		t.Fatalf("list err: %v len=%d", err, len(items))
	}
	if err := uc.Unassign(context.Background(), u.ID, items[0].ID); err != nil {
		t.Fatalf("unassign err: %v", err)
	}
	// tenant mismatch
	if _, err := uc.AssignByName(context.Background(), "other-"+tenant, u.ID, "teacher"); err == nil {
		t.Fatalf("expected tenant mismatch error")
	}
}
