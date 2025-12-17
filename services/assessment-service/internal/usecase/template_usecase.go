package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
)

type templateUseCase struct {
	repo repository.TemplateRepository
}

func NewTemplateUseCase(repo repository.TemplateRepository) domainUseCase.TemplateUseCase {
	return &templateUseCase{
		repo: repo,
	}
}

func (u *templateUseCase) Create(ctx context.Context, template *entity.ReportCardTemplate) error {
	if template.Name == "" {
		return errors.New("template name is required")
	}

	template.ID = uuid.New()
	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now

	// If this is the first template for the tenant, make it default
	existing, err := u.repo.GetByTenantID(ctx, template.TenantID)
	if err == nil && len(existing) == 0 {
		template.IsDefault = true
	}

	return u.repo.Create(ctx, template)
}

func (u *templateUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.ReportCardTemplate, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *templateUseCase) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*entity.ReportCardTemplate, error) {
	return u.repo.GetByTenantID(ctx, tenantID)
}

func (u *templateUseCase) Update(ctx context.Context, template *entity.ReportCardTemplate) error {
	existing, err := u.repo.GetByID(ctx, template.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("template not found")
	}

	if template.Name != "" {
		existing.Name = template.Name
	}
	// Update config fields if provided (simplified for now, assumes full config replacement)
	existing.Config = template.Config
	existing.IsDefault = template.IsDefault

	return u.repo.Update(ctx, existing)
}

func (u *templateUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
