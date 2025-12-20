package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/domain"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgres(t *testing.T) (testcontainers.Container, *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		Env:          map[string]string{"POSTGRES_USER": "dev", "POSTGRES_PASSWORD": "dev", "POSTGRES_DB": "devdb"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	rc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Skipf("skip: cannot start postgres container: %v", err)
	}
	host, err := rc.Host(ctx)
	if err != nil {
		_ = rc.Terminate(ctx)
		t.Skipf("skip: cannot get host: %v", err)
	}
	port, err := rc.MappedPort(ctx, "5432/tcp")
	if err != nil {
		_ = rc.Terminate(ctx)
		t.Skipf("skip: cannot get mapped port: %v", err)
	}
	connStr := "postgres://dev:dev@" + host + ":" + port.Port() + "/devdb?sslmode=disable"
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		_ = rc.Terminate(ctx)
		t.Fatalf("parse cfg: %v", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		_ = rc.Terminate(ctx)
		t.Fatalf("pool new: %v", err)
	}
	// wait for connections
	for i := 0; i < 20; i++ {
		if _, err := pool.Exec(ctx, "SELECT 1"); err == nil {
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
	t.Cleanup(func() {
		pool.Close()
		_ = rc.Terminate(context.Background())
	})
	return rc, pool
}

func runMigrations(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	sql1 := `
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    path TEXT NOT NULL,
    bucket VARCHAR(50) NOT NULL,
    uploaded_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_files_tenant_id ON files(tenant_id);
CREATE INDEX IF NOT EXISTS idx_files_bucket ON files(bucket);
`
	sql2 := `
ALTER TABLE files ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX IF NOT EXISTS idx_files_deleted_at ON files(deleted_at);
`
	if _, err := pool.Exec(ctx, sql1); err != nil {
		t.Fatalf("migrate 1: %v", err)
	}
	if _, err := pool.Exec(ctx, sql2); err != nil {
		t.Fatalf("migrate 2: %v", err)
	}
}

func TestPostgresFileRepository_CRUD(t *testing.T) {
	_, pool := startPostgres(t)
	runMigrations(t, pool)
	repo := NewPostgresFileRepository(pool)
	ctx := context.Background()

	tenantID := uuid.New()
	uploader := uuid.New()
	id := uuid.New()
	file := &domain.File{
		ID:           id,
		TenantID:     tenantID,
		Name:         "stored.txt",
		OriginalName: "orig.txt",
		MimeType:     "text/plain",
		Size:         12,
		Path:         "uploads/stored.txt",
		Bucket:       "uploads",
		UploadedBy:   uploader,
	}
	if err := repo.Create(ctx, file); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got == nil || got.ID != id || got.OriginalName != "orig.txt" {
		t.Fatalf("unexpected get result")
	}
	list, err := repo.List(ctx, tenantID, 10, 0)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("list len=%d want 1", len(list))
	}
	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	after, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	if after != nil {
		t.Fatalf("should be nil after soft delete")
	}
}

func TestPostgresFileRepository_GetByID_NoRows(t *testing.T) {
	_, pool := startPostgres(t)
	runMigrations(t, pool)
	repo := NewPostgresFileRepository(pool)
	ctx := context.Background()
	id := uuid.New()
	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil for unknown id")
	}
}

