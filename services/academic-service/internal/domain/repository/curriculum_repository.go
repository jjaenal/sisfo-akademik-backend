package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type CurriculumRepository interface {
	Create(ctx context.Context, curriculum *entity.Curriculum) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Curriculum, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Curriculum, int, error)
	Update(ctx context.Context, curriculum *entity.Curriculum) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Curriculum Subjects
	AddSubject(ctx context.Context, cs *entity.CurriculumSubject) error
	RemoveSubject(ctx context.Context, id uuid.UUID) error
	ListSubjects(ctx context.Context, curriculumID uuid.UUID) ([]entity.CurriculumSubject, error)

	// Grading Rules
	AddGradingRule(ctx context.Context, rule *entity.GradingRule) error
	ListGradingRules(ctx context.Context, curriculumID uuid.UUID) ([]entity.GradingRule, error)
	DeleteGradingRule(ctx context.Context, id uuid.UUID) error
}
