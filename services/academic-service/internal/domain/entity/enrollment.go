package entity

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  string     `json:"tenant_id"`
	ClassID   uuid.UUID  `json:"class_id"`
	StudentID uuid.UUID  `json:"student_id"`
	Status    string     `json:"status"` // active, dropped, moved
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`

	// Relations
	Student *Student `json:"student,omitempty"`
	Class   *Class   `json:"class,omitempty"`
}

func (e *Enrollment) Validate() map[string]string {
	errors := make(map[string]string)
	if e.ClassID == uuid.Nil {
		errors["class_id"] = "Class ID is required"
	}
	if e.StudentID == uuid.Nil {
		errors["student_id"] = "Student ID is required"
	}
	return errors
}
