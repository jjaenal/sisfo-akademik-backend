package entity

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "present"
	AttendanceStatusAbsent  AttendanceStatus = "absent"
	AttendanceStatusLate    AttendanceStatus = "late"
	AttendanceStatusExcused AttendanceStatus = "excused"
	AttendanceStatusSick    AttendanceStatus = "sick"
)

type StudentAttendance struct {
	ID             uuid.UUID        `json:"id"`
	TenantID       string           `json:"tenant_id"`
	StudentID      uuid.UUID        `json:"student_id"`
	ClassID        uuid.UUID        `json:"class_id"`
	SemesterID       uuid.UUID        `json:"semester_id"`
	AttendanceDate   time.Time        `json:"attendance_date"`
	Status           AttendanceStatus `json:"status"`
	Notes            string           `json:"notes"`
	CheckInLatitude  *float64         `json:"check_in_latitude"`
	CheckInLongitude *float64         `json:"check_in_longitude"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

func (StudentAttendance) TableName() string {
	return "student_attendance"
}

func (s *StudentAttendance) Validate() map[string]string {
	errors := make(map[string]string)
	if s.StudentID == uuid.Nil {
		errors["student_id"] = "Student ID is required"
	}
	if s.ClassID == uuid.Nil {
		errors["class_id"] = "Class ID is required"
	}
	if s.SemesterID == uuid.Nil {
		errors["semester_id"] = "Semester ID is required"
	}
	if s.Status == "" {
		errors["status"] = "Status is required"
	}
	return errors
}
