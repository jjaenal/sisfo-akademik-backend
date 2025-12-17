package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/repository"
)

type admissionPeriodRepository struct {
	db *pgxpool.Pool
}

// NewAdmissionPeriodRepository creates a new instance of AdmissionPeriodRepository
func NewAdmissionPeriodRepository(db *pgxpool.Pool) repository.AdmissionPeriodRepository {
	return &admissionPeriodRepository{db: db}
}

// Create inserts a new admission period into the database
func (r *admissionPeriodRepository) Create(ctx context.Context, period *entity.AdmissionPeriod) error {
	query := `
		INSERT INTO admission_periods (
			id, name, start_date, end_date, is_active, is_announced, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`
	if period.ID == uuid.Nil {
		period.ID = uuid.New()
	}
	now := time.Now()
	if period.CreatedAt.IsZero() {
		period.CreatedAt = now
	}
	if period.UpdatedAt.IsZero() {
		period.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		period.ID, period.Name, period.StartDate, period.EndDate, period.IsActive, period.IsAnnounced,
		period.CreatedAt, period.UpdatedAt,
	)
	return err
}

// GetByID retrieves an admission period by its ID
func (r *admissionPeriodRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AdmissionPeriod, error) {
	query := `
		SELECT 
			id, name, start_date, end_date, is_active, is_announced, created_at, updated_at
		FROM admission_periods 
		WHERE id = $1
	`
	var period entity.AdmissionPeriod
	err := r.db.QueryRow(ctx, query, id).Scan(
		&period.ID, &period.Name, &period.StartDate, &period.EndDate, &period.IsActive, &period.IsAnnounced,
		&period.CreatedAt, &period.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &period, nil
}

// GetActive retrieves the currently active admission period
func (r *admissionPeriodRepository) GetActive(ctx context.Context) (*entity.AdmissionPeriod, error) {
	query := `
		SELECT 
			id, name, start_date, end_date, is_active, is_announced, created_at, updated_at
		FROM admission_periods 
		WHERE is_active = true
		LIMIT 1
	`
	var period entity.AdmissionPeriod
	err := r.db.QueryRow(ctx, query).Scan(
		&period.ID, &period.Name, &period.StartDate, &period.EndDate, &period.IsActive, &period.IsAnnounced,
		&period.CreatedAt, &period.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &period, nil
}

// List retrieves all admission periods
func (r *admissionPeriodRepository) List(ctx context.Context) ([]*entity.AdmissionPeriod, error) {
	query := `
		SELECT 
			id, name, start_date, end_date, is_active, is_announced, created_at, updated_at
		FROM admission_periods
		ORDER BY start_date DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []*entity.AdmissionPeriod
	for rows.Next() {
		var period entity.AdmissionPeriod
		err := rows.Scan(
			&period.ID, &period.Name, &period.StartDate, &period.EndDate, &period.IsActive, &period.IsAnnounced,
			&period.CreatedAt, &period.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		periods = append(periods, &period)
	}
	return periods, nil
}

// Update updates an existing admission period
func (r *admissionPeriodRepository) Update(ctx context.Context, period *entity.AdmissionPeriod) error {
	query := `
		UPDATE admission_periods 
		SET 
			name = $2, start_date = $3, end_date = $4, is_active = $5, is_announced = $6, updated_at = $7
		WHERE id = $1
	`
	period.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		period.ID, period.Name, period.StartDate, period.EndDate, period.IsActive, period.IsAnnounced,
		period.UpdatedAt,
	)
	return err
}

// Delete deletes an admission period by its ID
func (r *admissionPeriodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM admission_periods WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
