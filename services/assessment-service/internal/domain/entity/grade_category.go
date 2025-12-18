package entity

import (
	"time"

	"github.com/google/uuid"
)

type GradeCategory struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Weight      float64    `json:"weight"` // Percentage
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (e *GradeCategory) Validate() map[string]string {
	errs := make(map[string]string)
	if e.Name == "" {
		errs["name"] = "name is required"
	}
	if e.Weight <= 0 {
		errs["weight"] = "weight must be greater than 0"
	}
	return errs
}
