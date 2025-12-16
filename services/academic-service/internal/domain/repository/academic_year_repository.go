package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type AcademicYearRepository interface {
	Create(ctx context.Context, academicYear *entity.AcademicYear) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AcademicYear, error)
	List(ctx context.Context, tenantID string) ([]entity.AcademicYear, error)
	Update(ctx context.Context, academicYear *entity.AcademicYear) error
	Delete(ctx context.Context, id uuid.UUID) error
}
