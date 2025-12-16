package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type StudentUseCase interface {
	Create(ctx context.Context, student *entity.Student) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Student, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Student, int, error)
	Update(ctx context.Context, student *entity.Student) error
	Delete(ctx context.Context, id uuid.UUID) error
}
