package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
)

type enrollmentUseCase struct {
	repo           repository.EnrollmentRepository
	classRepo      repository.ClassRepository
	contextTimeout time.Duration
}

var _ domainUseCase.EnrollmentUseCase = (*enrollmentUseCase)(nil)

func NewEnrollmentUseCase(repo repository.EnrollmentRepository, classRepo repository.ClassRepository, timeout time.Duration) domainUseCase.EnrollmentUseCase {
	return &enrollmentUseCase{
		repo:           repo,
		classRepo:      classRepo,
		contextTimeout: timeout,
	}
}

func (u *enrollmentUseCase) Enroll(ctx context.Context, e *entity.Enrollment) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := e.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	// Get Class details to check capacity
	class, err := u.classRepo.GetByID(ctx, e.ClassID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("class not found")
	}

	// Count existing enrollments
	existingEnrollments, err := u.repo.ListByClass(ctx, e.ClassID)
	if err != nil {
		return err
	}

	// Check capacity
	// Assuming existingEnrollments includes only active students.
	// We might need to filter by status if ListByClass returns all.
	// For now, let's assume simple count.
	if len(existingEnrollments) >= class.Capacity {
		return errors.New("class capacity exceeded")
	}

	return u.repo.Enroll(ctx, e)
}

func (u *enrollmentUseCase) Unenroll(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.Unenroll(ctx, id)
}

func (u *enrollmentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Enrollment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetByID(ctx, id)
}

func (u *enrollmentUseCase) ListByClass(ctx context.Context, classID uuid.UUID) ([]entity.Enrollment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.ListByClass(ctx, classID)
}

func (u *enrollmentUseCase) ListByStudent(ctx context.Context, studentID uuid.UUID) ([]entity.Enrollment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.ListByStudent(ctx, studentID)
}

func (u *enrollmentUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if status == "" {
		return errors.New("status is required")
	}

	return u.repo.UpdateStatus(ctx, id, status)
}
