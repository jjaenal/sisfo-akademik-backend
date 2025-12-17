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

type teacherAttendanceUseCase struct {
	repo           repository.TeacherAttendanceRepository
	contextTimeout time.Duration
}

var _ domainUseCase.TeacherAttendanceUseCase = (*teacherAttendanceUseCase)(nil)

func NewTeacherAttendanceUseCase(repo repository.TeacherAttendanceRepository, timeout time.Duration) domainUseCase.TeacherAttendanceUseCase {
	return &teacherAttendanceUseCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *teacherAttendanceUseCase) CheckIn(ctx context.Context, attendance *entity.TeacherAttendance) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if errMap := attendance.Validate(); len(errMap) > 0 {
		for _, v := range errMap {
			return errors.New(v)
		}
	}

	// Check if already checked in
	existing, err := u.repo.GetByTeacherAndDate(ctx, attendance.TeacherID, attendance.AttendanceDate)
	if err == nil && existing != nil {
		return errors.New("already checked in for this date")
	}

	return u.repo.Create(ctx, attendance)
}

func (u *teacherAttendanceUseCase) CheckOut(ctx context.Context, teacherID uuid.UUID, date time.Time, checkOutTime time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	attendance, err := u.repo.GetByTeacherAndDate(ctx, teacherID, date)
	if err != nil {
		return err
	}
	if attendance == nil {
		return errors.New("attendance record not found")
	}

	attendance.CheckOutTime = &checkOutTime
	attendance.UpdatedAt = time.Now()

	return u.repo.Update(ctx, attendance)
}

func (u *teacherAttendanceUseCase) GetByTeacherAndDate(ctx context.Context, teacherID uuid.UUID, date time.Time) (*entity.TeacherAttendance, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByTeacherAndDate(ctx, teacherID, date)
}

func (u *teacherAttendanceUseCase) List(ctx context.Context, filter map[string]interface{}) ([]*entity.TeacherAttendance, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.List(ctx, filter)
}
