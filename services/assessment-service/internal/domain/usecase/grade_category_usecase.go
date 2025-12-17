package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type GradeCategoryUseCase interface {
	Create(ctx context.Context, category *entity.GradeCategory) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.GradeCategory, error)
	List(ctx context.Context) ([]*entity.GradeCategory, error)
	Update(ctx context.Context, category *entity.GradeCategory) error
	Delete(ctx context.Context, id uuid.UUID) error
}
