package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDBPerm(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsPerm(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b1))
	}
}

func TestPermissionsRepo_Create_List_Assign_GetByRole(t *testing.T) {
	db := testDBPerm(t)
	ensureMigrationsPerm(t, db)
	pRepo := NewPermissionsRepo(db)
	rRepo := NewRolesRepo(db)
	uRepo := NewUsersRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := uRepo.Create(context.Background(), CreateUserParams{
		TenantID: tenant, Email: "perm1@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	role, err := rRepo.CreateRole(context.Background(), tenant, "teacher", false)
	if err != nil {
		t.Fatalf("create role err: %v", err)
	}
	res := "user_" + uuid.NewString()
	perm, err := pRepo.Create(context.Background(), res, "read")
	if err != nil {
		t.Fatalf("create perm err: %v", err)
	}
	items, err := pRepo.List(context.Background(), 10, 0)
	if err != nil || len(items) < 1 {
		t.Fatalf("list err: %v len=%d", err, len(items))
	}
	if err := pRepo.AssignToRole(context.Background(), role.ID, perm.ID); err != nil {
		t.Fatalf("assign to role err: %v", err)
	}
	if err := rRepo.AssignUserRole(context.Background(), u.ID, role.ID); err != nil {
		t.Fatalf("assign user role err: %v", err)
	}
	gotByRole, err := pRepo.GetByRole(context.Background(), role.ID)
	if err != nil || len(gotByRole) < 1 {
		t.Fatalf("get by role err=%v len=%d", err, len(gotByRole))
	}
	_ = pRepo.AssignToRole(context.Background(), role.ID, perm.ID)
	gotByRole2, err := pRepo.GetByRole(context.Background(), role.ID)
	if err != nil || len(gotByRole2) != 1 {
		t.Fatalf("duplicate assign should not increase count")
	}
}

func TestPermissionsRepo_GetByRole_Empty(t *testing.T) {
	db := testDBPerm(t)
	ensureMigrationsPerm(t, db)
	pRepo := NewPermissionsRepo(db)
	rRepo := NewRolesRepo(db)
	tenant := "t-" + uuid.NewString()
	role, err := rRepo.CreateRole(context.Background(), tenant, "empty-role", false)
	if err != nil {
		t.Fatalf("create role err: %v", err)
	}
	items, err := pRepo.GetByRole(context.Background(), role.ID)
	if err != nil {
		t.Fatalf("get by role err: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 permissions, got %d", len(items))
	}
}

func TestPermissionsRepo_FindByResourceAction_NotFound(t *testing.T) {
	db := testDBPerm(t)
	ensureMigrationsPerm(t, db)
	pRepo := NewPermissionsRepo(db)
	if _, err := pRepo.FindByResourceAction(context.Background(), "unknown_res", "read"); err == nil {
		t.Fatalf("expected not found error")
	}
}
