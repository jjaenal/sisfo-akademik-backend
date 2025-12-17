package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/repository"
)

type NotificationTemplateUseCase interface {
	Create(ctx context.Context, template *entity.NotificationTemplate) error
	Update(ctx context.Context, template *entity.NotificationTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error)
	List(ctx context.Context) ([]*entity.NotificationTemplate, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type notificationTemplateUseCase struct {
	repo    repository.NotificationTemplateRepository
	timeout time.Duration
}

func NewNotificationTemplateUseCase(repo repository.NotificationTemplateRepository, timeout time.Duration) NotificationTemplateUseCase {
	return &notificationTemplateUseCase{
		repo:    repo,
		timeout: timeout,
	}
}

func (u *notificationTemplateUseCase) Create(ctx context.Context, template *entity.NotificationTemplate) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if errs := template.Validate(); len(errs) > 0 {
		return errors.New("invalid input data")
	}

	template.ID = uuid.New()
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.IsActive = true

	return u.repo.Create(ctx, template)
}

func (u *notificationTemplateUseCase) Update(ctx context.Context, template *entity.NotificationTemplate) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if errs := template.Validate(); len(errs) > 0 {
		return errors.New("invalid input data")
	}

	existing, err := u.repo.GetByID(ctx, template.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("template not found")
	}

	template.UpdatedAt = time.Now()
	return u.repo.Update(ctx, template)
}

func (u *notificationTemplateUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *notificationTemplateUseCase) List(ctx context.Context) ([]*entity.NotificationTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.List(ctx)
}

func (u *notificationTemplateUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
