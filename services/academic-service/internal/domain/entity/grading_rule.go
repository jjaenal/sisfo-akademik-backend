package entity

import (
	"time"

	"github.com/google/uuid"
)

type GradingRule struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	TenantID     string     `json:"tenant_id" db:"tenant_id"`
	CurriculumID uuid.UUID  `json:"curriculum_id" db:"curriculum_id"`
	Grade        string     `json:"grade" db:"grade"`         // A, B, C, D, E
	MinScore     float64    `json:"min_score" db:"min_score"` // 85.0
	MaxScore     float64    `json:"max_score" db:"max_score"` // 100.0
	Points       float64    `json:"points" db:"points"`       // 4.0, 3.0, etc.
	Description  string     `json:"description" db:"description"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy    *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy    *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

func (g *GradingRule) Validate() map[string]string {
	errors := make(map[string]string)
	if g.CurriculumID == uuid.Nil {
		errors["curriculum_id"] = "Curriculum ID is required"
	}
	if g.Grade == "" {
		errors["grade"] = "Grade is required"
	}
	if g.MinScore < 0 {
		errors["min_score"] = "Min score cannot be negative"
	}
	if g.MaxScore < 0 {
		errors["max_score"] = "Max score cannot be negative"
	}
	if g.MinScore > g.MaxScore {
		errors["min_score"] = "Min score cannot be greater than max score"
	}
	return errors
}
