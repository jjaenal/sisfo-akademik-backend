package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type ReportCardRepository interface {
	Create(ctx context.Context, rc *entity.ReportCard) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportCard, error)
	GetByStudentAndSemester(ctx context.Context, studentID, semesterID uuid.UUID) (*entity.ReportCard, error)
	Update(ctx context.Context, rc *entity.ReportCard) error
	List(ctx context.Context, classID, semesterID uuid.UUID) ([]*entity.ReportCard, error)
}
