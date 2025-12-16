package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type ClassSubjectRepository interface {
	Create(ctx context.Context, classSubject *entity.ClassSubject) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ClassSubject, error)
	GetByClassAndSubject(ctx context.Context, classID, subjectID uuid.UUID) (*entity.ClassSubject, error)
	ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.ClassSubject, error)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]entity.ClassSubject, error)
	Update(ctx context.Context, classSubject *entity.ClassSubject) error
	Delete(ctx context.Context, id uuid.UUID) error
}
