package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
)

type curriculumRepository struct {
	db *pgxpool.Pool
}

var _ repository.CurriculumRepository = (*curriculumRepository)(nil)

func NewCurriculumRepository(db *pgxpool.Pool) repository.CurriculumRepository {
	return &curriculumRepository{db: db}
}

func (r *curriculumRepository) Create(ctx context.Context, c *entity.Curriculum) error {
	query := `
		INSERT INTO curricula (tenant_id, name, description, year, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		c.TenantID, c.Name, c.Description, c.Year, c.IsActive,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *curriculumRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Curriculum, error) {
	query := `
		SELECT id, tenant_id, name, description, year, is_active, created_at, updated_at
		FROM curricula
		WHERE id = $1 AND deleted_at IS NULL
	`
	var c entity.Curriculum
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.TenantID, &c.Name, &c.Description, &c.Year, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *curriculumRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Curriculum, int, error) {
	query := `
		SELECT id, tenant_id, name, description, year, is_active, created_at, updated_at
		FROM curricula
		WHERE tenant_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var curricula []entity.Curriculum
	for rows.Next() {
		var c entity.Curriculum
		if err := rows.Scan(
			&c.ID, &c.TenantID, &c.Name, &c.Description, &c.Year, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		curricula = append(curricula, c)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM curricula WHERE tenant_id = $1 AND deleted_at IS NULL`
	err = r.db.QueryRow(ctx, countQuery, tenantID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return curricula, total, nil
}

func (r *curriculumRepository) Update(ctx context.Context, c *entity.Curriculum) error {
	query := `
		UPDATE curricula
		SET name = $1, description = $2, year = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING updated_at
	`
	return r.db.QueryRow(ctx, query,
		c.Name, c.Description, c.Year, c.IsActive, c.ID,
	).Scan(&c.UpdatedAt)
}

func (r *curriculumRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE curricula SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *curriculumRepository) AddSubject(ctx context.Context, cs *entity.CurriculumSubject) error {
	query := `
		INSERT INTO curriculum_subjects (tenant_id, curriculum_id, subject_id, grade_level, semester, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		cs.TenantID, cs.CurriculumID, cs.SubjectID, cs.GradeLevel, cs.Semester,
	).Scan(&cs.ID, &cs.CreatedAt, &cs.UpdatedAt)
}

func (r *curriculumRepository) RemoveSubject(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE curriculum_subjects SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *curriculumRepository) ListSubjects(ctx context.Context, curriculumID uuid.UUID) ([]entity.CurriculumSubject, error) {
	query := `
		SELECT cs.id, cs.tenant_id, cs.curriculum_id, cs.subject_id, cs.grade_level, cs.semester, cs.created_at, cs.updated_at,
		       s.id, s.code, s.name, s.credit_units, s.type
		FROM curriculum_subjects cs
		JOIN subjects s ON cs.subject_id = s.id
		WHERE cs.curriculum_id = $1 AND cs.deleted_at IS NULL
		ORDER BY cs.grade_level, cs.semester
	`
	rows, err := r.db.Query(ctx, query, curriculumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []entity.CurriculumSubject
	for rows.Next() {
		var cs entity.CurriculumSubject
		var s entity.Subject
		// Scan both CS and Subject fields (Subject fields are mostly for display)
		// Note: We need to scan into Subject struct properly
		// s.ID, s.Code, s.Name, s.CreditUnits, s.Type
		if err := rows.Scan(
			&cs.ID, &cs.TenantID, &cs.CurriculumID, &cs.SubjectID, &cs.GradeLevel, &cs.Semester, &cs.CreatedAt, &cs.UpdatedAt,
			&s.ID, &s.Code, &s.Name, &s.CreditUnits, &s.Type,
		); err != nil {
			return nil, err
		}
		cs.Subject = &s
		subjects = append(subjects, cs)
	}

	return subjects, nil
}

func (r *curriculumRepository) AddGradingRule(ctx context.Context, rule *entity.GradingRule) error {
	query := `
		INSERT INTO grading_rules (
			id, tenant_id, curriculum_id, grade, min_score, max_score, points, description,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12
		)
	`
	if rule.ID == uuid.Nil {
		rule.ID = uuid.New()
	}
	now := time.Now()
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = now
	}
	if rule.UpdatedAt.IsZero() {
		rule.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		rule.ID, rule.TenantID, rule.CurriculumID, rule.Grade, rule.MinScore, rule.MaxScore, rule.Points, rule.Description,
		rule.CreatedAt, rule.UpdatedAt, rule.CreatedBy, rule.UpdatedBy,
	)
	return err
}

func (r *curriculumRepository) ListGradingRules(ctx context.Context, curriculumID uuid.UUID) ([]entity.GradingRule, error) {
	query := `
		SELECT 
			id, tenant_id, curriculum_id, grade, min_score, max_score, points, description,
			created_at, updated_at, created_by, updated_by
		FROM grading_rules
		WHERE curriculum_id = $1 AND deleted_at IS NULL
		ORDER BY min_score DESC
	`
	rows, err := r.db.Query(ctx, query, curriculumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []entity.GradingRule
	for rows.Next() {
		var rule entity.GradingRule
		err := rows.Scan(
			&rule.ID, &rule.TenantID, &rule.CurriculumID, &rule.Grade, &rule.MinScore, &rule.MaxScore, &rule.Points, &rule.Description,
			&rule.CreatedAt, &rule.UpdatedAt, &rule.CreatedBy, &rule.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *curriculumRepository) DeleteGradingRule(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE grading_rules SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
