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

type teacherRepository struct {
	db *pgxpool.Pool
}

var _ repository.TeacherRepository = (*teacherRepository)(nil)

func NewTeacherRepository(db *pgxpool.Pool) repository.TeacherRepository {
	return &teacherRepository{db: db}
}

func (r *teacherRepository) Create(ctx context.Context, t *entity.Teacher) error {
	query := `
		INSERT INTO teachers (
			id, tenant_id, user_id, nip, name, gender, title_front, title_back,
			phone, email, status, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13, $14, $15
		)
	`
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		t.ID, t.TenantID, t.UserID, t.NIP, t.Name, t.Gender, t.TitleFront, t.TitleBack,
		t.Phone, t.Email, t.Status, t.CreatedAt, t.UpdatedAt, t.CreatedBy, t.UpdatedBy,
	)
	return err
}

func (r *teacherRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Teacher, error) {
	query := `
		SELECT 
			id, tenant_id, user_id, nip, name, gender, title_front, title_back,
			phone, email, status, created_at, updated_at, created_by, updated_by, deleted_at
		FROM teachers
		WHERE id = $1 AND deleted_at IS NULL
	`
	var t entity.Teacher
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.TenantID, &t.UserID, &t.NIP, &t.Name, &t.Gender, &t.TitleFront, &t.TitleBack,
		&t.Phone, &t.Email, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy, &t.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *teacherRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Teacher, int, error) {
	countQuery := `SELECT COUNT(*) FROM teachers WHERE tenant_id = $1 AND deleted_at IS NULL`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			id, tenant_id, user_id, nip, name, gender, title_front, title_back,
			phone, email, status, created_at, updated_at, created_by, updated_by, deleted_at
		FROM teachers
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var teachers []entity.Teacher
	for rows.Next() {
		var t entity.Teacher
		if err := rows.Scan(
			&t.ID, &t.TenantID, &t.UserID, &t.NIP, &t.Name, &t.Gender, &t.TitleFront, &t.TitleBack,
			&t.Phone, &t.Email, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy, &t.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		teachers = append(teachers, t)
	}
	return teachers, total, nil
}

func (r *teacherRepository) Update(ctx context.Context, t *entity.Teacher) error {
	query := `
		UPDATE teachers SET
			user_id = $1, nip = $2, name = $3, gender = $4,
			title_front = $5, title_back = $6, phone = $7, email = $8,
			status = $9, updated_at = $10, updated_by = $11
		WHERE id = $12 AND deleted_at IS NULL
	`
	t.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		t.UserID, t.NIP, t.Name, t.Gender,
		t.TitleFront, t.TitleBack, t.Phone, t.Email,
		t.Status, t.UpdatedAt, t.UpdatedBy, t.ID,
	)
	return err
}

func (r *teacherRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE teachers SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
