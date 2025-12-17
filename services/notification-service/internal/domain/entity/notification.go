package entity

import (
	"time"

	"github.com/google/uuid"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSent    NotificationStatus = "SENT"
	NotificationStatusFailed  NotificationStatus = "FAILED"
)

type Notification struct {
	ID           uuid.UUID           `json:"id"`
	TemplateID   *uuid.UUID          `json:"template_id,omitempty"`
	Channel      NotificationChannel `json:"channel"`
	Recipient    string              `json:"recipient"` // Email address or Phone number
	Subject      string              `json:"subject"`   // For Email
	Body         string              `json:"body"`
	Status       NotificationStatus  `json:"status"`
	ErrorMessage string              `json:"error_message,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	SentAt       *time.Time          `json:"sent_at,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) Validate() map[string]string {
	errors := make(map[string]string)
	if n.Channel == "" {
		errors["channel"] = "Channel is required"
	}
	if n.Recipient == "" {
		errors["recipient"] = "Recipient is required"
	}
	if n.Body == "" {
		errors["body"] = "Body is required"
	}
	return errors
}
