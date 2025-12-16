package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type EnrollmentUseCase interface {
	Enroll(ctx context.Context, enrollment *entity.Enrollment) error
	Unenroll(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Enrollment, error)
	ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.Enrollment, error)
	ListByStudent(ctx context.Context, studentID uuid.UUID) ([]entity.Enrollment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
