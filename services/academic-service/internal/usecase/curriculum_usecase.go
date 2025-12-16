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

type curriculumUseCase struct {
	repo           repository.CurriculumRepository
	contextTimeout time.Duration
}

var _ domainUseCase.CurriculumUseCase = (*curriculumUseCase)(nil)

func NewCurriculumUseCase(repo repository.CurriculumRepository, timeout time.Duration) domainUseCase.CurriculumUseCase {
	return &curriculumUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *curriculumUseCase) Create(ctx context.Context, c *entity.Curriculum) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := c.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, c)
}

func (u *curriculumUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Curriculum, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *curriculumUseCase) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Curriculum, int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID, limit, offset)
}

func (u *curriculumUseCase) Update(ctx context.Context, c *entity.Curriculum) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := c.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, c)
}

func (u *curriculumUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}

func (u *curriculumUseCase) AddSubject(ctx context.Context, cs *entity.CurriculumSubject) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := cs.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.AddSubject(ctx, cs)
}

func (u *curriculumUseCase) RemoveSubject(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.RemoveSubject(ctx, id)
}

func (u *curriculumUseCase) ListSubjects(ctx context.Context, curriculumID uuid.UUID) ([]entity.CurriculumSubject, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.ListSubjects(ctx, curriculumID)
}
