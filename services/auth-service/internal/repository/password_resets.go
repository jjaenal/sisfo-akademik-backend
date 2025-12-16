package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PasswordReset struct {
	ID        uuid.UUID
	TenantID  string
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
}

type PasswordResetRepo struct {
	db *pgxpool.Pool
}

func NewPasswordResetRepo(db *pgxpool.Pool) *PasswordResetRepo {
	return &PasswordResetRepo{db: db}
}

func (r *PasswordResetRepo) EnsureSchema(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS password_resets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`)
	return err
}

func (r *PasswordResetRepo) Create(ctx context.Context, tenantID string, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*PasswordReset, error) {
	var pr PasswordReset
	err := r.db.QueryRow(ctx, `
		INSERT INTO password_resets (tenant_id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, tenant_id, user_id, token_hash, expires_at, used_at
	`, tenantID, userID, tokenHash, expiresAt).Scan(&pr.ID, &pr.TenantID, &pr.UserID, &pr.TokenHash, &pr.ExpiresAt, &pr.UsedAt)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *PasswordResetRepo) FindValidByTokenHash(ctx context.Context, tokenHash string, now time.Time) (*PasswordReset, error) {
	var pr PasswordReset
	err := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, user_id, token_hash, expires_at, used_at
		FROM password_resets
		WHERE token_hash=$1 AND used_at IS NULL AND expires_at > $2
		LIMIT 1
	`, tokenHash, now).Scan(&pr.ID, &pr.TenantID, &pr.UserID, &pr.TokenHash, &pr.ExpiresAt, &pr.UsedAt)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *PasswordResetRepo) MarkUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		UPDATE password_resets SET used_at=$1 WHERE id=$2 AND used_at IS NULL
	`, usedAt, id)
	return err
}
