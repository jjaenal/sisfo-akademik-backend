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

type teacherUseCase struct {
	repo           repository.TeacherRepository
	contextTimeout time.Duration
}

var _ domainUseCase.TeacherUseCase = (*teacherUseCase)(nil)

func NewTeacherUseCase(repo repository.TeacherRepository, timeout time.Duration) domainUseCase.TeacherUseCase {
	return &teacherUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *teacherUseCase) Create(ctx context.Context, teacher *entity.Teacher) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := teacher.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, teacher)
}

func (u *teacherUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Teacher, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *teacherUseCase) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Teacher, int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID, limit, offset)
}

func (u *teacherUseCase) Update(ctx context.Context, teacher *entity.Teacher) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := teacher.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, teacher)
}

func (u *teacherUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
