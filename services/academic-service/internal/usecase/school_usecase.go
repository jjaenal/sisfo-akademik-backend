package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
)

type schoolUseCase struct {
	repo           repository.SchoolRepository
	contextTimeout time.Duration
}

// Ensure interface implementation
var _ domainUseCase.SchoolUseCase = (*schoolUseCase)(nil)

func NewSchoolUseCase(repo repository.SchoolRepository, timeout time.Duration) domainUseCase.SchoolUseCase {
	return &schoolUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *schoolUseCase) Create(ctx context.Context, school *entity.School) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate
	if errMap := school.Validate(); len(errMap) > 0 {
		// Convert map to error string or custom error type
		// For simplicity, just returning first error
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, school)
}

func (u *schoolUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.School, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByID(ctx, id)
}

func (u *schoolUseCase) GetByTenantID(ctx context.Context, tenantID string) (*entity.School, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByTenantID(ctx, tenantID)
}

func (u *schoolUseCase) Update(ctx context.Context, school *entity.School) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Validate
	if errMap := school.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, school)
}

func (u *schoolUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.Delete(ctx, id)
}
