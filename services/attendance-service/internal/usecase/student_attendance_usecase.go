package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/repository"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/usecase"
)

type studentAttendanceUseCase struct {
	repo           repository.StudentAttendanceRepository
	contextTimeout time.Duration
}

// Ensure interface implementation
var _ domainUseCase.StudentAttendanceUseCase = (*studentAttendanceUseCase)(nil)

func NewStudentAttendanceUseCase(repo repository.StudentAttendanceRepository, timeout time.Duration) domainUseCase.StudentAttendanceUseCase {
	return &studentAttendanceUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *studentAttendanceUseCase) Create(ctx context.Context, attendance *entity.StudentAttendance) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := attendance.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Create(ctx, attendance)
}

func (u *studentAttendanceUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.StudentAttendance, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByID(ctx, id)
}

func (u *studentAttendanceUseCase) GetByClassAndDate(ctx context.Context, classID uuid.UUID, date time.Time) ([]*entity.StudentAttendance, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByClassAndDate(ctx, classID, date)
}

func (u *studentAttendanceUseCase) Update(ctx context.Context, attendance *entity.StudentAttendance) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := attendance.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	return u.repo.Update(ctx, attendance)
}
