package entity

import (
	"time"

	"github.com/google/uuid"
)

type Assessment struct {
	ID              uuid.UUID `json:"id"`
	TeacherID       uuid.UUID `json:"teacher_id"`
	SubjectID       uuid.UUID `json:"subject_id"`
	ClassID         uuid.UUID `json:"class_id"`
	GradeCategoryID uuid.UUID `json:"grade_category_id"`
	Name            string    `json:"name"`
	Date            time.Time `json:"date"`
	MaxScore        int       `json:"max_score"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Assessment) TableName() string {
	return "assessments"
}

func (a *Assessment) Validate() map[string]string {
	errors := make(map[string]string)
	if a.TeacherID == uuid.Nil {
		errors["teacher_id"] = "Teacher ID is required"
	}
	if a.SubjectID == uuid.Nil {
		errors["subject_id"] = "Subject ID is required"
	}
	if a.ClassID == uuid.Nil {
		errors["class_id"] = "Class ID is required"
	}
	if a.GradeCategoryID == uuid.Nil {
		errors["grade_category_id"] = "Grade Category ID is required"
	}
	if a.Name == "" {
		errors["name"] = "Name is required"
	}
	if a.MaxScore <= 0 {
		errors["max_score"] = "Max score must be greater than 0"
	}
	return errors
}
