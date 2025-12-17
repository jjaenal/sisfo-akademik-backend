package entity

import (
	"time"

	"github.com/google/uuid"
)

type AdmissionPeriod struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
	IsAnnounced bool    `json:"is_announced"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AdmissionPeriod) TableName() string {
	return "admission_periods"
}

func (a *AdmissionPeriod) Validate() map[string]string {
	errors := make(map[string]string)
	if a.Name == "" {
		errors["name"] = "Name is required"
	}
	if a.StartDate.IsZero() {
		errors["start_date"] = "Start date is required"
	}
	if a.EndDate.IsZero() {
		errors["end_date"] = "End date is required"
	}
	if a.EndDate.Before(a.StartDate) {
		errors["end_date"] = "End date must be after start date"
	}
	return errors
}
