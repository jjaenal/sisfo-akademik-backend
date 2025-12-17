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

type gradeRepository struct {
	db *pgxpool.Pool
}

func NewGradeRepository(db *pgxpool.Pool) repository.GradeRepository {
	return &gradeRepository{db: db}
}

func (r *gradeRepository) Create(ctx context.Context, grade *entity.Grade) error {
	query := `
		INSERT INTO grades (
			id, assessment_id, student_id, score, status, notes, graded_by,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9
		)
	`
	if grade.ID == uuid.Nil {
		grade.ID = uuid.New()
	}
	now := time.Now()
	if grade.CreatedAt.IsZero() {
		grade.CreatedAt = now
	}
	if grade.UpdatedAt.IsZero() {
		grade.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		grade.ID, grade.AssessmentID, grade.StudentID, grade.Score, grade.Status, grade.Notes, grade.GradedBy,
		grade.CreatedAt, grade.UpdatedAt,
	)
	return err
}

func (r *gradeRepository) Update(ctx context.Context, grade *entity.Grade) error {
	query := `
		UPDATE grades SET
			assessment_id = $2, student_id = $3, score = $4, status = $5, notes = $6, graded_by = $7,
			updated_at = $8
		WHERE id = $1
	`
	grade.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		grade.ID, grade.AssessmentID, grade.StudentID, grade.Score, grade.Status, grade.Notes, grade.GradedBy,
		grade.UpdatedAt,
	)
	return err
}

func (r *gradeRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Grade, error) {
	query := `
		SELECT 
			id, assessment_id, student_id, score, status, notes, graded_by,
			created_at, updated_at
		FROM grades 
		WHERE id = $1
	`
	var grade entity.Grade
	err := r.db.QueryRow(ctx, query, id).Scan(
		&grade.ID, &grade.AssessmentID, &grade.StudentID, &grade.Score, &grade.Status, &grade.Notes, &grade.GradedBy,
		&grade.CreatedAt, &grade.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

func (r *gradeRepository) GetByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (*entity.Grade, error) {
	query := `
		SELECT 
			id, assessment_id, student_id, score, status, notes, graded_by,
			created_at, updated_at
		FROM grades 
		WHERE student_id = $1 AND assessment_id = $2
	`
	var grade entity.Grade
	err := r.db.QueryRow(ctx, query, studentID, assessmentID).Scan(
		&grade.ID, &grade.AssessmentID, &grade.StudentID, &grade.Score, &grade.Status, &grade.Notes, &grade.GradedBy,
		&grade.CreatedAt, &grade.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

func (r *gradeRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.Grade, error) {
	query := `
		SELECT 
			id, assessment_id, student_id, score, status, notes, graded_by,
			created_at, updated_at
		FROM grades
	`
	
	var conditions []string
	var args []interface{}
	argCount := 1

	if assessmentID, ok := filter["assessment_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("assessment_id = $%d", argCount))
		args = append(args, assessmentID)
		argCount++
	}

	if studentID, ok := filter["student_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("student_id = $%d", argCount))
		args = append(args, studentID)
		argCount++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []*entity.Grade
	for rows.Next() {
		var grade entity.Grade
		err := rows.Scan(
			&grade.ID, &grade.AssessmentID, &grade.StudentID, &grade.Score, &grade.Status, &grade.Notes, &grade.GradedBy,
			&grade.CreatedAt, &grade.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		grades = append(grades, &grade)
	}
	return grades, nil
}

func (r *gradeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM grades WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
