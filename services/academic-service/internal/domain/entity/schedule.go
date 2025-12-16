package entity

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	TenantID  string     `json:"tenant_id" db:"tenant_id"`
	ClassID   uuid.UUID  `json:"class_id" db:"class_id"`
	SubjectID uuid.UUID  `json:"subject_id" db:"subject_id"`
	TeacherID uuid.UUID  `json:"teacher_id" db:"teacher_id"`
	DayOfWeek int        `json:"day_of_week" db:"day_of_week"`
	StartTime string     `json:"start_time" db:"start_time"` // Format: "15:04:05" or "15:04"
	EndTime   string     `json:"end_time" db:"end_time"`     // Format: "15:04:05" or "15:04"
	Room      string     `json:"room" db:"room"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (s *Schedule) Validate() map[string]string {
	errors := make(map[string]string)
	if s.ClassID == uuid.Nil {
		errors["class_id"] = "Class ID is required"
	}
	if s.SubjectID == uuid.Nil {
		errors["subject_id"] = "Subject ID is required"
	}
	if s.TeacherID == uuid.Nil {
		errors["teacher_id"] = "Teacher ID is required"
	}
	if s.DayOfWeek < 1 || s.DayOfWeek > 7 {
		errors["day_of_week"] = "Day of week must be between 1 and 7"
	}
	if s.StartTime == "" {
		errors["start_time"] = "Start time is required"
	}
	if s.EndTime == "" {
		errors["end_time"] = "End time is required"
	}
	if s.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	return errors
}
