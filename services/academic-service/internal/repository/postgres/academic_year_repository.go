package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
)

type academicYearRepository struct {
	db *pgxpool.Pool
}

var _ repository.AcademicYearRepository = (*academicYearRepository)(nil)

func NewAcademicYearRepository(db *pgxpool.Pool) repository.AcademicYearRepository {
	return &academicYearRepository{db: db}
}

func (r *academicYearRepository) Create(ctx context.Context, a *entity.AcademicYear) error {
	query := `
		INSERT INTO academic_years (
			id, tenant_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10
		)
	`
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		a.ID, a.TenantID, a.Name, a.StartDate, a.EndDate, a.IsActive,
		a.CreatedAt, a.UpdatedAt, a.CreatedBy, a.UpdatedBy,
	)
	return err
}

func (r *academicYearRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AcademicYear, error) {
	query := `
		SELECT 
			id, tenant_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM academic_years
		WHERE id = $1 AND deleted_at IS NULL
	`
	var a entity.AcademicYear
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.TenantID, &a.Name, &a.StartDate, &a.EndDate, &a.IsActive,
		&a.CreatedAt, &a.UpdatedAt, &a.CreatedBy, &a.UpdatedBy, &a.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *academicYearRepository) List(ctx context.Context, tenantID string) ([]entity.AcademicYear, error) {
	query := `
		SELECT 
			id, tenant_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM academic_years
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY start_date DESC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []entity.AcademicYear
	for rows.Next() {
		var a entity.AcademicYear
		if err := rows.Scan(
			&a.ID, &a.TenantID, &a.Name, &a.StartDate, &a.EndDate, &a.IsActive,
			&a.CreatedAt, &a.UpdatedAt, &a.CreatedBy, &a.UpdatedBy, &a.DeletedAt,
		); err != nil {
			return nil, err
		}
		years = append(years, a)
	}
	return years, nil
}

func (r *academicYearRepository) Update(ctx context.Context, a *entity.AcademicYear) error {
	query := `
		UPDATE academic_years SET
			name = $1, start_date = $2, end_date = $3, is_active = $4,
			updated_at = $5, updated_by = $6
		WHERE id = $7 AND deleted_at IS NULL
	`
	a.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		a.Name, a.StartDate, a.EndDate, a.IsActive,
		a.UpdatedAt, a.UpdatedBy, a.ID,
	)
	return err
}

func (r *academicYearRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE academic_years SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
