package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
)

type TeacherAttendanceUseCase interface {
	CheckIn(ctx context.Context, attendance *entity.TeacherAttendance) error
	CheckOut(ctx context.Context, teacherID uuid.UUID, date time.Time, checkOutTime time.Time) error
	GetByTeacherAndDate(ctx context.Context, teacherID uuid.UUID, date time.Time) (*entity.TeacherAttendance, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*entity.TeacherAttendance, error)
}
