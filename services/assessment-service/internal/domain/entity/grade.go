package entity

import (
	"time"

	"github.com/google/uuid"
)

type GradeStatus string

const (
	GradeStatusDraft     GradeStatus = "DRAFT"
	GradeStatusSubmitted GradeStatus = "SUBMITTED"
	GradeStatusFinal     GradeStatus = "FINAL"
)

type Grade struct {
	ID           uuid.UUID   `json:"id"`
	TenantID     string      `json:"tenant_id"`
	AssessmentID uuid.UUID   `json:"assessment_id"`
	StudentID    uuid.UUID   `json:"student_id"`
	Score        float64     `json:"score"`
	Feedback     string      `json:"feedback"`
	Notes        string      `json:"notes"`
	Status       GradeStatus `json:"status"`
	GradedBy     uuid.UUID   `json:"graded_by"`
	ApprovedBy   *uuid.UUID  `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time  `json:"approved_at,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`

	// Relationships
	Assessment *Assessment `json:"assessment,omitempty"`
}

func (e *Grade) Validate() map[string]string {
	errs := make(map[string]string)
	if e.Score < 0 {
		errs["score"] = "score cannot be negative"
	}
	return errs
}
