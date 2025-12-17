package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
)

type ApplicationDocumentRepository interface {
	Create(ctx context.Context, doc *entity.ApplicationDocument) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ApplicationDocument, error)
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationDocument, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
