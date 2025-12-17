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

type scheduleTemplateUseCase struct {
	repo           repository.ScheduleTemplateRepository
	contextTimeout time.Duration
}

var _ domainUseCase.ScheduleTemplateUseCase = (*scheduleTemplateUseCase)(nil)

func NewScheduleTemplateUseCase(repo repository.ScheduleTemplateRepository, timeout time.Duration) domainUseCase.ScheduleTemplateUseCase {
	return &scheduleTemplateUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *scheduleTemplateUseCase) Create(c context.Context, t *entity.ScheduleTemplate) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if err := t.Validate(); len(err) > 0 {
		return errors.New("validation error")
	}

	return u.repo.Create(ctx, t)
}

func (u *scheduleTemplateUseCase) GetByID(c context.Context, id uuid.UUID) (*entity.ScheduleTemplate, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	template, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, nil
	}

	items, err := u.repo.ListItems(ctx, id)
	if err != nil {
		return nil, err
	}
	template.Items = items

	return template, nil
}

func (u *scheduleTemplateUseCase) List(c context.Context, tenantID string) ([]entity.ScheduleTemplate, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.repo.List(ctx, tenantID)
}

func (u *scheduleTemplateUseCase) Update(c context.Context, t *entity.ScheduleTemplate) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.repo.Update(ctx, t)
}

func (u *scheduleTemplateUseCase) Delete(c context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}

func (u *scheduleTemplateUseCase) AddItem(c context.Context, item *entity.ScheduleTemplateItem) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if err := item.Validate(); len(err) > 0 {
		return errors.New("validation error")
	}

	return u.repo.AddItem(ctx, item)
}

func (u *scheduleTemplateUseCase) RemoveItem(c context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	return u.repo.RemoveItem(ctx, id)
}
