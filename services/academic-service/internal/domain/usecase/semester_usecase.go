package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type SemesterUseCase interface {
	Create(ctx context.Context, semester *entity.Semester) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Semester, error)
	List(ctx context.Context, tenantID string) ([]entity.Semester, error)
	ListByAcademicYear(ctx context.Context, academicYearID uuid.UUID) ([]entity.Semester, error)
	Update(ctx context.Context, semester *entity.Semester) error
	Delete(ctx context.Context, id uuid.UUID) error
	SetActive(ctx context.Context, id uuid.UUID) error
}
