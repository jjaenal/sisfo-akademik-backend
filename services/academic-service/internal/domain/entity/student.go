package entity

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	TenantID      string     `json:"tenant_id" db:"tenant_id"`
	UserID        *uuid.UUID `json:"user_id" db:"user_id"`
	NIS           string     `json:"nis" db:"nis"`
	NISN          string     `json:"nisn" db:"nisn"`
	Name          string     `json:"name" db:"name"`
	Gender        string     `json:"gender" db:"gender"`
	BirthPlace    string     `json:"birth_place" db:"birth_place"`
	BirthDate     *time.Time `json:"birth_date" db:"birth_date"`
	Address       string     `json:"address" db:"address"`
	Phone         string     `json:"phone" db:"phone"`
	Email         string     `json:"email" db:"email"`
	ParentName    string     `json:"parent_name" db:"parent_name"`
	ParentPhone   string     `json:"parent_phone" db:"parent_phone"`
	AdmissionDate *time.Time `json:"admission_date" db:"admission_date"`
	Status        string     `json:"status" db:"status"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy     *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy     *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (s *Student) Validate() map[string]string {
	errors := make(map[string]string)
	if s.Name == "" {
		errors["name"] = "Name is required"
	}
	if s.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	if s.Status == "" {
		s.Status = "active"
	}
	return errors
}
