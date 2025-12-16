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

type studentUseCase struct {
	repo           repository.StudentRepository
	contextTimeout time.Duration
}

var _ domainUseCase.StudentUseCase = (*studentUseCase)(nil)

func NewStudentUseCase(repo repository.StudentRepository, timeout time.Duration) domainUseCase.StudentUseCase {
	return &studentUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *studentUseCase) Create(ctx context.Context, student *entity.Student) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := student.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, student)
}

func (u *studentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Student, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *studentUseCase) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Student, int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID, limit, offset)
}

func (u *studentUseCase) Update(ctx context.Context, student *entity.Student) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := student.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, student)
}

func (u *studentUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
