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

func TestRolesRepo_FindRoleByName_NotFound(t *testing.T) {
	db := testDB2(t)
	ensureMigrations(t, db)
	roles := NewRolesRepo(db)
	tenant := "test-" + uuid.NewString()
	if _, err := roles.FindRoleByName(context.Background(), tenant, "unknown"); err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestRolesRepo_UnassignNonExistentRole_NoError(t *testing.T) {
	db := testDB2(t)
	ensureMigrations(t, db)
	roles := NewRolesRepo(db)
	if err := roles.UnassignUserRole(context.Background(), uuid.New(), uuid.New()); err != nil {
		t.Fatalf("unassign non-existent should not error, got %v", err)
	}
}

func TestRolesRepo_ListUserRoles_Empty(t *testing.T) {
	db := testDB2(t)
	ensureMigrations(t, db)
	users := NewUsersRepo(db)
	roles := NewRolesRepo(db)
	tenant := "test-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "roles-empty@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	list, err := roles.ListUserRoles(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("expected empty roles, got %d", len(list))
	}
}

func TestRolesRepo_Assign_List_Unassign(t *testing.T) {
	db := testDB2(t)
	ensureMigrations(t, db)
	users := NewUsersRepo(db)
	roles := NewRolesRepo(db)
	tenant := "test-" + uuid.NewString()
	u, err := users.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "roles-assign@example.com",
		Password: "Password123!",
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
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	if len(list) != 1 || list[0].ID != r.ID {
		t.Fatalf("expected 1 role assigned")
	}
	if err := roles.UnassignUserRole(context.Background(), u.ID, r.ID); err != nil {
		t.Fatalf("unassign err: %v", err)
	}
	list2, err := roles.ListUserRoles(context.Background(), u.ID)
	if err != nil {
		t.Fatalf("list2 err: %v", err)
	}
	if len(list2) != 0 {
		t.Fatalf("expected empty after unassign")
	}
}
