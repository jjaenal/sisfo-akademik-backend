package usecase

import (
	"context"
	"errors"
	"math"
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

	if attendance.CheckInLatitude != nil && attendance.CheckInLongitude != nil {
		if err := u.validateDistance(ctx, *attendance.CheckInLatitude, *attendance.CheckInLongitude, attendance.ClassID); err != nil {
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
			if err := u.validateDistance(ctx, *att.CheckInLatitude, *att.CheckInLongitude, att.ClassID); err != nil {
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

func (u *studentAttendanceUseCase) validateDistance(ctx context.Context, lat, lon float64, classID uuid.UUID) error {
	// TODO: Get school location for this class from Academic Service
	// For now, we assume a fixed location (e.g., school center) or skip if not available.
	// This is a placeholder.
	schoolLat := -6.2088 // Example Latitude
	schoolLon := 106.8456 // Example Longitude
	maxDist := 100.0 // meters

	dist := haversine(lat, lon, schoolLat, schoolLon)
	if dist > maxDist {
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
