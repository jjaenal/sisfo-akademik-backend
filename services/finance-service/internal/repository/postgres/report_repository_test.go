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

func TestReportRepository_Revenue(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewReportRepository(db)
	paymentRepo := NewPaymentRepository(db)
	invoiceRepo := NewInvoiceRepository(db)
	billingRepo := NewBillingConfigRepository(db)

	tenantID := uuid.New()
	studentID := uuid.New()

	// Setup Billing & Invoice
	billingConfig := &entity.BillingConfig{
		TenantID:  tenantID,
		Name:      "Revenue Test",
		Amount:    100000,
		Frequency: entity.BillingFrequencyMonthly,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = billingRepo.Create(context.Background(), billingConfig)

	invoice := &entity.Invoice{
		TenantID:        tenantID,
		StudentID:       studentID,
		BillingConfigID: billingConfig.ID,
		InvoiceNumber:   fmt.Sprintf("INV/REV/%d", time.Now().UnixNano()),
		Amount:          100000,
		Status:          entity.InvoiceStatusPaid,
		DueDate:         time.Now(),
		PaidAmount:      100000,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = invoiceRepo.Create(context.Background(), invoice)

	// Create Payments
	// Payment 1: Today
	payment1 := &entity.Payment{
		TenantID:        tenantID,
		InvoiceID:       invoice.ID,
		Amount:          50000,
		PaymentMethod:   entity.PaymentMethodTransfer,
		ReferenceNumber: "REV1",
		Status:          entity.PaymentStatusSuccess,
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = paymentRepo.Create(context.Background(), payment1)

	// Payment 2: Today
	payment2 := &entity.Payment{
		TenantID:        tenantID,
		InvoiceID:       invoice.ID,
		Amount:          50000,
		PaymentMethod:   entity.PaymentMethodTransfer,
		ReferenceNumber: "REV2",
		Status:          entity.PaymentStatusSuccess,
		TransactionDate: time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = paymentRepo.Create(context.Background(), payment2)

	ctx := context.Background()
	
	// Test Daily Revenue
	startDate := time.Now().UTC().Add(-24 * time.Hour)
	endDate := time.Now().UTC().Add(24 * time.Hour)
	daily, err := repo.GetDailyRevenue(ctx, tenantID, startDate, endDate)
	assert.NoError(t, err)
	assert.NotEmpty(t, daily)
	// We expect sum of 100000 for today
	todayStr := time.Now().UTC().Format("2006-01-02")
	foundToday := false
	for _, d := range daily {
		if d.Date == todayStr {
			foundToday = true
			assert.Equal(t, 100000.0, d.Amount)
		}
	}
	assert.True(t, foundToday)

	// Test Monthly Revenue
	monthly, err := repo.GetMonthlyRevenue(ctx, tenantID, time.Now().Year())
	assert.NoError(t, err)
	assert.NotEmpty(t, monthly)
	
	monthStr := time.Now().Format("2006-01")
	foundMonth := false
	for _, m := range monthly {
		if m.Month == monthStr {
			foundMonth = true
			assert.Equal(t, 100000.0, m.Amount)
		}
	}
	assert.True(t, foundMonth)
}

func TestReportRepository_OutstandingInvoices(t *testing.T) {
	db := testDB(t)
	ensureMigrations(t, db)
	repo := NewReportRepository(db)
	studentRepo := NewStudentRepository(db)
	invoiceRepo := NewInvoiceRepository(db)
	billingRepo := NewBillingConfigRepository(db)

	tenantID := uuid.New()
	studentID := uuid.New()

	// Setup Student
	student := &entity.Student{
		ID:        studentID,
		TenantID:  tenantID,
		Name:      "Outstanding Student",
		Status:    entity.StudentStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := studentRepo.Create(context.Background(), student)
	assert.NoError(t, err)

	// Setup Billing
	billingConfig := &entity.BillingConfig{
		TenantID:  tenantID,
		Name:      "Outstanding Test",
		Amount:    500000,
		Frequency: entity.BillingFrequencyMonthly,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = billingRepo.Create(context.Background(), billingConfig)

	// Setup Invoice (Unpaid)
	invoice := &entity.Invoice{
		TenantID:        tenantID,
		StudentID:       studentID,
		BillingConfigID: billingConfig.ID,
		InvoiceNumber:   fmt.Sprintf("INV/OUT/%d", time.Now().UnixNano()),
		Amount:          500000,
		Status:          entity.InvoiceStatusUnpaid,
		DueDate:         time.Now().Add(7 * 24 * time.Hour),
		PaidAmount:      0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err = invoiceRepo.Create(context.Background(), invoice)
	assert.NoError(t, err)

	// Test GetOutstandingInvoices
	ctx := context.Background()
	outstanding, err := repo.GetOutstandingInvoices(ctx, tenantID)
	assert.NoError(t, err)
	assert.NotEmpty(t, outstanding)
	
	found := false
	for _, o := range outstanding {
		if o.InvoiceNumber == invoice.InvoiceNumber {
			found = true
			assert.Equal(t, student.Name, o.StudentName)
			assert.Equal(t, 500000.0, o.Remaining)
			assert.Equal(t, string(entity.InvoiceStatusUnpaid), o.Status)
			break
		}
	}
	assert.True(t, found)
}
