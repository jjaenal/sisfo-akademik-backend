package usecase

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
)

type ApplicationDocumentUseCase interface {
	Upload(ctx context.Context, applicationID uuid.UUID, documentType string, file *multipart.FileHeader) (*entity.ApplicationDocument, error)
	GetByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationDocument, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
