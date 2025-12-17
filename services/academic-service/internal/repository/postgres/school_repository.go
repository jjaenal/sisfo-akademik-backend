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

type schoolRepository struct {
	db *pgxpool.Pool
}

// Ensure interface implementation
var _ repository.SchoolRepository = (*schoolRepository)(nil)

// NewSchoolRepository creates a new instance of SchoolRepository
func NewSchoolRepository(db *pgxpool.Pool) repository.SchoolRepository {
	return &schoolRepository{db: db}
}

// Create inserts a new school into the database
func (r *schoolRepository) Create(ctx context.Context, school *entity.School) error {
	query := `
		INSERT INTO schools (
			id, tenant_id, name, address, phone, email, website, 
			logo_url, latitude, longitude, accreditation, headmaster, created_at, updated_at, 
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, 
			$8, $9, $10, $11, $12, $13, $14, 
			$15, $16
		)
	`
	if school.ID == uuid.Nil {
		school.ID = uuid.New()
	}
	now := time.Now()
	if school.CreatedAt.IsZero() {
		school.CreatedAt = now
	}
	if school.UpdatedAt.IsZero() {
		school.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		school.ID, school.TenantID, school.Name, school.Address, school.Phone, school.Email, school.Website,
		school.LogoURL, school.Latitude, school.Longitude, school.Accreditation, school.Headmaster, school.CreatedAt, school.UpdatedAt,
		school.CreatedBy, school.UpdatedBy,
	)
	return err
}

// GetByID retrieves a school by its ID
func (r *schoolRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.School, error) {
	query := `
		SELECT 
			id, tenant_id, name, address, phone, email, website, 
			logo_url, latitude, longitude, accreditation, headmaster, created_at, updated_at, 
			created_by, updated_by, deleted_at
		FROM schools 
		WHERE id = $1 AND deleted_at IS NULL
	`
	var school entity.School
	err := r.db.QueryRow(ctx, query, id).Scan(
		&school.ID, &school.TenantID, &school.Name, &school.Address, &school.Phone, &school.Email, &school.Website,
		&school.LogoURL, &school.Latitude, &school.Longitude, &school.Accreditation, &school.Headmaster, &school.CreatedAt, &school.UpdatedAt,
		&school.CreatedBy, &school.UpdatedBy, &school.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil if not found, let usecase handle error if needed
		}
		return nil, err
	}
	return &school, nil
}

// GetByTenantID retrieves a school by its Tenant ID
func (r *schoolRepository) GetByTenantID(ctx context.Context, tenantID string) (*entity.School, error) {
	query := `
		SELECT 
			id, tenant_id, name, address, phone, email, website, 
			logo_url, latitude, longitude, accreditation, headmaster, created_at, updated_at, 
			created_by, updated_by, deleted_at
		FROM schools 
		WHERE tenant_id = $1 AND deleted_at IS NULL
	`
	var school entity.School
	err := r.db.QueryRow(ctx, query, tenantID).Scan(
		&school.ID, &school.TenantID, &school.Name, &school.Address, &school.Phone, &school.Email, &school.Website,
		&school.LogoURL, &school.Latitude, &school.Longitude, &school.Accreditation, &school.Headmaster, &school.CreatedAt, &school.UpdatedAt,
		&school.CreatedBy, &school.UpdatedBy, &school.DeletedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &school, nil
}

// Update updates an existing school
func (r *schoolRepository) Update(ctx context.Context, school *entity.School) error {
	query := `
		UPDATE schools SET
			name = $1, address = $2, phone = $3, email = $4, website = $5,
			logo_url = $6, latitude = $7, longitude = $8, accreditation = $9, headmaster = $10, updated_at = $11, updated_by = $12
		WHERE id = $13 AND deleted_at IS NULL
	`
	school.UpdatedAt = time.Now()
	_, err := r.db.Exec(ctx, query,
		school.Name, school.Address, school.Phone, school.Email, school.Website,
		school.LogoURL, school.Latitude, school.Longitude, school.Accreditation, school.Headmaster, school.UpdatedAt, school.UpdatedBy,
		school.ID,
	)
	return err
}

// Delete soft-deletes a school
func (r *schoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE schools SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), id)
	return err
}
