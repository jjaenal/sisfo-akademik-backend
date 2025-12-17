package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
)

type enrollmentRepository struct {
	db *pgxpool.Pool
}

var _ repository.EnrollmentRepository = (*enrollmentRepository)(nil)

func NewEnrollmentRepository(db *pgxpool.Pool) repository.EnrollmentRepository {
	return &enrollmentRepository{db: db}
}

func (r *enrollmentRepository) Enroll(ctx context.Context, e *entity.Enrollment) error {
	query := `
		INSERT INTO class_students (tenant_id, class_id, student_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	if e.Status == "" {
		e.Status = "active"
	}
	return r.db.QueryRow(ctx, query,
		e.TenantID, e.ClassID, e.StudentID, e.Status,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *enrollmentRepository) Unenroll(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE class_students SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *enrollmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Enrollment, error) {
	query := `
		SELECT id, tenant_id, class_id, student_id, status, created_at, updated_at
		FROM class_students
		WHERE id = $1 AND deleted_at IS NULL
	`
	var e entity.Enrollment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.TenantID, &e.ClassID, &e.StudentID, &e.Status, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *enrollmentRepository) ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.Enrollment, error) {
	query := `
		SELECT cs.id, cs.tenant_id, cs.class_id, cs.student_id, cs.status, cs.created_at, cs.updated_at,
		       s.id, s.name, s.nis, s.nisn, s.gender, s.email
		FROM class_students cs
		JOIN students s ON cs.student_id = s.id
		WHERE cs.class_id = $1 AND cs.deleted_at IS NULL
		ORDER BY s.name
	`
	rows, err := r.db.Query(ctx, query, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []entity.Enrollment
	for rows.Next() {
		var e entity.Enrollment
		var s entity.Student
		if err := rows.Scan(
			&e.ID, &e.TenantID, &e.ClassID, &e.StudentID, &e.Status, &e.CreatedAt, &e.UpdatedAt,
			&s.ID, &s.Name, &s.NIS, &s.NISN, &s.Gender, &s.Email,
		); err != nil {
			return nil, err
		}
		e.Student = &s
		enrollments = append(enrollments, e)
	}

	return enrollments, nil
}

func (r *enrollmentRepository) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]entity.Enrollment, error) {
	query := `
		SELECT cs.id, cs.tenant_id, cs.class_id, cs.student_id, cs.status, cs.created_at, cs.updated_at,
		       c.id, c.name, c.level, c.major
		FROM class_students cs
		JOIN classes c ON cs.class_id = c.id
		WHERE cs.student_id = $1 AND cs.deleted_at IS NULL
		ORDER BY cs.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []entity.Enrollment
	for rows.Next() {
		var e entity.Enrollment
		var c entity.Class
		if err := rows.Scan(
			&e.ID, &e.TenantID, &e.ClassID, &e.StudentID, &e.Status, &e.CreatedAt, &e.UpdatedAt,
			&c.ID, &c.Name, &c.Level, &c.Major,
		); err != nil {
			return nil, err
		}
		e.Class = &c
		enrollments = append(enrollments, e)
	}

	return enrollments, nil
}

func (r *enrollmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE class_students
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *enrollmentRepository) BulkEnroll(ctx context.Context, enrollments []*entity.Enrollment) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		INSERT INTO class_students (tenant_id, class_id, student_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	for _, e := range enrollments {
		if e.Status == "" {
			e.Status = "active"
		}
		err := tx.QueryRow(ctx, query,
			e.TenantID, e.ClassID, e.StudentID, e.Status,
		).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
