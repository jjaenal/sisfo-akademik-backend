package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
)

type NotificationTemplateRepository interface {
	Create(ctx context.Context, template *entity.NotificationTemplate) error
	Update(ctx context.Context, template *entity.NotificationTemplate) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.NotificationTemplate, error)
	GetByName(ctx context.Context, name string) (*entity.NotificationTemplate, error)
	List(ctx context.Context) ([]*entity.NotificationTemplate, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
