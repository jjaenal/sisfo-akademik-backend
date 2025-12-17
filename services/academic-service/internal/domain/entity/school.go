package entity

import (
	"time"

	"github.com/google/uuid"
)

type School struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	TenantID      string     `json:"tenant_id" db:"tenant_id"`
	Name          string     `json:"name" db:"name"`
	Address       string     `json:"address" db:"address"`
	Phone         string     `json:"phone" db:"phone"`
	Email         string     `json:"email" db:"email"`
	Website       string     `json:"website" db:"website"`
	LogoURL       string     `json:"logo_url" db:"logo_url"`
	Latitude      float64    `json:"latitude" db:"latitude"`
	Longitude     float64    `json:"longitude" db:"longitude"`
	Accreditation string     `json:"accreditation" db:"accreditation"`
	Headmaster    string     `json:"headmaster" db:"headmaster"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy     *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy     *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (s *School) Validate() map[string]string {
	errors := make(map[string]string)
	if s.Name == "" {
		errors["name"] = "Name is required"
	}
	if s.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	return errors
}
