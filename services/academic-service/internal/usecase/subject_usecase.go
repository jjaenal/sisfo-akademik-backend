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

type subjectUseCase struct {
	repo           repository.SubjectRepository
	contextTimeout time.Duration
}

var _ domainUseCase.SubjectUseCase = (*subjectUseCase)(nil)

func NewSubjectUseCase(repo repository.SubjectRepository, timeout time.Duration) domainUseCase.SubjectUseCase {
	return &subjectUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *subjectUseCase) Create(ctx context.Context, subject *entity.Subject) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := subject.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, subject)
}

func (u *subjectUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *subjectUseCase) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Subject, int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID, limit, offset)
}

func (u *subjectUseCase) Update(ctx context.Context, subject *entity.Subject) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := subject.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, subject)
}

func (u *subjectUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
