package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClassSubject struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	TenantID  string     `json:"tenant_id" db:"tenant_id"`
	ClassID   uuid.UUID  `json:"class_id" db:"class_id"`
	SubjectID uuid.UUID  `json:"subject_id" db:"subject_id"`
	TeacherID *uuid.UUID `json:"teacher_id" db:"teacher_id"` // Teacher can be null initially
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`

	// Relations
	Class   *Class   `json:"class,omitempty" db:"-"`
	Subject *Subject `json:"subject,omitempty" db:"-"`
	Teacher *Teacher `json:"teacher,omitempty" db:"-"`
}

func (c *ClassSubject) TableName() string {
	return "class_subjects"
}
