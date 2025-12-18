package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
)

type GradeRepository interface {
	Create(ctx context.Context, grade *entity.Grade) error
	CreateBulk(ctx context.Context, grades []*entity.Grade) error
	GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) ([]*entity.Grade, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Grade, error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) ([]*entity.Grade, error)
	GetByStudentAndAssessment(ctx context.Context, studentID, assessmentID uuid.UUID) (*entity.Grade, error)
	Update(ctx context.Context, grade *entity.Grade) error
}
