package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
)

type StudentAttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.StudentAttendance) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.StudentAttendance, error)
	GetByClassAndDate(ctx context.Context, classID uuid.UUID, date time.Time) ([]*entity.StudentAttendance, error)
	Update(ctx context.Context, attendance *entity.StudentAttendance) error
	BulkCreate(ctx context.Context, attendances []*entity.StudentAttendance) error
	GetSummary(ctx context.Context, studentID uuid.UUID, semesterID uuid.UUID) (map[string]int, error)
}
