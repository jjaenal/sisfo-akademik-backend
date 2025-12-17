package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationChannel string

const (
	NotificationChannelEmail    NotificationChannel = "EMAIL"
	NotificationChannelWhatsApp NotificationChannel = "WHATSAPP"
)

type NotificationTemplate struct {
	ID              uuid.UUID           `json:"id"`
	Name            string              `json:"name"`
	Channel         NotificationChannel `json:"channel"`
	SubjectTemplate string              `json:"subject_template"` // For Email
	BodyTemplate    string              `json:"body_template"`
	IsActive        bool                `json:"is_active"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

func (t *NotificationTemplate) Validate() map[string]string {
	errors := make(map[string]string)
	if t.Name == "" {
		errors["name"] = "Name is required"
	}
	if t.Channel == "" {
		errors["channel"] = "Channel is required"
	}
	if t.BodyTemplate == "" {
		errors["body_template"] = "Body template is required"
	}
	return errors
}
