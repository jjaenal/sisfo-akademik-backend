package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/usecase"
)

type admissionPeriodUseCase struct {
	repo           repository.AdmissionPeriodRepository
	appRepo        repository.ApplicationRepository
	contextTimeout time.Duration
}

// Ensure interface implementation
var _ domainUseCase.AdmissionPeriodUseCase = (*admissionPeriodUseCase)(nil)

func NewAdmissionPeriodUseCase(repo repository.AdmissionPeriodRepository, appRepo repository.ApplicationRepository, timeout time.Duration) domainUseCase.AdmissionPeriodUseCase {
	return &admissionPeriodUseCase{
		repo:           repo,
		appRepo:        appRepo,
		contextTimeout: timeout,
	}
}

func (u *admissionPeriodUseCase) Create(ctx context.Context, period *entity.AdmissionPeriod) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := period.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, period)
}

func (u *admissionPeriodUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.AdmissionPeriod, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByID(ctx, id)
}

func (u *admissionPeriodUseCase) GetActive(ctx context.Context) (*entity.AdmissionPeriod, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetActive(ctx)
}

func (u *admissionPeriodUseCase) List(ctx context.Context) ([]*entity.AdmissionPeriod, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.List(ctx)
}

func (u *admissionPeriodUseCase) Update(ctx context.Context, period *entity.AdmissionPeriod) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := period.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, period)
}

func (u *admissionPeriodUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.Delete(ctx, id)
}

func (u *admissionPeriodUseCase) AnnounceResults(ctx context.Context, id uuid.UUID, passingGrade float64) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	period, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if period == nil {
		return errors.New("admission period not found")
	}

	if period.IsAnnounced {
		return errors.New("results already announced")
	}

	// Get applications
	apps, err := u.appRepo.List(ctx, map[string]interface{}{"admission_period_id": id})
	if err != nil {
		return err
	}

	for _, app := range apps {
		if app.FinalScore == nil {
			continue
		}

		if *app.FinalScore >= passingGrade {
			app.Status = entity.ApplicationStatusAccepted
		} else {
			app.Status = entity.ApplicationStatusRejected
		}

		if err := u.appRepo.Update(ctx, app); err != nil {
			return err
		}
	}

	period.IsAnnounced = true
	return u.repo.Update(ctx, period)
}
