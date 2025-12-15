package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDB2(t *testing.T) *pgxpool.Pool {
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

func TestRolesRepo_AssignListUnassign(t *testing.T) {
	db := testDB2(t)
	ensureMigrations(t, db)
	users := NewUsersRepo(db)
	roles := NewRolesRepo(db)
	tenant := "test-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "roles1@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	r, err := roles.CreateRole(context.Background(), tenant, "teacher", false)
	if err != nil {
		t.Fatalf("create role err: %v", err)
	}
	if err := roles.AssignUserRole(context.Background(), u.ID, r.ID); err != nil {
		t.Fatalf("assign err: %v", err)
	}
	list, err := roles.ListUserRoles(context.Background(), u.ID)
	if err != nil || len(list) == 0 {
		t.Fatalf("list err: %v len=%d", err, len(list))
	}
	if err := roles.UnassignUserRole(context.Background(), u.ID, r.ID); err != nil {
		t.Fatalf("unassign err: %v", err)
	}
}
