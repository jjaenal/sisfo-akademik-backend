package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type AssessmentRepository interface {
	Create(ctx context.Context, assessment *entity.Assessment) error
	Update(ctx context.Context, assessment *entity.Assessment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Assessment, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*entity.Assessment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
