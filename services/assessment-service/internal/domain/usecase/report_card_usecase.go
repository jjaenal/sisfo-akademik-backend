package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type ReportCardUseCase interface {
	Generate(ctx context.Context, tenantID string, studentID, classID, semesterID uuid.UUID) (*entity.ReportCard, error)
	GetByStudent(ctx context.Context, studentID, semesterID uuid.UUID) (*entity.ReportCard, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportCard, error)
	Publish(ctx context.Context, id uuid.UUID) error
	GetPDF(ctx context.Context, id uuid.UUID) ([]byte, error)
}
