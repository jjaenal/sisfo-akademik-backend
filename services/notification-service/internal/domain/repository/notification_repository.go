package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	Update(ctx context.Context, notification *entity.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	ListByRecipient(ctx context.Context, recipient string) ([]*entity.Notification, error)
}
