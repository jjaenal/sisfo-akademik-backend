package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type GradeCategoryRepository interface {
	Create(ctx context.Context, category *entity.GradeCategory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.GradeCategory, error)
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.GradeCategory, error)
	Update(ctx context.Context, category *entity.GradeCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}
