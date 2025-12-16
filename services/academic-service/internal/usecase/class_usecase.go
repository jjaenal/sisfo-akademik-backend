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

type classUseCase struct {
	repo           repository.ClassRepository
	contextTimeout time.Duration
}

var _ domainUseCase.ClassUseCase = (*classUseCase)(nil)

func NewClassUseCase(repo repository.ClassRepository, timeout time.Duration) domainUseCase.ClassUseCase {
	return &classUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *classUseCase) Create(ctx context.Context, class *entity.Class) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := class.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, class)
}

func (u *classUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Class, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *classUseCase) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Class, int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID, limit, offset)
}

func (u *classUseCase) Update(ctx context.Context, class *entity.Class) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := class.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, class)
}

func (u *classUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
