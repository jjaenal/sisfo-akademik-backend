package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReportUseCase_GetDailyRevenue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)

	u := NewReportUseCase(mockReportRepo, mockInvoiceRepo, mockPaymentRepo, 2*time.Second)

	tenantID := uuid.New()
	now := time.Now()
	revenue := []*entity.DailyRevenue{{Date: now.Format("2006-01-02"), Amount: 1000}}

	t.Run("success", func(t *testing.T) {
		mockReportRepo.EXPECT().GetDailyRevenue(gomock.Any(), tenantID, now, now).Return(revenue, nil)
		res, err := u.GetDailyRevenue(context.Background(), tenantID, now, now)
		assert.NoError(t, err)
		assert.Equal(t, revenue, res)
	})

	t.Run("error", func(t *testing.T) {
		mockReportRepo.EXPECT().GetDailyRevenue(gomock.Any(), tenantID, now, now).Return(nil, errors.New("db error"))
		res, err := u.GetDailyRevenue(context.Background(), tenantID, now, now)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestReportUseCase_GetMonthlyRevenue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)

	u := NewReportUseCase(mockReportRepo, mockInvoiceRepo, mockPaymentRepo, 2*time.Second)

	tenantID := uuid.New()
	year := 2023
	revenue := []*entity.MonthlyRevenue{{Month: "January", Amount: 1000}}

	t.Run("success", func(t *testing.T) {
		mockReportRepo.EXPECT().GetMonthlyRevenue(gomock.Any(), tenantID, year).Return(revenue, nil)
		res, err := u.GetMonthlyRevenue(context.Background(), tenantID, year)
		assert.NoError(t, err)
		assert.Equal(t, revenue, res)
	})
}

func TestReportUseCase_GetOutstandingInvoices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)

	u := NewReportUseCase(mockReportRepo, mockInvoiceRepo, mockPaymentRepo, 2*time.Second)

	tenantID := uuid.New()
	outstanding := []*entity.OutstandingReport{{StudentName: "John", Amount: 1000}}

	t.Run("success", func(t *testing.T) {
		mockReportRepo.EXPECT().GetOutstandingInvoices(gomock.Any(), tenantID).Return(outstanding, nil)
		res, err := u.GetOutstandingInvoices(context.Background(), tenantID)
		assert.NoError(t, err)
		assert.Equal(t, outstanding, res)
	})
}

func TestReportUseCase_GetStudentHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportRepository(ctrl)
	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)

	u := NewReportUseCase(mockReportRepo, mockInvoiceRepo, mockPaymentRepo, 2*time.Second)

	tenantID := uuid.New()
	studentID := uuid.New()
	invoiceID := uuid.New()

	invoices := []*entity.Invoice{{ID: invoiceID, Amount: 1000}}
	payments := []*entity.Payment{{InvoiceID: invoiceID, Amount: 1000}}

	t.Run("success", func(t *testing.T) {
		mockInvoiceRepo.EXPECT().List(gomock.Any(), tenantID, studentID, entity.InvoiceStatus("")).Return(invoices, nil)
		mockPaymentRepo.EXPECT().ListByInvoiceID(gomock.Any(), invoiceID).Return(payments, nil)

		res, err := u.GetStudentHistory(context.Background(), tenantID, studentID)
		assert.NoError(t, err)
		assert.Equal(t, invoices, res.Invoices)
		assert.Equal(t, payments, res.Payments)
	})

	t.Run("invoice error", func(t *testing.T) {
		mockInvoiceRepo.EXPECT().List(gomock.Any(), tenantID, studentID, entity.InvoiceStatus("")).Return(nil, errors.New("db error"))

		res, err := u.GetStudentHistory(context.Background(), tenantID, studentID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("payment error", func(t *testing.T) {
		mockInvoiceRepo.EXPECT().List(gomock.Any(), tenantID, studentID, entity.InvoiceStatus("")).Return(invoices, nil)
		mockPaymentRepo.EXPECT().ListByInvoiceID(gomock.Any(), invoiceID).Return(nil, errors.New("db error"))

		res, err := u.GetStudentHistory(context.Background(), tenantID, studentID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