func TestPostgresFileRepository_List_OrderAndPagination(t *testing.T) {
	_, pool := startPostgres(t)
	runMigrations(t, pool)
	repo := NewPostgresFileRepository(pool)
	ctx := context.Background()

	tenantID := uuid.New()
	uploader := uuid.New()

	file1 := &domain.File{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         "a.txt",
		OriginalName: "a.txt",
		MimeType:     "text/plain",
		Size:         10,
		Path:         "uploads/a.txt",
		Bucket:       "uploads",
		UploadedBy:   uploader,
	}
	if err := repo.Create(ctx, file1); err != nil {
		t.Fatalf("create f1: %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	file2 := &domain.File{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         "b.txt",
		OriginalName: "b.txt",
		MimeType:     "text/plain",
		Size:         20,
		Path:         "uploads/b.txt",
		Bucket:       "uploads",
		UploadedBy:   uploader,
	}
	if err := repo.Create(ctx, file2); err != nil {
		t.Fatalf("create f2: %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	file3 := &domain.File{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         "c.txt",
		OriginalName: "c.txt",
		MimeType:     "text/plain",
		Size:         30,
		Path:         "uploads/c.txt",
		Bucket:       "uploads",
		UploadedBy:   uploader,
	}
	if err := repo.Create(ctx, file3); err != nil {
		t.Fatalf("create f3: %v", err)
	}

	if err := repo.Delete(ctx, file2.ID); err != nil {
		t.Fatalf("delete f2: %v", err)
	}

	all, err := repo.List(ctx, tenantID, 10, 0)
	if err != nil {
		t.Fatalf("list all: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("list count=%d want 2", len(all))
	}
	if all[0].ID != file3.ID {
		t.Fatalf("order unexpected: first=%s want=%s", all[0].ID.String(), file3.ID.String())
	}

	page1, err := repo.List(ctx, tenantID, 1, 0)
	if err != nil {
		t.Fatalf("list page1: %v", err)
	}
	if len(page1) != 1 || page1[0].ID != file3.ID {
		t.Fatalf("page1 unexpected")
	}

	page2, err := repo.List(ctx, tenantID, 1, 1)
	if err != nil {
		t.Fatalf("list page2: %v", err)
	}
	if len(page2) != 1 || page2[0].ID != file1.ID {
		t.Fatalf("page2 unexpected")
	}
}

func TestPostgresFileRepository_DBErrors(t *testing.T) {
	_, pool := startPostgres(t)
	runMigrations(t, pool)
	repo := NewPostgresFileRepository(pool)
	ctx := context.Background()

	pool.Close()

	tenantID := uuid.New()
	uploader := uuid.New()
	id := uuid.New()
	file := &domain.File{
		ID:           id,
		TenantID:     tenantID,
		Name:         "err.txt",
		OriginalName: "err.txt",
		MimeType:     "text/plain",
		Size:         1,
		Path:         "uploads/err.txt",
		Bucket:       "uploads",
		UploadedBy:   uploader,
	}

	if err := repo.Create(ctx, file); err == nil {
		t.Fatalf("expected error on create with closed pool")
	}
	if _, err := repo.GetByID(ctx, id); err == nil {
		t.Fatalf("expected error on get with closed pool")
	}
	if err := repo.Delete(ctx, id); err == nil {
		t.Fatalf("expected error on delete with closed pool")
	}
	if _, err := repo.List(ctx, tenantID, 10, 0); err == nil {
		t.Fatalf("expected error on list with closed pool")
	}
}

func TestPostgresFileRepository_List_ScanError(t *testing.T) {
	rc, pool := startPostgres(t)
	runMigrations(t, pool)
	repo := NewPostgresFileRepository(pool)
	ctx := context.Background()

	tenantID := uuid.New()
	uploader := uuid.New()
	id := uuid.New()

	// Allow NULL in mime_type to force scan error into string field
	if _, err := pool.Exec(ctx, `ALTER TABLE files ALTER COLUMN mime_type DROP NOT NULL`); err != nil {
		_ = rc.Terminate(ctx)
		t.Fatalf("alter column: %v", err)
	}
	// Insert row with NULL mime_type
	ins := `
		INSERT INTO files (id, tenant_id, name, original_name, mime_type, size, path, bucket, uploaded_by, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, NULL, $5, $6, $7, $8, NOW(), NOW(), NULL)
	`
	if _, err := pool.Exec(ctx, ins, id, tenantID, "null-mime.txt", "null-mime.txt", int64(1), "uploads/null-mime.txt", "uploads", uploader); err != nil {
		_ = rc.Terminate(ctx)
		t.Fatalf("insert null mime: %v", err)
	}

	_, err := repo.List(ctx, tenantID, 10, 0)
	if err == nil {
		_ = rc.Terminate(ctx)
		t.Fatalf("expected scan error for NULL mime_type")
	}
}
