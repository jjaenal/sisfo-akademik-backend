package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
)

type gradeCategoryUseCase struct {
	repo           repository.GradeCategoryRepository
	contextTimeout time.Duration
}

// Ensure interface implementation
var _ domainUseCase.GradeCategoryUseCase = (*gradeCategoryUseCase)(nil)

func NewGradeCategoryUseCase(repo repository.GradeCategoryRepository, timeout time.Duration) domainUseCase.GradeCategoryUseCase {
	return &gradeCategoryUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *gradeCategoryUseCase) Create(ctx context.Context, category *entity.GradeCategory) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := category.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, category)
}

func (u *gradeCategoryUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.GradeCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByID(ctx, id)
}

func (u *gradeCategoryUseCase) List(ctx context.Context) ([]*entity.GradeCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.List(ctx)
}

func (u *gradeCategoryUseCase) Update(ctx context.Context, category *entity.GradeCategory) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := category.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, category)
}

func (u *gradeCategoryUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.Delete(ctx, id)
}
