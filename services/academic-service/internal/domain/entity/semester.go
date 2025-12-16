package entity

import (
	"time"

	"github.com/google/uuid"
)

type SemesterType string

const (
	SemesterOdd  SemesterType = "ODD"
	SemesterEven SemesterType = "EVEN"
)

type Semester struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	TenantID       string       `json:"tenant_id" db:"tenant_id"`
	AcademicYearID uuid.UUID    `json:"academic_year_id" db:"academic_year_id"`
	Name           string     `json:"name" db:"name"`
	StartDate      time.Time  `json:"start_date" db:"start_date"`
	EndDate        time.Time  `json:"end_date" db:"end_date"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy      *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy      *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt      *time.Time `json:"deleted_at" db:"deleted_at"`

	// Relations
	AcademicYear *AcademicYear `json:"academic_year,omitempty" db:"-"`
}

func (s *Semester) Validate() map[string]string {
	errors := make(map[string]string)
	if s.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	if s.AcademicYearID == uuid.Nil {
		errors["academic_year_id"] = "Academic Year ID is required"
	}
	if s.Name == "" {
		errors["name"] = "Name is required"
	}
	if s.StartDate.IsZero() {
		errors["start_date"] = "Start date is required"
	}
	if s.EndDate.IsZero() {
		errors["end_date"] = "End date is required"
	}
	if !s.StartDate.Before(s.EndDate) {
		errors["date_range"] = "Start date must be before end date"
	}
	return errors
}
