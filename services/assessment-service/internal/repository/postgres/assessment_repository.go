package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
)

type assessmentRepository struct {
	db *pgxpool.Pool
}

func NewAssessmentRepository(db *pgxpool.Pool) repository.AssessmentRepository {
	return &assessmentRepository{db: db}
}

func (r *assessmentRepository) Create(ctx context.Context, assessment *entity.Assessment) error {
	query := `
		INSERT INTO assessments (
			id, teacher_id, subject_id, class_id, grade_category_id, 
			name, date, max_score, description, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, 
			$10, $11
		)
	`
	if assessment.ID == uuid.Nil {
		assessment.ID = uuid.New()
	}
	now := time.Now()
	if assessment.CreatedAt.IsZero() {
		assessment.CreatedAt = now
	}
	if assessment.UpdatedAt.IsZero() {
		assessment.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		assessment.ID, assessment.TeacherID, assessment.SubjectID, assessment.ClassID, assessment.GradeCategoryID,
		assessment.Name, assessment.Date, assessment.MaxScore, assessment.Description,
		assessment.CreatedAt, assessment.UpdatedAt,
	)
	return err
}

func (r *assessmentRepository) Update(ctx context.Context, assessment *entity.Assessment) error {
	query := `
		UPDATE assessments SET
			teacher_id = $2, subject_id = $3, class_id = $4, grade_category_id = $5,
			name = $6, date = $7, max_score = $8, description = $9,
			updated_at = $10
		WHERE id = $1
	`
	assessment.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		assessment.ID, assessment.TeacherID, assessment.SubjectID, assessment.ClassID, assessment.GradeCategoryID,
		assessment.Name, assessment.Date, assessment.MaxScore, assessment.Description,
		assessment.UpdatedAt,
	)
	return err
}

func (r *assessmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Assessment, error) {
	query := `
		SELECT 
			id, teacher_id, subject_id, class_id, grade_category_id, 
			name, date, max_score, description, 
			created_at, updated_at
		FROM assessments 
		WHERE id = $1
	`
	var assessment entity.Assessment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&assessment.ID, &assessment.TeacherID, &assessment.SubjectID, &assessment.ClassID, &assessment.GradeCategoryID,
		&assessment.Name, &assessment.Date, &assessment.MaxScore, &assessment.Description,
		&assessment.CreatedAt, &assessment.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &assessment, nil
}

func (r *assessmentRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.Assessment, error) {
	query := `
		SELECT 
			id, teacher_id, subject_id, class_id, grade_category_id, 
			name, date, max_score, description, 
			created_at, updated_at
		FROM assessments
	`
	
	var conditions []string
	var args []interface{}
	argCount := 1

	if classID, ok := filter["class_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("class_id = $%d", argCount))
		args = append(args, classID)
		argCount++
	}

	if subjectID, ok := filter["subject_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("subject_id = $%d", argCount))
		args = append(args, subjectID)
		argCount++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query += " ORDER BY date DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assessments []*entity.Assessment
	for rows.Next() {
		var assessment entity.Assessment
		err := rows.Scan(
			&assessment.ID, &assessment.TeacherID, &assessment.SubjectID, &assessment.ClassID, &assessment.GradeCategoryID,
			&assessment.Name, &assessment.Date, &assessment.MaxScore, &assessment.Description,
			&assessment.CreatedAt, &assessment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assessments = append(assessments, &assessment)
	}
	return assessments, nil
}

func (r *assessmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM assessments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
