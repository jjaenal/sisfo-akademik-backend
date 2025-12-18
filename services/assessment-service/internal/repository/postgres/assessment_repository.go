package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
)

type AssessmentRepository struct {
	db DBPool
}

func NewAssessmentRepository(db DBPool) repository.AssessmentRepository {
	return &AssessmentRepository{db: db}
}

func (r *AssessmentRepository) Create(ctx context.Context, assessment *entity.Assessment) error {
	query := `
		INSERT INTO assessments (id, tenant_id, grade_category_id, teacher_id, subject_id, class_id, semester_id, name, description, max_score, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(ctx, query,
		assessment.ID,
		assessment.TenantID,
		assessment.GradeCategoryID,
		assessment.TeacherID,
		assessment.SubjectID,
		assessment.ClassID,
		assessment.SemesterID,
		assessment.Name,
		assessment.Description,
		assessment.MaxScore,
		assessment.Date,
		assessment.CreatedAt,
		assessment.UpdatedAt,
	)
	return err
}

func (r *AssessmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Assessment, error) {
	query := `
		SELECT id, tenant_id, grade_category_id, teacher_id, subject_id, class_id, semester_id, name, description, max_score, date, created_at, updated_at, deleted_at
		FROM assessments
		WHERE id = $1 AND deleted_at IS NULL
	`
	var a entity.Assessment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.TenantID,
		&a.GradeCategoryID,
		&a.TeacherID,
		&a.SubjectID,
		&a.ClassID,
		&a.SemesterID,
		&a.Name,
		&a.Description,
		&a.MaxScore,
		&a.Date,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AssessmentRepository) GetByClassAndSubject(ctx context.Context, classID, subjectID uuid.UUID) ([]*entity.Assessment, error) {
	query := `
		SELECT id, tenant_id, grade_category_id, teacher_id, subject_id, class_id, semester_id, name, description, max_score, date, created_at, updated_at, deleted_at
		FROM assessments
		WHERE class_id = $1 AND subject_id = $2 AND deleted_at IS NULL
		ORDER BY date DESC
	`
	rows, err := r.db.Query(ctx, query, classID, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []*entity.Assessment
	for rows.Next() {
		var a entity.Assessment
		if err := rows.Scan(
			&a.ID,
			&a.TenantID,
			&a.GradeCategoryID,
			&a.TeacherID,
			&a.SubjectID,
			&a.ClassID,
			&a.SemesterID,
			&a.Name,
			&a.Description,
			&a.MaxScore,
			&a.Date,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.DeletedAt,
		); err != nil {
			return nil, err
		}
		assessments = append(assessments, &a)
	}
	return assessments, nil
}

func (r *AssessmentRepository) Update(ctx context.Context, assessment *entity.Assessment) error {
	query := `
		UPDATE assessments
		SET name = $1, description = $2, max_score = $3, date = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query,
		assessment.Name,
		assessment.Description,
		assessment.MaxScore,
		assessment.Date,
		assessment.UpdatedAt,
		assessment.ID,
	)
	return err
}

func (r *AssessmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE assessments
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
