package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
)

type TeacherAttendanceRepository interface {
	Create(ctx context.Context, attendance *entity.TeacherAttendance) error
	Update(ctx context.Context, attendance *entity.TeacherAttendance) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.TeacherAttendance, error)
	GetByTeacherAndDate(ctx context.Context, teacherID uuid.UUID, date time.Time) (*entity.TeacherAttendance, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*entity.TeacherAttendance, error)
}
