package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
)

type templateRepository struct {
	db *pgxpool.Pool
}

func NewTemplateRepository(db *pgxpool.Pool) repository.TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(ctx context.Context, template *entity.ReportCardTemplate) error {
	query := `
		INSERT INTO report_card_templates (id, tenant_id, name, config, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	configBytes, err := json.Marshal(template.Config)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query,
		template.ID,
		template.TenantID,
		template.Name,
		configBytes,
		template.IsDefault,
		template.CreatedAt,
		template.UpdatedAt,
	)
	return err
}

func (r *templateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportCardTemplate, error) {
	query := `
		SELECT id, tenant_id, name, config, is_default, created_at, updated_at
		FROM report_card_templates
		WHERE id = $1
	`

	var template entity.ReportCardTemplate
	var configBytes []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&template.ID,
		&template.TenantID,
		&template.Name,
		&configBytes,
		&template.IsDefault,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(configBytes, &template.Config); err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *templateRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.ReportCardTemplate, error) {
	query := `
		SELECT id, tenant_id, name, config, is_default, created_at, updated_at
		FROM report_card_templates
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*entity.ReportCardTemplate
	for rows.Next() {
		var template entity.ReportCardTemplate
		var configBytes []byte

		if err := rows.Scan(
			&template.ID,
			&template.TenantID,
			&template.Name,
			&configBytes,
			&template.IsDefault,
			&template.CreatedAt,
			&template.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(configBytes, &template.Config); err != nil {
			return nil, err
		}

		templates = append(templates, &template)
	}

	return templates, nil
}

func (r *templateRepository) GetDefault(ctx context.Context, tenantID string) (*entity.ReportCardTemplate, error) {
	query := `
		SELECT id, tenant_id, name, config, is_default, created_at, updated_at
		FROM report_card_templates
		WHERE tenant_id = $1 AND is_default = true
		LIMIT 1
	`

	var template entity.ReportCardTemplate
	var configBytes []byte

	err := r.db.QueryRow(ctx, query, tenantID).Scan(
		&template.ID,
		&template.TenantID,
		&template.Name,
		&configBytes,
		&template.IsDefault,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(configBytes, &template.Config); err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *templateRepository) Update(ctx context.Context, template *entity.ReportCardTemplate) error {
	query := `
		UPDATE report_card_templates
		SET name = $1, config = $2, is_default = $3, updated_at = $4
		WHERE id = $5
	`

	configBytes, err := json.Marshal(template.Config)
	if err != nil {
		return err
	}

	// If setting as default, unset others first (simple transaction could be better here)
	if template.IsDefault {
		unsetQuery := `UPDATE report_card_templates SET is_default = false WHERE tenant_id = $1 AND id != $2`
		if _, err := r.db.Exec(ctx, unsetQuery, template.TenantID, template.ID); err != nil {
			return err
		}
	}

	template.UpdatedAt = time.Now()
	_, err = r.db.Exec(ctx, query,
		template.Name,
		configBytes,
		template.IsDefault,
		template.UpdatedAt,
		template.ID,
	)
	return err
}

func (r *templateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM report_card_templates WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
