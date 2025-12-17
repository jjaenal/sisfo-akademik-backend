package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
)

type ApplicationRepository interface {
	Create(ctx context.Context, application *entity.Application) error
	Update(ctx context.Context, application *entity.Application) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Application, error)
	GetByRegistrationNumber(ctx context.Context, regNum string) (*entity.Application, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*entity.Application, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
