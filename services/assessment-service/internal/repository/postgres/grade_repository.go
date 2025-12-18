package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
)

type GradeRepository struct {
	db DBPool
}

func NewGradeRepository(db DBPool) repository.GradeRepository {
	return &GradeRepository{db: db}
}

func (r *GradeRepository) Create(ctx context.Context, grade *entity.Grade) error {
	query := `
		INSERT INTO grades (id, tenant_id, assessment_id, student_id, score, feedback, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query,
		grade.ID,
		grade.TenantID,
		grade.AssessmentID,
		grade.StudentID,
		grade.Score,
		grade.Feedback,
		grade.CreatedAt,
		grade.UpdatedAt,
	)
	return err
}

func (r *GradeRepository) CreateBulk(ctx context.Context, grades []*entity.Grade) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	batch := &pgx.Batch{}
	query := `
		INSERT INTO grades (id, tenant_id, assessment_id, student_id, score, feedback, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (assessment_id, student_id) 
		DO UPDATE SET score = EXCLUDED.score, feedback = EXCLUDED.feedback, updated_at = EXCLUDED.updated_at
	`

	for _, grade := range grades {
		batch.Queue(query,
			grade.ID,
			grade.TenantID,
			grade.AssessmentID,
			grade.StudentID,
			grade.Score,
			grade.Feedback,
			grade.CreatedAt,
			grade.UpdatedAt,
		)
	}

	br := tx.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *GradeRepository) GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) ([]*entity.Grade, error) {
	query := `
		SELECT id, tenant_id, assessment_id, student_id, score, feedback, created_at, updated_at, deleted_at
		FROM grades
		WHERE assessment_id = $1 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, assessmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []*entity.Grade
	for rows.Next() {
		var g entity.Grade
		if err := rows.Scan(
			&g.ID,
			&g.TenantID,
			&g.AssessmentID,
			&g.StudentID,
			&g.Score,
			&g.Feedback,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.DeletedAt,
		); err != nil {
			return nil, err
		}
		grades = append(grades, &g)
	}
	return grades, nil
}

func (r *GradeRepository) GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]*entity.Grade, error) {
	query := `
		SELECT id, tenant_id, assessment_id, student_id, score, feedback, created_at, updated_at, deleted_at
		FROM grades
		WHERE student_id = $1 AND deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []*entity.Grade
	for rows.Next() {
		var g entity.Grade
		if err := rows.Scan(
			&g.ID,
			&g.TenantID,
			&g.AssessmentID,
			&g.StudentID,
			&g.Score,
			&g.Feedback,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.DeletedAt,
		); err != nil {
			return nil, err
		}
		grades = append(grades, &g)
	}
	return grades, nil
}

func (r *GradeRepository) GetByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (*entity.Grade, error) {
	query := `
		SELECT id, tenant_id, assessment_id, student_id, score, feedback, created_at, updated_at, deleted_at
		FROM grades
		WHERE student_id = $1 AND assessment_id = $2 AND deleted_at IS NULL
	`
	var g entity.Grade
	err := r.db.QueryRow(ctx, query, studentID, assessmentID).Scan(
		&g.ID,
		&g.TenantID,
		&g.AssessmentID,
		&g.StudentID,
		&g.Score,
		&g.Feedback,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *GradeRepository) Update(ctx context.Context, grade *entity.Grade) error {
	query := `
		UPDATE grades
		SET score = $1, feedback = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query,
		grade.Score,
		grade.Feedback,
		grade.UpdatedAt,
		grade.ID,
	)
	return err
}
