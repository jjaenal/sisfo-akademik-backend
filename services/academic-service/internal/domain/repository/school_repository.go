package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type SchoolRepository interface {
	Create(ctx context.Context, school *entity.School) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.School, error)
	GetByTenantID(ctx context.Context, tenantID string) (*entity.School, error)
	Update(ctx context.Context, school *entity.School) error
	Delete(ctx context.Context, id uuid.UUID) error
}
