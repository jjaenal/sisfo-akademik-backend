package usecase

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/service"
	domainUseCase "github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/usecase"
)

type studentAttendanceUseCase struct {
	repo           repository.StudentAttendanceRepository
	schoolService  service.SchoolService
	contextTimeout time.Duration
}

// Ensure interface implementation
var _ domainUseCase.StudentAttendanceUseCase = (*studentAttendanceUseCase)(nil)

func NewStudentAttendanceUseCase(repo repository.StudentAttendanceRepository, schoolService service.SchoolService, timeout time.Duration) domainUseCase.StudentAttendanceUseCase {
	return &studentAttendanceUseCase{
		repo:           repo,
		schoolService:  schoolService,
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

	if attendance.CheckInLatitude != nil && attendance.CheckInLongitude != nil {
		if err := u.validateDistance(ctx, *attendance.CheckInLatitude, *attendance.CheckInLongitude, attendance.TenantID); err != nil {
			return err
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

func (u *studentAttendanceUseCase) BulkCreate(ctx context.Context, attendances []*entity.StudentAttendance) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if len(attendances) == 0 {
		return nil
	}

	// Validate all records
	for _, att := range attendances {
		if errMap := att.Validate(); len(errMap) > 0 {
			for _, v := range errMap {
				return errors.New(v)
			}
		}
		if att.CheckInLatitude != nil && att.CheckInLongitude != nil {
			if err := u.validateDistance(ctx, *att.CheckInLatitude, *att.CheckInLongitude, att.TenantID); err != nil {
				return err
			}
		}
	}

	return u.repo.BulkCreate(ctx, attendances)
}

func (u *studentAttendanceUseCase) GetSummary(ctx context.Context, studentID uuid.UUID, semesterID uuid.UUID) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	return u.repo.GetSummary(ctx, studentID, semesterID)
}

func (u *studentAttendanceUseCase) GetDailyReport(ctx context.Context, tenantID uuid.UUID, date time.Time) ([]*entity.StudentAttendance, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByTenantAndDate(ctx, tenantID, date)
}

func (u *studentAttendanceUseCase) GetMonthlyReport(ctx context.Context, tenantID uuid.UUID, month int, year int) ([]*entity.StudentAttendance, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// Calculate start and end date of the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	// Add one month to start date, then subtract one day to get the last day of the month
	endDate := startDate.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)

	return u.repo.GetByDateRange(ctx, tenantID, startDate, endDate, nil)
}

func (u *studentAttendanceUseCase) GetClassReport(ctx context.Context, tenantID uuid.UUID, classID uuid.UUID, startDate, endDate time.Time) ([]*entity.StudentAttendance, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.repo.GetByDateRange(ctx, tenantID, startDate, endDate, &classID)
}

func (u *studentAttendanceUseCase) validateDistance(ctx context.Context, lat, lon float64, tenantID string) error {
	if u.schoolService == nil {
		return nil
	}

	location, err := u.schoolService.GetLocation(ctx, tenantID)
	if err != nil {
		// If location service fails or not found, we currently skip validation
		// In strict mode, we might want to return error
		return nil
	}

	if location == nil {
		return nil
	}

	dist := haversine(lat, lon, location.Latitude, location.Longitude)
	// Use location radius if available, otherwise default to 100m
	radius := location.Radius
	if radius <= 0 {
		radius = 100.0
	}

	if dist > radius {
		return errors.New("location too far from school")
	}
	return nil
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters
	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
