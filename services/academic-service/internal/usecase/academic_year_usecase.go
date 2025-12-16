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

type academicYearUseCase struct {
	repo           repository.AcademicYearRepository
	contextTimeout time.Duration
}

var _ domainUseCase.AcademicYearUseCase = (*academicYearUseCase)(nil)

func NewAcademicYearUseCase(repo repository.AcademicYearRepository, timeout time.Duration) domainUseCase.AcademicYearUseCase {
	return &academicYearUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *academicYearUseCase) Create(ctx context.Context, academicYear *entity.AcademicYear) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := academicYear.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, academicYear)
}

func (u *academicYearUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.AcademicYear, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *academicYearUseCase) List(ctx context.Context, tenantID string) ([]entity.AcademicYear, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID)
}

func (u *academicYearUseCase) Update(ctx context.Context, academicYear *entity.AcademicYear) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := academicYear.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, academicYear)
}

func (u *academicYearUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
