package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Role struct {
	ID           uuid.UUID
	TenantID     string
	Name         string
	IsSystemRole bool
}

type RolesRepo struct {
	db *pgxpool.Pool
}

func NewRolesRepo(db *pgxpool.Pool) *RolesRepo {
	return &RolesRepo{db: db}
}

func (r *RolesRepo) CreateRole(ctx context.Context, tenantID, name string, system bool) (*Role, error) {
	var out Role
	err := r.db.QueryRow(ctx, `
		INSERT INTO roles (tenant_id, name, is_system_role)
		VALUES ($1, $2, $3)
		RETURNING id, tenant_id, name, is_system_role
	`, tenantID, name, system).Scan(&out.ID, &out.TenantID, &out.Name, &out.IsSystemRole)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *RolesRepo) FindRoleByName(ctx context.Context, tenantID, name string) (*Role, error) {
	var out Role
	err := r.db.QueryRow(ctx, `
		SELECT id, tenant_id, name, is_system_role
		FROM roles
		WHERE tenant_id=$1 AND name=$2 AND deleted_at IS NULL
		LIMIT 1
	`, tenantID, name).Scan(&out.ID, &out.TenantID, &out.Name, &out.IsSystemRole)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *RolesRepo) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	rows, err := r.db.Query(ctx, `
		SELECT r.id, r.tenant_id, r.name, r.is_system_role
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Role
	for rows.Next() {
		var rr Role
		if err := rows.Scan(&rr.ID, &rr.TenantID, &rr.Name, &rr.IsSystemRole); err != nil {
			return nil, err
		}
		out = append(out, rr)
	}
	return out, nil
}

func (r *RolesRepo) AssignUserRole(ctx context.Context, userID, roleID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`, userID, roleID)
	return err
}

func (r *RolesRepo) UnassignUserRole(ctx context.Context, userID, roleID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM user_roles WHERE user_id=$1 AND role_id=$2
	`, userID, roleID)
	return err
}
