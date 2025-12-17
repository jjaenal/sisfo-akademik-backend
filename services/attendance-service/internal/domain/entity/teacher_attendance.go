package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeacherAttendanceStatus string

const (
	TeacherAttendanceStatusPresent TeacherAttendanceStatus = "present"
	TeacherAttendanceStatusAbsent  TeacherAttendanceStatus = "absent"
	TeacherAttendanceStatusLate    TeacherAttendanceStatus = "late"
	TeacherAttendanceStatusExcused TeacherAttendanceStatus = "excused"
	TeacherAttendanceStatusSick    TeacherAttendanceStatus = "sick"
)

type TeacherAttendance struct {
	ID             uuid.UUID               `json:"id"`
	TeacherID      uuid.UUID               `json:"teacher_id"`
	SemesterID     uuid.UUID               `json:"semester_id"`
	AttendanceDate time.Time               `json:"attendance_date"`
	CheckInTime    *time.Time              `json:"check_in_time"`
	CheckOutTime   *time.Time              `json:"check_out_time"`
	Status         TeacherAttendanceStatus `json:"status"`
	Notes          string                  `json:"notes"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

func (TeacherAttendance) TableName() string {
	return "teacher_attendance"
}

func (t *TeacherAttendance) Validate() map[string]string {
	errors := make(map[string]string)
	if t.TeacherID == uuid.Nil {
		errors["teacher_id"] = "Teacher ID is required"
	}
	if t.SemesterID == uuid.Nil {
		errors["semester_id"] = "Semester ID is required"
	}
	if t.Status == "" {
		errors["status"] = "Status is required"
	}
	return errors
}
