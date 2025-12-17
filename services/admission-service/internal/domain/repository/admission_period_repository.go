package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
)

type AdmissionPeriodRepository interface {
	Create(ctx context.Context, period *entity.AdmissionPeriod) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AdmissionPeriod, error)
	GetActive(ctx context.Context) (*entity.AdmissionPeriod, error)
	List(ctx context.Context) ([]*entity.AdmissionPeriod, error)
	Update(ctx context.Context, period *entity.AdmissionPeriod) error
	Delete(ctx context.Context, id uuid.UUID) error
}
