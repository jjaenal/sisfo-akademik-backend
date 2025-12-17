package entity

import (
	"time"

	"github.com/google/uuid"
)

type DailyRevenue struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type MonthlyRevenue struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

type OutstandingReport struct {
	StudentID     uuid.UUID `json:"student_id"`
	StudentName   string    `json:"student_name"`
	InvoiceNumber string    `json:"invoice_number"`
	Amount        float64   `json:"amount"`
	PaidAmount    float64   `json:"paid_amount"`
	Remaining     float64   `json:"remaining"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
}

type StudentHistory struct {
	Invoices []*Invoice `json:"invoices"`
	Payments []*Payment `json:"payments"`
}
