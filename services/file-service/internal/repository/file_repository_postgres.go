package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/domain"
)

type PostgresFileRepository struct {
	db *pgxpool.Pool
}

func NewPostgresFileRepository(db *pgxpool.Pool) *PostgresFileRepository {
	return &PostgresFileRepository{db: db}
}

func (r *PostgresFileRepository) Create(ctx context.Context, file *domain.File) error {
	query := `
		INSERT INTO files (id, tenant_id, name, original_name, mime_type, size, path, bucket, uploaded_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		file.ID,
		file.TenantID,
		file.Name,
		file.OriginalName,
		file.MimeType,
		file.Size,
		file.Path,
		file.Bucket,
		file.UploadedBy,
		file.CreatedAt,
		file.UpdatedAt,
	)
	return err
}

func (r *PostgresFileRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	query := `
		SELECT id, tenant_id, name, original_name, mime_type, size, path, bucket, uploaded_by, created_at, updated_at, deleted_at
		FROM files
		WHERE id = $1 AND deleted_at IS NULL
	`
	var file domain.File
	err := r.db.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.TenantID,
		&file.Name,
		&file.OriginalName,
		&file.MimeType,
		&file.Size,
		&file.Path,
		&file.Bucket,
		&file.UploadedBy,
		&file.CreatedAt,
		&file.UpdatedAt,
		&file.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *PostgresFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE files SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *PostgresFileRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.File, error) {
	query := `
		SELECT id, tenant_id, name, original_name, mime_type, size, path, bucket, uploaded_by, created_at, updated_at, deleted_at
		FROM files
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*domain.File
	for rows.Next() {
		var file domain.File
		if err := rows.Scan(
			&file.ID,
			&file.TenantID,
			&file.Name,
			&file.OriginalName,
			&file.MimeType,
			&file.Size,
			&file.Path,
			&file.Bucket,
			&file.UploadedBy,
			&file.CreatedAt,
			&file.UpdatedAt,
			&file.DeletedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, &file)
	}
	return files, nil
}
