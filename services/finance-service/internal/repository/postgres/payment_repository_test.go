package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestPaymentRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewPaymentRepository(db)
	invoiceRepo := NewInvoiceRepository(db)
	billingRepo := NewBillingConfigRepository(db)

	tenantID := uuid.New()
	studentID := uuid.New()

	// Setup Billing Config
	billingConfig := &entity.BillingConfig{
		TenantID:  tenantID,
		Name:      "SPP Payment Test",
		Amount:    100000,
		Frequency: entity.BillingFrequencyMonthly,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := billingRepo.Create(context.Background(), billingConfig)
	assert.NoError(t, err)

	// Setup Invoice
	invoice := &entity.Invoice{
		TenantID:        tenantID,
		StudentID:       studentID,
		BillingConfigID: billingConfig.ID,
		InvoiceNumber:   fmt.Sprintf("INV/%d/PAY001", time.Now().UnixNano()),
		Amount:          100000,
		Status:          entity.InvoiceStatusUnpaid,
		DueDate:         time.Now().Add(7 * 24 * time.Hour),
		PaidAmount:      0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = invoiceRepo.Create(context.Background(), invoice)
	assert.NoError(t, err)

	// Create Payment
	payment := &entity.Payment{
		TenantID:        tenantID,
		InvoiceID:       invoice.ID,
		Amount:          100000,
		PaymentMethod:   entity.PaymentMethodTransfer,
		ReferenceNumber: "REF123456",
		Status:          entity.PaymentStatusSuccess,
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx := context.Background()
	err = repo.Create(ctx, payment)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, payment.ID)

	// GetByID
	fetched, err := repo.GetByID(ctx, payment.ID)
	assert.NoError(t, err)
	if fetched == nil {
		t.Fatal("fetched payment is nil")
	}
	assert.Equal(t, payment.ID, fetched.ID)
	assert.Equal(t, payment.InvoiceID, fetched.InvoiceID)
	assert.Equal(t, payment.Amount, fetched.Amount)

	// ListByInvoiceID
	// Add another payment (e.g. partial payment scenario, though above was full)
	payment2 := &entity.Payment{
		TenantID:        tenantID,
		InvoiceID:       invoice.ID,
		Amount:          50000,
		PaymentMethod:   entity.PaymentMethodCash,
		ReferenceNumber: "CASH001",
		Status:          entity.PaymentStatusSuccess,
		TransactionDate: time.Now().Add(-1 * time.Hour),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = repo.Create(ctx, payment2)
	assert.NoError(t, err)

	list, err := repo.ListByInvoiceID(ctx, invoice.ID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
	
	// Check ordering (DESC by transaction date)
	// payment2 is older (transaction date -1 hour) than payment (now)? 
	// Wait, payment.TransactionDate is Now(). payment2 is Now() - 1h.
	// So payment (newer) should be first.
	assert.Equal(t, payment.ID, list[0].ID)
}
