package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
)

type GradeCategoryRepository struct {
	db DBPool
}

func NewGradeCategoryRepository(db DBPool) repository.GradeCategoryRepository {
	return &GradeCategoryRepository{db: db}
}

func (r *GradeCategoryRepository) Create(ctx context.Context, category *entity.GradeCategory) error {
	query := `
		INSERT INTO grade_categories (id, tenant_id, name, weight, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		category.ID,
		category.TenantID,
		category.Name,
		category.Weight,
		category.CreatedAt,
		category.UpdatedAt,
	)
	return err
}

func (r *GradeCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.GradeCategory, error) {
	query := `
		SELECT id, tenant_id, name, weight, created_at, updated_at, deleted_at
		FROM grade_categories
		WHERE id = $1 AND deleted_at IS NULL
	`
	var c entity.GradeCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.TenantID,
		&c.Name,
		&c.Weight,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GradeCategoryRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.GradeCategory, error) {
	query := `
		SELECT id, tenant_id, name, weight, created_at, updated_at, deleted_at
		FROM grade_categories
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.GradeCategory
	for rows.Next() {
		var c entity.GradeCategory
		if err := rows.Scan(
			&c.ID,
			&c.TenantID,
			&c.Name,
			&c.Weight,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *GradeCategoryRepository) Update(ctx context.Context, category *entity.GradeCategory) error {
	query := `
		UPDATE grade_categories
		SET name = $1, weight = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query,
		category.Name,
		category.Weight,
		category.UpdatedAt,
		category.ID,
	)
	return err
}

func (r *GradeCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE grade_categories
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
