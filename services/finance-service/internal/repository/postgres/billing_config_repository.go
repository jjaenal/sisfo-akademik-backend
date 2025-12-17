package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type billingConfigRepository struct {
	db *pgxpool.Pool
}

func NewBillingConfigRepository(db *pgxpool.Pool) repository.BillingConfigRepository {
	return &billingConfigRepository{db: db}
}

func (r *billingConfigRepository) Create(ctx context.Context, config *entity.BillingConfig) error {
	query := `
		INSERT INTO billing_configurations (
			id, tenant_id, name, amount, frequency, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`
	if config.ID == uuid.Nil {
		config.ID = uuid.New()
	}
	now := time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = now
	}
	if config.UpdatedAt.IsZero() {
		config.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		config.ID, config.TenantID, config.Name, config.Amount,
		config.Frequency, config.IsActive, config.CreatedAt, config.UpdatedAt,
	)
	return err
}

func (r *billingConfigRepository) Update(ctx context.Context, config *entity.BillingConfig) error {
	query := `
		UPDATE billing_configurations
		SET name = $1, amount = $2, frequency = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	config.UpdatedAt = time.Now()
	tag, err := r.db.Exec(ctx, query,
		config.Name, config.Amount, config.Frequency, config.IsActive, config.UpdatedAt, config.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("billing config not found")
	}
	return nil
}

func (r *billingConfigRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.BillingConfig, error) {
	query := `
		SELECT id, tenant_id, name, amount, frequency, is_active, created_at, updated_at
		FROM billing_configurations
		WHERE id = $1
	`
	var config entity.BillingConfig
	err := r.db.QueryRow(ctx, query, id).Scan(
		&config.ID, &config.TenantID, &config.Name, &config.Amount,
		&config.Frequency, &config.IsActive, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("billing config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *billingConfigRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*entity.BillingConfig, error) {
	query := `
		SELECT id, tenant_id, name, amount, frequency, is_active, created_at, updated_at
		FROM billing_configurations
		WHERE tenant_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*entity.BillingConfig
	for rows.Next() {
		var config entity.BillingConfig
		if err := rows.Scan(
			&config.ID, &config.TenantID, &config.Name, &config.Amount,
			&config.Frequency, &config.IsActive, &config.CreatedAt, &config.UpdatedAt,
		); err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

func (r *billingConfigRepository) ListAllActiveMonthly(ctx context.Context) ([]*entity.BillingConfig, error) {
	query := `
		SELECT id, tenant_id, name, amount, frequency, is_active, created_at, updated_at
		FROM billing_configurations
		WHERE is_active = true AND frequency = 'MONTHLY'
		ORDER BY tenant_id, created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*entity.BillingConfig
	for rows.Next() {
		var config entity.BillingConfig
		if err := rows.Scan(
			&config.ID, &config.TenantID, &config.Name, &config.Amount,
			&config.Frequency, &config.IsActive, &config.CreatedAt, &config.UpdatedAt,
		); err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

func (r *billingConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE billing_configurations SET is_active = false, updated_at = NOW() WHERE id = $1`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("billing config not found")
	}
	return nil
}
