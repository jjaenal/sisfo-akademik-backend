package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/repository"
)

type applicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) repository.ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, application *entity.Application) error {
	query := `
		INSERT INTO applications (
			id, tenant_id, admission_period_id, registration_number, first_name, last_name, 
			email, phone_number, status, previous_school, average_score, 
			submission_date, test_score, interview_score, final_score, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, 
			$7, $8, $9, $10, $11, 
			$12, $13, $14, $15, $16, $17
		)
	`
	if application.ID == uuid.Nil {
		application.ID = uuid.New()
	}
	if application.TenantID == "" {
		application.TenantID = "default"
	}
	now := time.Now()
	if application.CreatedAt.IsZero() {
		application.CreatedAt = now
	}
	if application.UpdatedAt.IsZero() {
		application.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		application.ID, application.TenantID, application.AdmissionPeriodID, application.RegistrationNumber, application.FirstName, application.LastName,
		application.Email, application.PhoneNumber, application.Status, application.PreviousSchool, application.AverageScore,
		application.SubmissionDate, application.TestScore, application.InterviewScore, application.FinalScore, application.CreatedAt, application.UpdatedAt,
	)
	return err
}

func (r *applicationRepository) Update(ctx context.Context, application *entity.Application) error {
	query := `
		UPDATE applications SET
			admission_period_id = $2, registration_number = $3, first_name = $4, last_name = $5,
			email = $6, phone_number = $7, status = $8, previous_school = $9, average_score = $10,
			submission_date = $11, test_score = $12, interview_score = $13, final_score = $14, updated_at = $15
		WHERE id = $1
	`
	application.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		application.ID, application.AdmissionPeriodID, application.RegistrationNumber, application.FirstName, application.LastName,
		application.Email, application.PhoneNumber, application.Status, application.PreviousSchool, application.AverageScore,
		application.SubmissionDate, application.TestScore, application.InterviewScore, application.FinalScore, application.UpdatedAt,
	)
	return err
}

func (r *applicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Application, error) {
	query := `
		SELECT 
			id, tenant_id, admission_period_id, registration_number, first_name, last_name, 
			email, phone_number, status, previous_school, average_score, 
			submission_date, test_score, interview_score, final_score, created_at, updated_at
		FROM applications 
		WHERE id = $1
	`
	var application entity.Application
	err := r.db.QueryRow(ctx, query, id).Scan(
		&application.ID, &application.TenantID, &application.AdmissionPeriodID, &application.RegistrationNumber, &application.FirstName, &application.LastName,
		&application.Email, &application.PhoneNumber, &application.Status, &application.PreviousSchool, &application.AverageScore,
		&application.SubmissionDate, &application.TestScore, &application.InterviewScore, &application.FinalScore, &application.CreatedAt, &application.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) GetByRegistrationNumber(ctx context.Context, regNum string) (*entity.Application, error) {
	query := `
		SELECT 
			id, admission_period_id, registration_number, first_name, last_name, 
			email, phone_number, status, previous_school, average_score, 
			submission_date, test_score, interview_score, final_score, created_at, updated_at
		FROM applications 
		WHERE registration_number = $1
	`
	var application entity.Application
	err := r.db.QueryRow(ctx, query, regNum).Scan(
		&application.ID, &application.AdmissionPeriodID, &application.RegistrationNumber, &application.FirstName, &application.LastName,
		&application.Email, &application.PhoneNumber, &application.Status, &application.PreviousSchool, &application.AverageScore,
		&application.SubmissionDate, &application.TestScore, &application.InterviewScore, &application.FinalScore, &application.CreatedAt, &application.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) List(ctx context.Context, filter map[string]interface{}) ([]*entity.Application, error) {
	query := `
		SELECT 
			id, admission_period_id, registration_number, first_name, last_name, 
			email, phone_number, status, previous_school, average_score, 
			submission_date, test_score, interview_score, final_score, created_at, updated_at
		FROM applications
	`
	
	var conditions []string
	var args []interface{}
	argCount := 1

	if periodID, ok := filter["admission_period_id"]; ok {
		conditions = append(conditions, fmt.Sprintf("admission_period_id = $%d", argCount))
		args = append(args, periodID)
		argCount++
	}

	if status, ok := filter["status"]; ok {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, status)
		argCount++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []*entity.Application
	for rows.Next() {
		var application entity.Application
		err := rows.Scan(
			&application.ID, &application.AdmissionPeriodID, &application.RegistrationNumber, &application.FirstName, &application.LastName,
			&application.Email, &application.PhoneNumber, &application.Status, &application.PreviousSchool, &application.AverageScore,
			&application.SubmissionDate, &application.TestScore, &application.InterviewScore, &application.FinalScore, &application.CreatedAt, &application.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		applications = append(applications, &application)
	}
	return applications, nil
}

func (r *applicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM applications WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
