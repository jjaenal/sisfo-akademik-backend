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

func TestInvoiceRepository_CRUD(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewInvoiceRepository(db)
	billingRepo := NewBillingConfigRepository(db)

	tenantID := uuid.New()
	studentID := uuid.New()

	// Setup Billing Config
	billingConfig := &entity.BillingConfig{
		TenantID:  tenantID,
		Name:      "SPP Test",
		Amount:    100000,
		Frequency: entity.BillingFrequencyMonthly,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := billingRepo.Create(context.Background(), billingConfig)
	assert.NoError(t, err)

	// Create Invoice
	invoice := &entity.Invoice{
		TenantID:        tenantID,
		StudentID:       studentID,
		BillingConfigID: billingConfig.ID,
		InvoiceNumber:   fmt.Sprintf("INV/%d", time.Now().UnixNano()),
		Amount:          100000,
		Status:          entity.InvoiceStatusUnpaid,
		DueDate:         time.Now().Add(7 * 24 * time.Hour),
		PaidAmount:      0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx := context.Background()
	err = repo.Create(ctx, invoice)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, invoice.ID)

	// GetByID
	fetched, err := repo.GetByID(ctx, invoice.ID)
	assert.NoError(t, err)
	if fetched == nil {
		t.Fatal("fetched invoice is nil")
	}
	assert.Equal(t, invoice.ID, fetched.ID)
	assert.Equal(t, invoice.InvoiceNumber, fetched.InvoiceNumber)

	// GetByInvoiceNumber
	fetchedByNo, err := repo.GetByInvoiceNumber(ctx, invoice.InvoiceNumber)
	assert.NoError(t, err)
	assert.Equal(t, invoice.ID, fetchedByNo.ID)

	// Update
	invoice.Status = entity.InvoiceStatusPartial
	invoice.PaidAmount = 50000
	err = repo.Update(ctx, invoice)
	assert.NoError(t, err)

	fetchedAfterUpdate, err := repo.GetByID(ctx, invoice.ID)
	assert.NoError(t, err)
	assert.Equal(t, entity.InvoiceStatusPartial, fetchedAfterUpdate.Status)
	assert.Equal(t, 50000.0, fetchedAfterUpdate.PaidAmount)

	// List
	list, err := repo.List(ctx, tenantID, uuid.Nil, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, list)

	listByStudent, err := repo.List(ctx, tenantID, studentID, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, listByStudent)
	assert.Equal(t, invoice.ID, listByStudent[0].ID)

	// Exists
	exists, err := repo.Exists(ctx, studentID, billingConfig.ID, int(time.Now().Month()), time.Now().Year())
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestInvoiceRepository_UpdateOverdueStatus(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewInvoiceRepository(db)
	billingRepo := NewBillingConfigRepository(db)

	tenantID := uuid.New()
	studentID := uuid.New()

	// Setup Billing Config
	billingConfig := &entity.BillingConfig{
		TenantID:  tenantID,
		Name:      "SPP Overdue",
		Amount:    100000,
		Frequency: entity.BillingFrequencyMonthly,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := billingRepo.Create(context.Background(), billingConfig)
	assert.NoError(t, err)

	// Create Overdue Invoice
	invoice := &entity.Invoice{
		TenantID:        tenantID,
		StudentID:       studentID,
		BillingConfigID: billingConfig.ID,
		InvoiceNumber:   fmt.Sprintf("INV/%d/OVERDUE", time.Now().UnixNano()),
		Amount:          100000,
		Status:          entity.InvoiceStatusUnpaid,
		DueDate:         time.Now().Add(-24 * time.Hour), // Yesterday
		PaidAmount:      0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	ctx := context.Background()
	err = repo.Create(ctx, invoice)
	assert.NoError(t, err)

	// Run UpdateOverdueStatus
	count, err := repo.UpdateOverdueStatus(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, count, int64(1))

	// Verify status
	fetched, err := repo.GetByID(ctx, invoice.ID)
	assert.NoError(t, err)
	if fetched == nil {
		t.Fatal("fetched invoice is nil")
	}
	assert.Equal(t, entity.InvoiceStatusOverdue, fetched.Status)
}
