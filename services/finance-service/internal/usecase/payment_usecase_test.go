package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentUseCase_RecordPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewPaymentUseCase(mockPaymentRepo, mockInvoiceRepo, timeout)

	invoiceID := uuid.New()
	payment := &entity.Payment{
		InvoiceID:     invoiceID,
		Amount:        50000,
		PaymentMethod: "TRANSFER",
	}

	t.Run("success full payment", func(t *testing.T) {
		invoice := &entity.Invoice{
			ID:         invoiceID,
			Amount:     50000,
			PaidAmount: 0,
			Status:     entity.InvoiceStatusUnpaid,
		}

		mockInvoiceRepo.EXPECT().GetByID(gomock.Any(), invoiceID).Return(invoice, nil)
		mockPaymentRepo.EXPECT().Create(gomock.Any(), payment).Return(nil)
		
		// Expect invoice update to PAID
		mockInvoiceRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, inv *entity.Invoice) error {
			assert.Equal(t, entity.InvoiceStatusPaid, inv.Status)
			assert.Equal(t, 50000.0, inv.PaidAmount)
			return nil
		})

		err := u.RecordPayment(context.Background(), payment)
		assert.NoError(t, err)
	})

	t.Run("success partial payment", func(t *testing.T) {
		invoice := &entity.Invoice{
			ID:         invoiceID,
			Amount:     100000,
			PaidAmount: 0,
			Status:     entity.InvoiceStatusUnpaid,
		}

		mockInvoiceRepo.EXPECT().GetByID(gomock.Any(), invoiceID).Return(invoice, nil)
		mockPaymentRepo.EXPECT().Create(gomock.Any(), payment).Return(nil)

		// Expect invoice update to PARTIAL
		mockInvoiceRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, inv *entity.Invoice) error {
			assert.Equal(t, entity.InvoiceStatusPartial, inv.Status)
			assert.Equal(t, 50000.0, inv.PaidAmount)
			return nil
		})

		err := u.RecordPayment(context.Background(), payment)
		assert.NoError(t, err)
	})

	t.Run("invoice not found", func(t *testing.T) {
		mockInvoiceRepo.EXPECT().GetByID(gomock.Any(), invoiceID).Return(nil, errors.New("not found"))

		err := u.RecordPayment(context.Background(), payment)
		assert.Error(t, err)
	})

	t.Run("invoice already paid", func(t *testing.T) {
		invoice := &entity.Invoice{
			ID:     invoiceID,
			Status: entity.InvoiceStatusPaid,
		}

		mockInvoiceRepo.EXPECT().GetByID(gomock.Any(), invoiceID).Return(invoice, nil)

		err := u.RecordPayment(context.Background(), payment)
		assert.Error(t, err)
		assert.Equal(t, "invoice is already paid", err.Error())
	})

	t.Run("payment exceeds remaining", func(t *testing.T) {
		invoice := &entity.Invoice{
			ID:         invoiceID,
			Amount:     50000,
			PaidAmount: 20000, // Remaining 30000
		}

		paymentExceed := &entity.Payment{
			InvoiceID: invoiceID,
			Amount:    40000, // Exceeds 30000
		}

		mockInvoiceRepo.EXPECT().GetByID(gomock.Any(), invoiceID).Return(invoice, nil)
		// We expect NO call to paymentRepo.Create because it should fail before that

		err := u.RecordPayment(context.Background(), paymentExceed)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds remaining amount")
	})
}

func TestPaymentUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewPaymentUseCase(mockPaymentRepo, mockInvoiceRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		payment := &entity.Payment{ID: id}
		mockPaymentRepo.EXPECT().GetByID(gomock.Any(), id).Return(payment, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, payment, res)
	})

	t.Run("not found", func(t *testing.T) {
		mockPaymentRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("not found"))

		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestPaymentUseCase_ListByInvoiceID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewPaymentUseCase(mockPaymentRepo, mockInvoiceRepo, timeout)

	invoiceID := uuid.New()

	t.Run("success", func(t *testing.T) {
		payments := []*entity.Payment{{InvoiceID: invoiceID}}
		mockPaymentRepo.EXPECT().ListByInvoiceID(gomock.Any(), invoiceID).Return(payments, nil)

		res, err := u.ListByInvoiceID(context.Background(), invoiceID)
		assert.NoError(t, err)
		assert.Equal(t, payments, res)
	})

	t.Run("error", func(t *testing.T) {
		mockPaymentRepo.EXPECT().ListByInvoiceID(gomock.Any(), invoiceID).Return(nil, errors.New("db error"))

		res, err := u.ListByInvoiceID(context.Background(), invoiceID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
