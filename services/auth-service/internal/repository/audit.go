package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditRepo struct {
	db *pgxpool.Pool
}

func NewAuditRepo(db *pgxpool.Pool) *AuditRepo {
	return &AuditRepo{db: db}
}

func (r *AuditRepo) Log(ctx context.Context, tenantID string, userID *uuid.UUID, action, resourceType string, resourceID *uuid.UUID, newValues any) error {
	var uid *uuid.UUID
	if userID != nil {
		uid = userID
	}
	var rid *uuid.UUID
	if resourceID != nil {
		rid = resourceID
	}
	j, _ := json.Marshal(newValues)
	_, err := r.db.Exec(ctx, `
		INSERT INTO audit_logs (tenant_id, user_id, action, resource_type, resource_id, new_values)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb)
	`, tenantID, uid, action, resourceType, rid, string(j))
	return err
}

