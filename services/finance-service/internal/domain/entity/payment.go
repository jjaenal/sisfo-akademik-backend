package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
	PaymentStatusFailed  PaymentStatus = "FAILED"
)

type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "CASH"
	PaymentMethodTransfer PaymentMethod = "TRANSFER"
	PaymentMethodVA       PaymentMethod = "VA"
	PaymentMethodEWallet  PaymentMethod = "EWALLET"
)

type Payment struct {
	ID              uuid.UUID     `json:"id"`
	TenantID        uuid.UUID     `json:"tenant_id"`
	InvoiceID       uuid.UUID     `json:"invoice_id"`
	Amount          float64       `json:"amount"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	ReferenceNumber string        `json:"reference_number"`
	TransactionDate time.Time     `json:"transaction_date"`
	Status          PaymentStatus `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	// Relations
	Invoice *Invoice `json:"invoice,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

func (p *Payment) Validate() map[string]string {
	errors := make(map[string]string)
	if p.InvoiceID == uuid.Nil {
		errors["invoice_id"] = "Invoice ID is required"
	}
	if p.Amount <= 0 {
		errors["amount"] = "Amount must be greater than 0"
	}
	if p.PaymentMethod == "" {
		errors["payment_method"] = "Payment method is required"
	}
	return errors
}
