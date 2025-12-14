package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           uuid.UUID
	TenantID     string
	Email        string
	PasswordHash string
	IsActive     bool
}

type UsersRepo struct {
	db *pgxpool.Pool
}

func NewUsersRepo(db *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, password_hash, is_active
		FROM users
		WHERE id=$1 AND deleted_at IS NULL
		LIMIT 1
	`, id)
	var u User
	if err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.IsActive); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsersRepo) FindByEmail(ctx context.Context, tenantID, email string) (*User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, email, password_hash, is_active
		FROM users
		WHERE tenant_id=$1 AND email=$2 AND deleted_at IS NULL
		LIMIT 1
	`, tenantID, email)
	var u User
	if err := row.Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.IsActive); err != nil {
		return nil, err
	}
	return &u, nil
}
