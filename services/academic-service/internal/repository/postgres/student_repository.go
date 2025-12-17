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

type studentRepository struct {
	db *pgxpool.Pool
}

var _ repository.StudentRepository = (*studentRepository)(nil)

func NewStudentRepository(db *pgxpool.Pool) repository.StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, s *entity.Student) error {
	query := `
		INSERT INTO students (
			id, tenant_id, user_id, nis, nisn, name, gender, birth_place, birth_date,
			address, phone, email, parent_name, parent_phone, admission_date, status,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20
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
		s.ID, s.TenantID, s.UserID, s.NIS, s.NISN, s.Name, s.Gender, s.BirthPlace, s.BirthDate,
		s.Address, s.Phone, s.Email, s.ParentName, s.ParentPhone, s.AdmissionDate, s.Status,
		s.CreatedAt, s.UpdatedAt, s.CreatedBy, s.UpdatedBy,
	)
	return err
}

func (r *studentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Student, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, nis, nisn, name, gender, birth_place, birth_date,
			address, phone, email, parent_name, parent_phone, admission_date, status,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM students
		WHERE id = $1 AND deleted_at IS NULL
	`
	var s entity.Student
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.TenantID, &s.UserID, &s.NIS, &s.NISN, &s.Name, &s.Gender, &s.BirthPlace, &s.BirthDate,
		&s.Address, &s.Phone, &s.Email, &s.ParentName, &s.ParentPhone, &s.AdmissionDate, &s.Status,
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

func (r *studentRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Student, int, error) {
	// Count query
	countQuery := `SELECT COUNT(*) FROM students WHERE tenant_id = $1 AND deleted_at IS NULL`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	query := `
		SELECT 
			id, tenant_id, user_id, nis, nisn, name, gender, birth_place, birth_date,
			address, phone, email, parent_name, parent_phone, admission_date, status,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM students
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []entity.Student
	for rows.Next() {
		var s entity.Student
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.UserID, &s.NIS, &s.NISN, &s.Name, &s.Gender, &s.BirthPlace, &s.BirthDate,
			&s.Address, &s.Phone, &s.Email, &s.ParentName, &s.ParentPhone, &s.AdmissionDate, &s.Status,
			&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		students = append(students, s)
	}
	return students, total, nil
}

func (r *studentRepository) Update(ctx context.Context, s *entity.Student) error {
	query := `
		UPDATE students SET
			user_id = $1, nis = $2, nisn = $3, name = $4, gender = $5,
			birth_place = $6, birth_date = $7, address = $8, phone = $9,
			email = $10, parent_name = $11, parent_phone = $12,
			admission_date = $13, status = $14, updated_at = $15, updated_by = $16
		WHERE id = $17 AND deleted_at IS NULL
	`
	s.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		s.UserID, s.NIS, s.NISN, s.Name, s.Gender,
		s.BirthPlace, s.BirthDate, s.Address, s.Phone,
		s.Email, s.ParentName, s.ParentPhone,
		s.AdmissionDate, s.Status, s.UpdatedAt, s.UpdatedBy, s.ID,
	)
	return err
}

func (r *studentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE students SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
