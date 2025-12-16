package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
)

type SubjectRepository interface {
	Create(ctx context.Context, subject *entity.Subject) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Subject, error)
	List(ctx context.Context, tenantID string, limit, offset int) ([]entity.Subject, int, error)
	Update(ctx context.Context, subject *entity.Subject) error
	Delete(ctx context.Context, id uuid.UUID) error
}
