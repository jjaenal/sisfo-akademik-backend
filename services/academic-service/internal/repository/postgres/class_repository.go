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

type classRepository struct {
	db *pgxpool.Pool
}

var _ repository.ClassRepository = (*classRepository)(nil)

func NewClassRepository(db *pgxpool.Pool) repository.ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(ctx context.Context, c *entity.Class) error {
	query := `
		INSERT INTO classes (
			id, tenant_id, school_id, academic_year_id, name, level, major,
			homeroom_teacher_id, capacity, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13
		)
	`
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	if c.UpdatedAt.IsZero() {
		c.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		c.ID, c.TenantID, c.SchoolID, c.AcademicYearID, c.Name, c.Level, c.Major,
		c.HomeroomTeacherID, c.Capacity, c.CreatedAt, c.UpdatedAt, c.CreatedBy, c.UpdatedBy,
	)
	return err
}

func (r *classRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Class, error) {
	query := `
		SELECT 
			id, tenant_id, school_id, academic_year_id, name, level, major,
			homeroom_teacher_id, capacity, created_at, updated_at, created_by, updated_by, deleted_at
		FROM classes
		WHERE id = $1 AND deleted_at IS NULL
	`
	var c entity.Class
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.TenantID, &c.SchoolID, &c.AcademicYearID, &c.Name, &c.Level, &c.Major,
		&c.HomeroomTeacherID, &c.Capacity, &c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy, &c.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *classRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Class, int, error) {
	countQuery := `SELECT COUNT(*) FROM classes WHERE tenant_id = $1 AND deleted_at IS NULL`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			id, tenant_id, school_id, academic_year_id, name, level, major,
			homeroom_teacher_id, capacity, created_at, updated_at, created_by, updated_by, deleted_at
		FROM classes
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var classes []entity.Class
	for rows.Next() {
		var c entity.Class
		if err := rows.Scan(
			&c.ID, &c.TenantID, &c.SchoolID, &c.AcademicYearID, &c.Name, &c.Level, &c.Major,
			&c.HomeroomTeacherID, &c.Capacity, &c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy, &c.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		classes = append(classes, c)
	}
	return classes, total, nil
}

func (r *classRepository) Update(ctx context.Context, c *entity.Class) error {
	query := `
		UPDATE classes SET
			school_id = $1, academic_year_id = $2, name = $3, level = $4,
			major = $5, homeroom_teacher_id = $6, capacity = $7,
			updated_at = $8, updated_by = $9
		WHERE id = $10 AND deleted_at IS NULL
	`
	c.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		c.SchoolID, c.AcademicYearID, c.Name, c.Level,
		c.Major, c.HomeroomTeacherID, c.Capacity,
		c.UpdatedAt, c.UpdatedBy, c.ID,
	)
	return err
}

func (r *classRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE classes SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
