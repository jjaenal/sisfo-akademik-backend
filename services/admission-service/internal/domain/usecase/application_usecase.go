package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
)

type ApplicationUseCase interface {
	SubmitApplication(ctx context.Context, application *entity.Application) (*entity.Application, error)
	GetApplicationStatus(ctx context.Context, registrationNumber string) (*entity.Application, error)
	VerifyApplication(ctx context.Context, id uuid.UUID, status entity.ApplicationStatus) error
	ListApplications(ctx context.Context, filter map[string]interface{}) ([]*entity.Application, error)
	InputTestScore(ctx context.Context, id uuid.UUID, score float64) error
	InputInterviewScore(ctx context.Context, id uuid.UUID, score float64) error
	CalculateFinalScores(ctx context.Context, periodID uuid.UUID) error
	RegisterStudent(ctx context.Context, id uuid.UUID) error
}
