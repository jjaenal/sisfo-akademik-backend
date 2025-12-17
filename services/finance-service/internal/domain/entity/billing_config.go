package entity

import (
	"time"

	"github.com/google/uuid"
)

type BillingFrequency string

const (
	BillingFrequencyMonthly BillingFrequency = "MONTHLY"
	BillingFrequencyOnce    BillingFrequency = "ONCE"
	BillingFrequencyYearly  BillingFrequency = "YEARLY"
)

type BillingConfig struct {
	ID        uuid.UUID        `json:"id"`
	TenantID  uuid.UUID        `json:"tenant_id"`
	Name      string           `json:"name"`
	Amount    float64          `json:"amount"`
	Frequency BillingFrequency `json:"frequency"`
	IsActive  bool             `json:"is_active"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (BillingConfig) TableName() string {
	return "billing_configurations"
}

func (b *BillingConfig) Validate() map[string]string {
	errors := make(map[string]string)
	if b.Name == "" {
		errors["name"] = "Name is required"
	}
	if b.Amount <= 0 {
		errors["amount"] = "Amount must be greater than 0"
	}
	if b.Frequency == "" {
		errors["frequency"] = "Frequency is required"
	}
	return errors
}
