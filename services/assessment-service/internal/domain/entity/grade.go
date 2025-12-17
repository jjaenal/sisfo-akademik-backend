package entity

import (
	"time"

	"github.com/google/uuid"
)

type GradeStatus string

const (
	GradeStatusDraft     GradeStatus = "draft"
	GradeStatusSubmitted GradeStatus = "submitted"
	GradeStatusApproved  GradeStatus = "approved"
)

type Grade struct {
	ID           uuid.UUID   `json:"id"`
	AssessmentID uuid.UUID   `json:"assessment_id"`
	StudentID    uuid.UUID   `json:"student_id"`
	Score        float64     `json:"score"`
	Status       GradeStatus `json:"status"`
	Notes        string      `json:"notes"`
	GradedBy     uuid.UUID   `json:"graded_by"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (Grade) TableName() string {
	return "grades"
}

func (g *Grade) Validate() map[string]string {
	errors := make(map[string]string)
	if g.AssessmentID == uuid.Nil {
		errors["assessment_id"] = "Assessment ID is required"
	}
	if g.StudentID == uuid.Nil {
		errors["student_id"] = "Student ID is required"
	}
	if g.Score < 0 {
		errors["score"] = "Score cannot be negative"
	}
	return errors
}
