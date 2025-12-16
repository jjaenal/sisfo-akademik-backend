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

type subjectRepository struct {
	db *pgxpool.Pool
}

var _ repository.SubjectRepository = (*subjectRepository)(nil)

func NewSubjectRepository(db *pgxpool.Pool) repository.SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(ctx context.Context, s *entity.Subject) error {
	query := `
		INSERT INTO subjects (
			id, tenant_id, code, name, description, credit_units, type,
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
		s.ID, s.TenantID, s.Code, s.Name, s.Description, s.CreditUnits, s.Type,
		s.CreatedAt, s.UpdatedAt, s.CreatedBy, s.UpdatedBy,
	)
	return err
}

func (r *subjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error) {
	query := `
		SELECT 
			id, tenant_id, code, name, description, credit_units, type,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM subjects
		WHERE id = $1 AND deleted_at IS NULL
	`
	var s entity.Subject
	err := r.db.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.TenantID, &s.Code, &s.Name, &s.Description, &s.CreditUnits, &s.Type,
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

func (r *subjectRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Subject, int, error) {
	countQuery := `SELECT COUNT(*) FROM subjects WHERE tenant_id = $1 AND deleted_at IS NULL`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			id, tenant_id, code, name, description, credit_units, type,
			created_at, updated_at, created_by, updated_by, deleted_at
		FROM subjects
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subjects []entity.Subject
	for rows.Next() {
		var s entity.Subject
		if err := rows.Scan(
			&s.ID, &s.TenantID, &s.Code, &s.Name, &s.Description, &s.CreditUnits, &s.Type,
			&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		subjects = append(subjects, s)
	}
	return subjects, total, nil
}

func (r *subjectRepository) Update(ctx context.Context, s *entity.Subject) error {
	query := `
		UPDATE subjects SET
			code = $1, name = $2, description = $3, credit_units = $4,
			type = $5, updated_at = $6, updated_by = $7
		WHERE id = $8 AND deleted_at IS NULL
	`
	s.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		s.Code, s.Name, s.Description, s.CreditUnits,
		s.Type, s.UpdatedAt, s.UpdatedBy, s.ID,
	)
	return err
}

func (r *subjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE subjects SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
