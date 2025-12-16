package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type ClassRepository interface {
	Create(ctx context.Context, class *entity.Class) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Class, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Class, int, error)
	Update(ctx context.Context, class *entity.Class) error
	Delete(ctx context.Context, id uuid.UUID) error
}
