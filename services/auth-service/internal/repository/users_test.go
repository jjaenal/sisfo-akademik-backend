package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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
	// After soft delete, should not be found
	if _, err := repo.FindByID(context.Background(), u.ID); err == nil {
		t.Fatalf("expected find after soft delete to error")
	}
}

func TestUsersRepo_Create_Validation_MissingFields(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	if _, err := repo.Create(context.Background(), CreateUserParams{}); err == nil {
		t.Fatalf("expected validation error for missing fields")
	}
}

func TestUsersRepo_FindByEmail_NotFound(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	tenant := "test-" + uuid.NewString()
	if _, err := repo.FindByEmail(context.Background(), tenant, "notfound@example.com"); err == nil {
		t.Fatalf("expected not found error")
	}
}

func TestUsersRepo_Update_NoFields_ReturnsCurrent(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	tenant := "test-" + uuid.NewString()
	u, err := repo.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "user2@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	got, err := repo.Update(context.Background(), u.ID, UpdateUserParams{
		Now: time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if got.ID != u.ID || got.Email != u.Email {
		t.Fatalf("expected current user returned: %+v", got)
	}
}

func TestUsersRepo_Update_PasswordHashChanged(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	tenant := "test-" + uuid.NewString()
	u, err := repo.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "user3@example.com",
		Password: "OldPass123!",
	})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	oldHash := u.PasswordHash
	newPwd := "NewPass123!"
	upd, err := repo.Update(context.Background(), u.ID, UpdateUserParams{
		Password: &newPwd,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if upd.PasswordHash == oldHash {
		t.Fatalf("password hash should change")
	}
	if bcrypt.CompareHashAndPassword([]byte(upd.PasswordHash), []byte(newPwd)) != nil {
		t.Fatalf("new password does not match hash")
	}
	if bcrypt.CompareHashAndPassword([]byte(upd.PasswordHash), []byte("OldPass123!")) == nil {
		t.Fatalf("old password should not match hash")
	}
}

func TestUsersRepo_List_Pagination(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	tenant := "test-" + uuid.NewString()
	for i := 0; i < 12; i++ {
		email := "u" + uuid.NewString() + "@example.com"
		if _, err := repo.Create(context.Background(), CreateUserParams{
			TenantID: tenant,
			Email:    email,
			Password: "Password123!",
		}); err != nil {
			t.Fatalf("seed user %d err: %v", i, err)
		}
	}
	page1, total, err := repo.List(context.Background(), tenant, 5, 0)
	if err != nil || len(page1) != 5 || total < 12 {
		t.Fatalf("page1 err=%v len=%d total=%d", err, len(page1), total)
	}
	page2, _, err := repo.List(context.Background(), tenant, 5, 5)
	if err != nil || len(page2) != 5 {
		t.Fatalf("page2 err=%v len=%d", err, len(page2))
	}
	page3, _, err := repo.List(context.Background(), tenant, 5, 10)
	if err != nil || len(page3) < 2 {
		t.Fatalf("page3 err=%v len=%d", err, len(page3))
	}
}

func TestUsersRepo_SoftDelete_Excluded_From_List(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	tenant := "test-" + uuid.NewString()
	u, err := repo.Create(context.Background(), CreateUserParams{
		TenantID: tenant,
		Email:    "softdel1@example.com",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	if err := repo.SoftDelete(context.Background(), u.ID); err != nil {
		t.Fatalf("soft delete err: %v", err)
	}
	items, total, err := repo.List(context.Background(), tenant, 10, 0)
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	for _, it := range items {
		if it.ID == u.ID {
			t.Fatalf("soft-deleted user should be excluded from list")
		}
	}
	if total < len(items) {
		t.Fatalf("total should be >= items len")
	}
}

func TestUsersRepo_FindByID_NotFound(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewUsersRepo(db)
	if _, err := repo.FindByID(context.Background(), uuid.New()); err == nil {
		t.Fatalf("expected not found error")
	}
}
