package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type TeacherUseCase interface {
	Create(ctx context.Context, teacher *entity.Teacher) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Teacher, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Teacher, int, error)
	Update(ctx context.Context, teacher *entity.Teacher) error
	Delete(ctx context.Context, id uuid.UUID) error
}
