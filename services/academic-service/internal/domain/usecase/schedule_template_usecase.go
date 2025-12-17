package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type ScheduleTemplateUseCase interface {
	Create(ctx context.Context, template *entity.ScheduleTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ScheduleTemplate, error)
	List(ctx context.Context, tenantID string) ([]entity.ScheduleTemplate, error)
	Update(ctx context.Context, template *entity.ScheduleTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error

	AddItem(ctx context.Context, item *entity.ScheduleTemplateItem) error
	RemoveItem(ctx context.Context, id uuid.UUID) error
}
