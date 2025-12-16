package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *entity.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Schedule, int, error)
	ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.Schedule, error)
	Update(ctx context.Context, schedule *entity.Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	CheckConflicts(ctx context.Context, schedule *entity.Schedule) ([]entity.Schedule, error)
}
