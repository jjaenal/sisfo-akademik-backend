package entity

import (
	"time"

	"github.com/google/uuid"
)

type Curriculum struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Year        int        `json:"year"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedBy   *uuid.UUID `json:"created_by"`
	UpdatedBy   *uuid.UUID `json:"updated_by"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (c *Curriculum) Validate() map[string]string {
	errors := make(map[string]string)
	if c.Name == "" {
		errors["name"] = "Name is required"
	}
	if c.Year <= 0 {
		errors["year"] = "Year must be valid"
	}
	return errors
}

type CurriculumSubject struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     string     `json:"tenant_id"`
	CurriculumID uuid.UUID  `json:"curriculum_id"`
	SubjectID    uuid.UUID  `json:"subject_id"`
	GradeLevel   int        `json:"grade_level"`
	Semester     int        `json:"semester"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedBy    *uuid.UUID `json:"created_by"`
	UpdatedBy    *uuid.UUID `json:"updated_by"`
	DeletedAt    *time.Time `json:"deleted_at"`

	// Relations
	Subject *Subject `json:"subject,omitempty"`
}

func (cs *CurriculumSubject) Validate() map[string]string {
	errors := make(map[string]string)
	if cs.CurriculumID == uuid.Nil {
		errors["curriculum_id"] = "Curriculum ID is required"
	}
	if cs.SubjectID == uuid.Nil {
		errors["subject_id"] = "Subject ID is required"
	}
	return errors
}
