package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoAuthorizer struct {
	db *pgxpool.Pool
}

func NewRepoAuthorizer(db *pgxpool.Pool) *RepoAuthorizer {
	return &RepoAuthorizer{db: db}
}

func (a *RepoAuthorizer) Allow(subjectID uuid.UUID, tenantID string, permission string) (bool, error) {
	parts := strings.Split(permission, ":")
	if len(parts) != 2 {
		return false, nil
	}
	resource := strings.TrimSpace(parts[0])
	action := strings.TrimSpace(parts[1])
	if resource == "" || action == "" {
		return false, nil
	}
	var cnt int
	err := a.db.QueryRow(context.Background(), `
		SELECT COUNT(1)
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		JOIN role_permissions rp ON rp.role_id = r.id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id=$1 AND r.tenant_id=$2 AND r.deleted_at IS NULL AND p.resource=$3 AND p.action=$4
	`, subjectID, tenantID, resource, action).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}
