package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type GradingUseCase interface {
	CreateAssessment(ctx context.Context, assessment *entity.Assessment) error
	InputGrade(ctx context.Context, grade *entity.Grade) error
	GetStudentGrades(ctx context.Context, studentID, classID, semesterID uuid.UUID) ([]*entity.Grade, error)
	CalculateFinalScore(ctx context.Context, studentID, classID, subjectID, semesterID uuid.UUID) (float64, error)
}
