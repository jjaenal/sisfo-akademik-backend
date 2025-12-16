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

type semesterRepository struct {
	db *pgxpool.Pool
}

var _ repository.SemesterRepository = (*semesterRepository)(nil)

func NewSemesterRepository(db *pgxpool.Pool) repository.SemesterRepository {
	return &semesterRepository{db: db}
}

func (r *semesterRepository) Create(ctx context.Context, s *entity.Semester) error {
	query := `
		INSERT INTO semesters (
			id, tenant_id, academic_year_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11
		)
	`
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	now := time.Now()
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		s.ID, s.TenantID, s.AcademicYearID, s.Name, s.StartDate, s.EndDate, s.IsActive,
		s.CreatedAt, s.UpdatedAt, s.CreatedBy, s.UpdatedBy,
	)
	return err
}

func (r *semesterRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Semester, error) {
	query := `
		SELECT 
			id, tenant_id, academic_year_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM semesters
		WHERE id = $1 AND deleted_at IS NULL
	`
	var s entity.Semester
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.TenantID, &s.AcademicYearID, &s.Name, &s.StartDate, &s.EndDate, &s.IsActive,
		&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *semesterRepository) List(ctx context.Context, tenantID string) ([]entity.Semester, error) {
	query := `
		SELECT 
			id, tenant_id, academic_year_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM semesters
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY start_date DESC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var semesters []entity.Semester
	for rows.Next() {
		var s entity.Semester
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.AcademicYearID, &s.Name, &s.StartDate, &s.EndDate, &s.IsActive,
			&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
		); err != nil {
			return nil, err
		}
		semesters = append(semesters, s)
	}
	return semesters, nil
}

func (r *semesterRepository) ListByAcademicYear(ctx context.Context, academicYearID uuid.UUID) ([]entity.Semester, error) {
	query := `
		SELECT 
			id, tenant_id, academic_year_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM semesters
		WHERE academic_year_id = $1 AND deleted_at IS NULL
		ORDER BY start_date ASC
	`
	rows, err := r.db.Query(ctx, query, academicYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var semesters []entity.Semester
	for rows.Next() {
		var s entity.Semester
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.AcademicYearID, &s.Name, &s.StartDate, &s.EndDate, &s.IsActive,
			&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
		); err != nil {
			return nil, err
		}
		semesters = append(semesters, s)
	}
	return semesters, nil
}

func (r *semesterRepository) GetActive(ctx context.Context, tenantID string) (*entity.Semester, error) {
	query := `
		SELECT 
			id, tenant_id, academic_year_id, name, start_date, end_date, is_active,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM semesters
		WHERE tenant_id = $1 AND is_active = true AND deleted_at IS NULL
		LIMIT 1
	`
	var s entity.Semester
	err := r.db.QueryRow(ctx, query, tenantID).Scan(
		&s.ID, &s.TenantID, &s.AcademicYearID, &s.Name, &s.StartDate, &s.EndDate, &s.IsActive,
		&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *semesterRepository) Update(ctx context.Context, s *entity.Semester) error {
	query := `
		UPDATE semesters SET
			academic_year_id = $1, name = $2, start_date = $3, end_date = $4, is_active = $5,
			updated_at = $6, updated_by = $7
		WHERE id = $8 AND deleted_at IS NULL
	`
	s.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		s.AcademicYearID, s.Name, s.StartDate, s.EndDate, s.IsActive,
		s.UpdatedAt, s.UpdatedBy, s.ID,
	)
	return err
}

func (r *semesterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE semesters SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
