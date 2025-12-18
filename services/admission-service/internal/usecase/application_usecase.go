package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
)

type applicationUseCase struct {
	applicationRepo repository.ApplicationRepository
	rabbitClient    *rabbit.Client
	contextTimeout  time.Duration
}

func NewApplicationUseCase(applicationRepo repository.ApplicationRepository, rabbitClient *rabbit.Client, timeout time.Duration) domainUseCase.ApplicationUseCase {
	return &applicationUseCase{
		applicationRepo: applicationRepo,
		rabbitClient:    rabbitClient,
		contextTimeout:  timeout,
	}
}

func (u *applicationUseCase) SubmitApplication(ctx context.Context, application *entity.Application) (*entity.Application, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errs := application.Validate(); len(errs) > 0 {
		return nil, errors.New("validation error")
	}

	// Generate Registration Number
	// Format: REG-YYYYMMDD-XXXX
	now := time.Now()
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	application.RegistrationNumber = fmt.Sprintf("REG-%s-%04d", now.Format("20060102"), n.Int64())
	application.Status = entity.ApplicationStatusSubmitted
	application.SubmissionDate = &now

	err := u.applicationRepo.Create(ctx, application)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (u *applicationUseCase) GetApplicationStatus(ctx context.Context, registrationNumber string) (*entity.Application, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.applicationRepo.GetByRegistrationNumber(ctx, registrationNumber)
}

func (u *applicationUseCase) VerifyApplication(ctx context.Context, id uuid.UUID, status entity.ApplicationStatus) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	app, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}

	app.Status = status
	return u.applicationRepo.Update(ctx, app)
}

func (u *applicationUseCase) ListApplications(ctx context.Context, filter map[string]interface{}) ([]*entity.Application, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.applicationRepo.List(ctx, filter)
}

func (u *applicationUseCase) InputTestScore(ctx context.Context, id uuid.UUID, score float64) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	app, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}

	app.TestScore = &score
	return u.applicationRepo.Update(ctx, app)
}

func (u *applicationUseCase) InputInterviewScore(ctx context.Context, id uuid.UUID, score float64) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	app, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}

	app.InterviewScore = &score
	return u.applicationRepo.Update(ctx, app)
}

func (u *applicationUseCase) CalculateFinalScores(ctx context.Context, periodID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Get all applications for the period
	apps, err := u.applicationRepo.List(ctx, map[string]interface{}{"admission_period_id": periodID})
	if err != nil {
		return err
	}

	for _, app := range apps {
		if app.TestScore == nil || app.InterviewScore == nil {
			continue // Skip if scores are missing
		}

		// Simple formula: 40% Test, 40% Interview, 20% School Average
		test := *app.TestScore
		interview := *app.InterviewScore
		school := app.AverageScore

		final := (test * 0.4) + (interview * 0.4) + (school * 0.2)
		app.FinalScore = &final

		if err := u.applicationRepo.Update(ctx, app); err != nil {
			return err
		}
	}
	return nil
}

func (u *applicationUseCase) RegisterStudent(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	app, err := u.applicationRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("application not found")
	}

	if app.Status != entity.ApplicationStatusAccepted {
		return errors.New("application must be accepted to register")
	}

	if app.Status == entity.ApplicationStatusRegistered {
		return errors.New("student already registered")
	}

	// Publish event to create User and Student records
	eventPayload := map[string]interface{}{
		"tenant_id":           app.TenantID,
		"application_id":      app.ID,
		"registration_number": app.RegistrationNumber,
		"first_name":          app.FirstName,
		"last_name":           app.LastName,
		"email":               app.Email,
		"phone_number":        app.PhoneNumber,
		"timestamp":           time.Now(),
	}

	if u.rabbitClient != nil {
		if err := u.rabbitClient.PublishJSON("sisfo.events", "admission.student.registered", eventPayload); err != nil {
			return fmt.Errorf("failed to publish registration event: %w", err)
		}
	}

	app.Status = entity.ApplicationStatusRegistered
	return u.applicationRepo.Update(ctx, app)
}
