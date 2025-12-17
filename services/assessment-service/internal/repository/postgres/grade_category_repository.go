package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
)

type gradeCategoryRepository struct {
	db *pgxpool.Pool
}

// NewGradeCategoryRepository creates a new instance of GradeCategoryRepository
func NewGradeCategoryRepository(db *pgxpool.Pool) repository.GradeCategoryRepository {
	return &gradeCategoryRepository{db: db}
}

// Create inserts a new grade category into the database
func (r *gradeCategoryRepository) Create(ctx context.Context, category *entity.GradeCategory) error {
	query := `
		INSERT INTO grade_categories (
			id, name, description, weight, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	if category.UpdatedAt.IsZero() {
		category.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		category.ID, category.Name, category.Description, category.Weight,
		category.CreatedAt, category.UpdatedAt,
	)
	return err
}

// GetByID retrieves a grade category by its ID
func (r *gradeCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.GradeCategory, error) {
	query := `
		SELECT 
			id, name, description, weight, created_at, updated_at
		FROM grade_categories 
		WHERE id = $1
	`
	var category entity.GradeCategory
	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID, &category.Name, &category.Description, &category.Weight,
		&category.CreatedAt, &category.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// List retrieves all grade categories
func (r *gradeCategoryRepository) List(ctx context.Context) ([]*entity.GradeCategory, error) {
	query := `
		SELECT 
			id, name, description, weight, created_at, updated_at
		FROM grade_categories
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.GradeCategory
	for rows.Next() {
		var category entity.GradeCategory
		err := rows.Scan(
			&category.ID, &category.Name, &category.Description, &category.Weight,
			&category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

// Update updates an existing grade category
func (r *gradeCategoryRepository) Update(ctx context.Context, category *entity.GradeCategory) error {
	query := `
		UPDATE grade_categories 
		SET 
			name = $2, description = $3, weight = $4, updated_at = $5
		WHERE id = $1
	`
	category.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		category.ID, category.Name, category.Description, category.Weight,
		category.UpdatedAt,
	)
	return err
}

// Delete deletes a grade category by its ID
func (r *gradeCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM grade_categories WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
