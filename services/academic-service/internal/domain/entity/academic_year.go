package entity

import (
	"time"

	"github.com/google/uuid"
)

type AcademicYear struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	TenantID  string     `json:"tenant_id" db:"tenant_id"`
	Name      string     `json:"name" db:"name"`
	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   time.Time  `json:"end_date" db:"end_date"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (a *AcademicYear) Validate() map[string]string {
	errors := make(map[string]string)
	if a.Name == "" {
		errors["name"] = "Name is required"
	}
	if a.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	if a.StartDate.IsZero() {
		errors["start_date"] = "Start date is required"
	}
	if a.EndDate.IsZero() {
		errors["end_date"] = "End date is required"
	}
	if !a.StartDate.Before(a.EndDate) {
		errors["date_range"] = "Start date must be before end date"
	}
	return errors
}
