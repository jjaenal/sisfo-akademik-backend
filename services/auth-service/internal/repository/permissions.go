package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Permission struct {
	ID        uuid.UUID
	Resource  string
	Action    string
}

type PermissionsRepo struct {
	db *pgxpool.Pool
}

func NewPermissionsRepo(db *pgxpool.Pool) *PermissionsRepo {
	return &PermissionsRepo{db: db}
}

func (r *PermissionsRepo) Create(ctx context.Context, resource, action string) (*Permission, error) {
	var out Permission
	err := r.db.QueryRow(ctx, `
		INSERT INTO permissions (resource, action)
		VALUES ($1, $2)
		RETURNING id, resource, action
	`, resource, action).Scan(&out.ID, &out.Resource, &out.Action)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PermissionsRepo) FindByResourceAction(ctx context.Context, resource, action string) (*Permission, error) {
	var out Permission
	err := r.db.QueryRow(ctx, `
		SELECT id, resource, action
		FROM permissions
		WHERE resource=$1 AND action=$2
		LIMIT 1
	`, resource, action).Scan(&out.ID, &out.Resource, &out.Action)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *PermissionsRepo) List(ctx context.Context, limit, offset int) ([]Permission, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, resource, action
		FROM permissions
		ORDER BY resource ASC, action ASC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *PermissionsRepo) AssignToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`, roleID, permissionID)
	return err
}

func (r *PermissionsRepo) GetByRole(ctx context.Context, roleID uuid.UUID) ([]Permission, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.resource, p.action
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Resource, &p.Action); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}
