package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type ScheduleUseCase interface {
	Create(ctx context.Context, schedule *entity.Schedule) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Schedule, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Schedule, int, error)
	ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.Schedule, error)
	Update(ctx context.Context, schedule *entity.Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	BulkCreate(ctx context.Context, schedules []*entity.Schedule) error
	CreateFromTemplate(ctx context.Context, templateID uuid.UUID, classID uuid.UUID, teacherMap map[uuid.UUID]uuid.UUID) error
}
