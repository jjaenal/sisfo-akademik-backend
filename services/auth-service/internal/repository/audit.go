package repository

import (
	"context"
	"encoding/json"
	"strings"
	"time"

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

type AuditLog struct {
	ID           uuid.UUID       `json:"id"`
	TenantID     string          `json:"tenant_id"`
	UserID       *uuid.UUID      `json:"user_id"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resource_type"`
	ResourceID   *uuid.UUID      `json:"resource_id"`
	OldValues    json.RawMessage `json:"old_values"`
	NewValues    json.RawMessage `json:"new_values"`
	CreatedAt    time.Time       `json:"created_at"`
}

type ListParams struct {
	UserID       *uuid.UUID
	Action       string
	ResourceType string
	ResourceID   *uuid.UUID
	Start        *time.Time
	End          *time.Time
}

func (r *AuditRepo) List(ctx context.Context, tenantID string, p ListParams, limit, offset int) ([]AuditLog, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	conds := []string{"tenant_id=$1"}
	args := []any{tenantID}
	idx := 2
	if p.UserID != nil {
		conds = append(conds, "user_id=$"+itoa(idx))
		args = append(args, *p.UserID)
		idx++
	}
	if strings.TrimSpace(p.Action) != "" {
		conds = append(conds, "action=$"+itoa(idx))
		args = append(args, strings.TrimSpace(p.Action))
		idx++
	}
	if strings.TrimSpace(p.ResourceType) != "" {
		conds = append(conds, "resource_type=$"+itoa(idx))
		args = append(args, strings.TrimSpace(p.ResourceType))
		idx++
	}
	if p.ResourceID != nil {
		conds = append(conds, "resource_id=$"+itoa(idx))
		args = append(args, *p.ResourceID)
		idx++
	}
	if p.Start != nil {
		conds = append(conds, "created_at>=$"+itoa(idx))
		args = append(args, *p.Start)
		idx++
	}
	if p.End != nil {
		conds = append(conds, "created_at<=$"+itoa(idx))
		args = append(args, *p.End)
		idx++
	}
	where := "WHERE " + strings.Join(conds, " AND ")
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(1) FROM audit_logs "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	args2 := append(append([]any{}, args...), limit, offset)
	rows, err := r.db.Query(ctx, `
		SELECT id, tenant_id, user_id, action, resource_type, resource_id, old_values, new_values, created_at
		FROM audit_logs `+where+` ORDER BY created_at DESC LIMIT $`+itoa(idx)+` OFFSET $`+itoa(idx+1), args2...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []AuditLog
	for rows.Next() {
		var it AuditLog
		var oldb, newb []byte
		if err := rows.Scan(&it.ID, &it.TenantID, &it.UserID, &it.Action, &it.ResourceType, &it.ResourceID, &oldb, &newb, &it.CreatedAt); err != nil {
			return nil, 0, err
		}
		it.OldValues = json.RawMessage(oldb)
		it.NewValues = json.RawMessage(newb)
		out = append(out, it)
	}
	return out, total, nil
}

func (r *AuditRepo) Search(ctx context.Context, tenantID string, q string, limit, offset int) ([]AuditLog, int, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	pat := "%" + strings.ToLower(strings.TrimSpace(q)) + "%"
	var total int
	if err := r.db.QueryRow(ctx, `
		SELECT COUNT(1) FROM audit_logs
		WHERE tenant_id=$1 AND (
			LOWER(action) ILIKE $2 OR
			LOWER(resource_type) ILIKE $2 OR
			CAST(new_values AS TEXT) ILIKE $2 OR
			CAST(old_values AS TEXT) ILIKE $2
		)
	`, tenantID, pat).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, tenant_id, user_id, action, resource_type, resource_id, old_values, new_values, created_at
		FROM audit_logs
		WHERE tenant_id=$1 AND (
			LOWER(action) ILIKE $2 OR
			LOWER(resource_type) ILIKE $2 OR
			CAST(new_values AS TEXT) ILIKE $2 OR
			CAST(old_values AS TEXT) ILIKE $2
		)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, tenantID, pat, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var out []AuditLog
	for rows.Next() {
		var it AuditLog
		var oldb, newb []byte
		if err := rows.Scan(&it.ID, &it.TenantID, &it.UserID, &it.Action, &it.ResourceType, &it.ResourceID, &oldb, &newb, &it.CreatedAt); err != nil {
			return nil, 0, err
		}
		it.OldValues = json.RawMessage(oldb)
		it.NewValues = json.RawMessage(newb)
		out = append(out, it)
	}
	return out, total, nil
}

func (r *AuditRepo) CleanupOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	cmd, err := r.db.Exec(ctx, `DELETE FROM audit_logs WHERE created_at < $1`, cutoff)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}

func itoa(i int) string {
	var digits = "0123456789"
	if i == 0 {
		return "0"
	}
	out := make([]byte, 0, 3)
	for i > 0 {
		out = append([]byte{digits[i%10]}, out...)
		i /= 10
	}
	return string(out)
}
