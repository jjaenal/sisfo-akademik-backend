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

		err := u.RecordPayment(context.Background(), paymentExceed)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exceeds remaining amount")
	})
}
