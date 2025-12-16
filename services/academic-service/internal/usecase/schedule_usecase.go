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

type scheduleUseCase struct {
	repo           repository.ScheduleRepository
	contextTimeout time.Duration
}

var _ domainUseCase.ScheduleUseCase = (*scheduleUseCase)(nil)

func NewScheduleUseCase(repo repository.ScheduleRepository, timeout time.Duration) domainUseCase.ScheduleUseCase {
	return &scheduleUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *scheduleUseCase) Create(ctx context.Context, schedule *entity.Schedule) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := schedule.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	// Check for conflicts
	conflicts, err := u.repo.CheckConflicts(ctx, schedule)
	if err != nil {
		return err
	}
	if len(conflicts) > 0 {
		return errors.New("schedule conflict detected: overlapping with existing schedule")
	}

	return u.repo.Create(ctx, schedule)
}

func (u *scheduleUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *scheduleUseCase) List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Schedule, int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID, limit, offset)
}

func (u *scheduleUseCase) ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.Schedule, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.ListByClass(ctx, classID)
}

func (u *scheduleUseCase) Update(ctx context.Context, schedule *entity.Schedule) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := schedule.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	// Check for conflicts
	conflicts, err := u.repo.CheckConflicts(ctx, schedule)
	if err != nil {
		return err
	}
	if len(conflicts) > 0 {
		return errors.New("schedule conflict detected: overlapping with existing schedule")
	}

	return u.repo.Update(ctx, schedule)
}

func (u *scheduleUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
