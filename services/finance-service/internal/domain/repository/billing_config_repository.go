package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
)

type BillingConfigRepository interface {
	Create(ctx context.Context, config *entity.BillingConfig) error
	Update(ctx context.Context, config *entity.BillingConfig) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.BillingConfig, error)
	List(ctx context.Context, tenantID uuid.UUID) ([]*entity.BillingConfig, error)
	ListAllActiveMonthly(ctx context.Context) ([]*entity.BillingConfig, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
