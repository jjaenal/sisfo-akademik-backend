package entity

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceStatus string

const (
	InvoiceStatusUnpaid    InvoiceStatus = "UNPAID"
	InvoiceStatusPaid      InvoiceStatus = "PAID"
	InvoiceStatusPartial   InvoiceStatus = "PARTIAL"
	InvoiceStatusOverdue   InvoiceStatus = "OVERDUE"
	InvoiceStatusCancelled InvoiceStatus = "CANCELLED"
)

type Invoice struct {
	ID              uuid.UUID     `json:"id"`
	TenantID        uuid.UUID     `json:"tenant_id"`
	StudentID       uuid.UUID     `json:"student_id"`
	BillingConfigID uuid.UUID     `json:"billing_config_id"`
	InvoiceNumber   string        `json:"invoice_number"`
	Amount          float64       `json:"amount"`
	Status          InvoiceStatus `json:"status"`
	DueDate         time.Time     `json:"due_date"`
	PaidAmount      float64       `json:"paid_amount"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	// Relations
	BillingConfig *BillingConfig `json:"billing_config,omitempty"`
}

func (Invoice) TableName() string {
	return "invoices"
}

func (i *Invoice) Validate() map[string]string {
	errors := make(map[string]string)
	if i.StudentID == uuid.Nil {
		errors["student_id"] = "Student ID is required"
	}
	if i.Amount < 0 {
		errors["amount"] = "Amount cannot be negative"
	}
	if i.DueDate.IsZero() {
		errors["due_date"] = "Due date is required"
	}
	return errors
}
