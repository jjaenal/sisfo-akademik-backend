package entity

import (
	"time"

	"github.com/google/uuid"
)

type Class struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	TenantID          string     `json:"tenant_id" db:"tenant_id"`
	SchoolID          *uuid.UUID `json:"school_id" db:"school_id"`
	AcademicYearID    *uuid.UUID `json:"academic_year_id" db:"academic_year_id"`
	Name              string     `json:"name" db:"name"`
	Level             int        `json:"level" db:"level"`
	Major             string     `json:"major" db:"major"`
	HomeroomTeacherID *uuid.UUID `json:"homeroom_teacher_id" db:"homeroom_teacher_id"`
	Capacity          int        `json:"capacity" db:"capacity"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy         *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy         *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt         *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (c *Class) Validate() map[string]string {
	errors := make(map[string]string)
	if c.Name == "" {
		errors["name"] = "Name is required"
	}
	if c.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	return errors
}
