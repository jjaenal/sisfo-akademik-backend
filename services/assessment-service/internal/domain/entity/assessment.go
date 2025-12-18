package entity

import (
	"time"

	"github.com/google/uuid"
)

type Assessment struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	GradeCategoryID uuid.UUID  `json:"grade_category_id"`
	TeacherID       uuid.UUID  `json:"teacher_id"`
	SubjectID       uuid.UUID  `json:"subject_id"`
	ClassID         uuid.UUID  `json:"class_id"`
	SemesterID      uuid.UUID  `json:"semester_id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	MaxScore        float64    `json:"max_score"`
	Date            time.Time  `json:"date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`

	// Relationships
	GradeCategory *GradeCategory `json:"grade_category,omitempty"`
}

func (e *Assessment) Validate() map[string]string {
	errs := make(map[string]string)
	if e.Name == "" {
		errs["name"] = "name is required"
	}
	if e.MaxScore <= 0 {
		errs["max_score"] = "max_score must be greater than 0"
	}
	if e.Date.IsZero() {
		errs["date"] = "date is required"
	}
	return errs
}
