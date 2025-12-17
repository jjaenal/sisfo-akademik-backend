package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type BillingConfigUseCase interface {
	Create(ctx context.Context, config *entity.BillingConfig) error
	Update(ctx context.Context, config *entity.BillingConfig) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.BillingConfig, error)
	List(ctx context.Context, tenantID uuid.UUID) ([]*entity.BillingConfig, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type billingConfigUseCase struct {
	repo    repository.BillingConfigRepository
	timeout time.Duration
}

func NewBillingConfigUseCase(repo repository.BillingConfigRepository, timeout time.Duration) BillingConfigUseCase {
	return &billingConfigUseCase{
		repo:    repo,
		timeout: timeout,
	}
}

func (u *billingConfigUseCase) Create(ctx context.Context, config *entity.BillingConfig) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if errs := config.Validate(); len(errs) > 0 {
		return errors.New("invalid input data")
	}

	return u.repo.Create(ctx, config)
}

func (u *billingConfigUseCase) Update(ctx context.Context, config *entity.BillingConfig) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	if errs := config.Validate(); len(errs) > 0 {
		return errors.New("invalid input data")
	}

	return u.repo.Update(ctx, config)
}

func (u *billingConfigUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.BillingConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *billingConfigUseCase) List(ctx context.Context, tenantID uuid.UUID) ([]*entity.BillingConfig, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.List(ctx, tenantID)
}

func (u *billingConfigUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
