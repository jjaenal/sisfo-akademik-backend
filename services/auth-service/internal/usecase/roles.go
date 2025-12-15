package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
)

type Roles interface {
	AssignByName(ctx context.Context, tenantID string, userID uuid.UUID, roleName string) (*repository.Role, error)
	List(ctx context.Context, userID uuid.UUID) ([]repository.Role, error)
	Unassign(ctx context.Context, userID, roleID uuid.UUID) error
}

type rolesUC struct {
	users *repository.UsersRepo
	roles *repository.RolesRepo
}

func NewRoles(users *repository.UsersRepo, roles *repository.RolesRepo) Roles {
	return &rolesUC{users: users, roles: roles}
}

func (r *rolesUC) AssignByName(ctx context.Context, tenantID string, userID uuid.UUID, roleName string) (*repository.Role, error) {
	u, err := r.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(tenantID) != u.TenantID {
		return nil, repository.ErrValidation("tenant mismatch")
	}
	roleName = strings.TrimSpace(roleName)
	if roleName == "" {
		return nil, repository.ErrValidation("role_name required")
	}
	role, err := r.roles.FindRoleByName(ctx, tenantID, roleName)
	if err != nil {
		role, err = r.roles.CreateRole(ctx, tenantID, roleName, false)
		if err != nil {
			return nil, err
		}
	}
	if err := r.roles.AssignUserRole(ctx, userID, role.ID); err != nil {
		return nil, err
	}
	return role, nil
}

func (r *rolesUC) List(ctx context.Context, userID uuid.UUID) ([]repository.Role, error) {
	return r.roles.ListUserRoles(ctx, userID)
}

func (r *rolesUC) Unassign(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.roles.UnassignUserRole(ctx, userID, roleID)
}
