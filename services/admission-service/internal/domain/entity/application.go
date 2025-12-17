package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string

const (
	ApplicationStatusDraft     ApplicationStatus = "draft"
	ApplicationStatusSubmitted ApplicationStatus = "submitted"
	ApplicationStatusVerified  ApplicationStatus = "verified"
	ApplicationStatusAccepted  ApplicationStatus = "accepted"
	ApplicationStatusRejected  ApplicationStatus = "rejected"
	ApplicationStatusRegistered ApplicationStatus = "registered"
)

type Application struct {
	ID                 uuid.UUID         `json:"id"`
	TenantID           string            `json:"tenant_id"`
	AdmissionPeriodID  uuid.UUID         `json:"admission_period_id"`
	RegistrationNumber string            `json:"registration_number"`
	FirstName          string            `json:"first_name"`
	LastName           string            `json:"last_name"`
	Email              string            `json:"email"`
	PhoneNumber        string            `json:"phone_number"`
	Status             ApplicationStatus `json:"status"`
	PreviousSchool     string            `json:"previous_school"`
	AverageScore       float64           `json:"average_score"`
	TestScore          *float64          `json:"test_score"`
	InterviewScore     *float64          `json:"interview_score"`
	FinalScore         *float64          `json:"final_score"`
	SubmissionDate     *time.Time        `json:"submission_date"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

func (Application) TableName() string {
	return "applications"
}

func (a *Application) Validate() map[string]string {
	errors := make(map[string]string)
	if a.AdmissionPeriodID == uuid.Nil {
		errors["admission_period_id"] = "Admission Period ID is required"
	}
	if a.FirstName == "" {
		errors["first_name"] = "First name is required"
	}
	if a.Email == "" {
		errors["email"] = "Email is required"
	}
	return errors
}
