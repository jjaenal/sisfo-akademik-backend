package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subject struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	TenantID    string     `json:"tenant_id" db:"tenant_id"`
	Code        string     `json:"code" db:"code"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
	CreditUnits int        `json:"credit_units" db:"credit_units"`
	Type        string     `json:"type" db:"type"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy   *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy   *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (s *Subject) Validate() map[string]string {
	errors := make(map[string]string)
	if s.Name == "" {
		errors["name"] = "Name is required"
	}
	if s.Code == "" {
		errors["code"] = "Code is required"
	}
	if s.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	return errors
}
