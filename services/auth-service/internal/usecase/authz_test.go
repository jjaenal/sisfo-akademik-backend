package usecase

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
)

func testDBAuthz(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsAuthz(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b1))
	}
}

func TestRepoAuthorizer_Allow(t *testing.T) {
	db := testDBAuthz(t)
	ensureMigrationsAuthz(t, db)
	users := repository.NewUsersRepo(db)
	roles := repository.NewRolesRepo(db)
	perms := repository.NewPermissionsRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "authz1@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	role, err := roles.CreateRole(context.Background(), tenant, "teacher", false)
	if err != nil {
		t.Fatalf("create role err: %v", err)
	}
	res := "user_" + uuid.NewString()
	p, err := perms.Create(context.Background(), res, "read")
	if err != nil {
		t.Fatalf("create perm err: %v", err)
	}
	_ = roles.AssignUserRole(context.Background(), u.ID, role.ID)
	_ = perms.AssignToRole(context.Background(), role.ID, p.ID)
	authz := NewRepoAuthorizer(db)
	ok, err := authz.Allow(u.ID, tenant, res+":read")
	if err != nil || !ok {
		t.Fatalf("expected allow true, got ok=%v err=%v", ok, err)
	}
	ok, err = authz.Allow(u.ID, tenant, res+":write")
	if err != nil || ok {
		t.Fatalf("expected allow false for missing perm, got ok=%v err=%v", ok, err)
	}
	ok, err = authz.Allow(u.ID, "other-tenant", "user:read")
	if err != nil || ok {
		t.Fatalf("expected allow false for tenant mismatch, ok=%v err=%v", ok, err)
	}
	ok, err = authz.Allow(u.ID, tenant, "malformed")
	if err != nil || ok {
		t.Fatalf("expected allow false for malformed, ok=%v err=%v", ok, err)
	}
}

func TestRepoAuthorizer_Deny_On_DeletedRole(t *testing.T) {
	db := testDBAuthz(t)
	ensureMigrationsAuthz(t, db)
	users := repository.NewUsersRepo(db)
	roles := repository.NewRolesRepo(db)
	perms := repository.NewPermissionsRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "authz2@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	role, err := roles.CreateRole(context.Background(), tenant, "student", false)
	if err != nil {
		t.Fatalf("create role err: %v", err)
	}
	res := "grades_" + uuid.NewString()
	p, err := perms.Create(context.Background(), res, "read")
	if err != nil {
		t.Fatalf("create perm err: %v", err)
	}
	_ = roles.AssignUserRole(context.Background(), u.ID, role.ID)
	_ = perms.AssignToRole(context.Background(), role.ID, p.ID)
	authz := NewRepoAuthorizer(db)
	ok, err := authz.Allow(u.ID, tenant, res+":read")
	if err != nil || !ok {
		t.Fatalf("expected allow true before delete, ok=%v err=%v", ok, err)
	}
	_, _ = db.Exec(context.Background(), `UPDATE roles SET deleted_at=NOW() WHERE id=$1`, role.ID)
	ok, err = authz.Allow(u.ID, tenant, res+":read")
	if err != nil {
		t.Fatalf("allow err after delete: %v", err)
	}
	if ok {
		t.Fatalf("expected deny after role soft delete")
	}
}

func TestRepoAuthorizer_Allow_With_Multiple_Roles(t *testing.T) {
	db := testDBAuthz(t)
	ensureMigrationsAuthz(t, db)
	users := repository.NewUsersRepo(db)
	roles := repository.NewRolesRepo(db)
	perms := repository.NewPermissionsRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := users.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "authz3@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create user err: %v", err)
	}
	roleA, err := roles.CreateRole(context.Background(), tenant, "roleA", false)
	if err != nil {
		t.Fatalf("create roleA err: %v", err)
	}
	roleB, err := roles.CreateRole(context.Background(), tenant, "roleB", false)
	if err != nil {
		t.Fatalf("create roleB err: %v", err)
	}
	resA := "user_" + uuid.NewString()
	permA, err := perms.Create(context.Background(), resA, "read")
	if err != nil {
		t.Fatalf("create permA err: %v", err)
	}
	resB := "report_" + uuid.NewString()
	permB, err := perms.Create(context.Background(), resB, "export")
	if err != nil {
		t.Fatalf("create permB err: %v", err)
	}
	_ = roles.AssignUserRole(context.Background(), u.ID, roleA.ID)
	_ = roles.AssignUserRole(context.Background(), u.ID, roleB.ID)
	_ = perms.AssignToRole(context.Background(), roleA.ID, permA.ID)
	_ = perms.AssignToRole(context.Background(), roleB.ID, permB.ID)
	authz := NewRepoAuthorizer(db)
	ok, err := authz.Allow(u.ID, tenant, resA+":read")
	if err != nil || !ok {
		t.Fatalf("expected allow via roleA, ok=%v err=%v", ok, err)
	}
	ok, err = authz.Allow(u.ID, tenant, resB+":export")
	if err != nil || !ok {
		t.Fatalf("expected allow via roleB, ok=%v err=%v", ok, err)
	}
	ok, err = authz.Allow(u.ID, tenant, resA+":write")
	if err != nil || ok {
		t.Fatalf("expected deny for missing permission, ok=%v err=%v", ok, err)
	}
}
