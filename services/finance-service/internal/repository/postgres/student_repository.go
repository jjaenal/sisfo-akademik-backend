package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
)

type StudentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(ctx context.Context, student *entity.Student) error {
	query := `
		INSERT INTO students (id, tenant_id, name, status, class_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query,
		student.ID,
		student.TenantID,
		student.Name,
		student.Status,
		student.ClassID,
		student.CreatedAt,
		student.UpdatedAt,
	)
	return err
}

func (r *StudentRepository) GetActive(ctx context.Context, tenantID uuid.UUID) ([]*entity.Student, error) {
	query := `
		SELECT id, tenant_id, name, status, class_id, created_at, updated_at
		FROM students
		WHERE tenant_id = $1 AND status = 'active'
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []*entity.Student
	for rows.Next() {
		var s entity.Student
		if err := rows.Scan(
			&s.ID,
			&s.TenantID,
			&s.Name,
			&s.Status,
			&s.ClassID,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		students = append(students, &s)
	}
	return students, nil
}
