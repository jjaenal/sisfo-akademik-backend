package entity

import (
	"time"

	"github.com/google/uuid"
)

type ScheduleTemplate struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	TenantID    string                 `json:"tenant_id" db:"tenant_id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	IsActive    bool                   `json:"is_active" db:"is_active"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	CreatedBy   *uuid.UUID             `json:"created_by" db:"created_by"`
	UpdatedBy   *uuid.UUID             `json:"updated_by" db:"updated_by"`
	DeletedAt   *time.Time             `json:"deleted_at" db:"deleted_at"`
	Items       []ScheduleTemplateItem `json:"items,omitempty" db:"-"` // One-to-many relationship
}

type ScheduleTemplateItem struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	TemplateID         uuid.UUID  `json:"template_id" db:"template_id"`
	SubjectID          *uuid.UUID `json:"subject_id" db:"subject_id"` // Optional
	DayOfWeek          int        `json:"day_of_week" db:"day_of_week"`
	StartTime          string     `json:"start_time" db:"start_time"` // Format: "15:04:05"
	EndTime            string     `json:"end_time" db:"end_time"`     // Format: "15:04:05"
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

func (t *ScheduleTemplate) Validate() map[string]string {
	errors := make(map[string]string)
	if t.TenantID == "" {
		errors["tenant_id"] = "Tenant ID is required"
	}
	if t.Name == "" {
		errors["name"] = "Name is required"
	}
	return errors
}

func (i *ScheduleTemplateItem) Validate() map[string]string {
	errors := make(map[string]string)
	if i.DayOfWeek < 1 || i.DayOfWeek > 7 {
		errors["day_of_week"] = "Day of week must be between 1 and 7"
	}
	if i.StartTime == "" {
		errors["start_time"] = "Start time is required"
	}
	if i.EndTime == "" {
		errors["end_time"] = "End time is required"
	}
	return errors
}
