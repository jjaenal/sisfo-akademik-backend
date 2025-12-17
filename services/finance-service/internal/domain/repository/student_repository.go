package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
)

type StudentRepository interface {
	Create(ctx context.Context, student *entity.Student) error
	GetActive(ctx context.Context, tenantID uuid.UUID) ([]*entity.Student, error)
}
