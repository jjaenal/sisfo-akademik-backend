package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type TemplateRepository interface {
	Create(ctx context.Context, template *entity.ReportCardTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportCardTemplate, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.ReportCardTemplate, error)
	GetDefault(ctx context.Context, tenantID string) (*entity.ReportCardTemplate, error)
	Update(ctx context.Context, template *entity.ReportCardTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}
