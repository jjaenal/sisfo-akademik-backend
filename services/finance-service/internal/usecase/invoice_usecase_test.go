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

func TestInvoiceUseCase_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	mockBillingConfigRepo := mocks.NewMockBillingConfigRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewInvoiceUseCase(mockInvoiceRepo, mockBillingConfigRepo, timeout)

	tenantID := uuid.New()
	studentID := uuid.New()
	billingConfigID := uuid.New()

	t.Run("success", func(t *testing.T) {
		config := &entity.BillingConfig{
			ID:     billingConfigID,
			Amount: 100000,
		}

		mockBillingConfigRepo.EXPECT().GetByID(gomock.Any(), billingConfigID).Return(config, nil)

		mockInvoiceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, invoice *entity.Invoice) error {
			assert.Equal(t, tenantID, invoice.TenantID)
			assert.Equal(t, studentID, invoice.StudentID)
			assert.Equal(t, billingConfigID, invoice.BillingConfigID)
			assert.Equal(t, config.Amount, invoice.Amount)
			assert.Equal(t, entity.InvoiceStatusUnpaid, invoice.Status)
			assert.NotEmpty(t, invoice.InvoiceNumber)
			return nil
		})

		invoice, err := u.Generate(context.Background(), tenantID, studentID, billingConfigID)
		assert.NoError(t, err)
		assert.NotNil(t, invoice)
	})

	t.Run("billing config not found", func(t *testing.T) {
		mockBillingConfigRepo.EXPECT().GetByID(gomock.Any(), billingConfigID).Return(nil, errors.New("not found"))

		invoice, err := u.Generate(context.Background(), tenantID, studentID, billingConfigID)
		assert.Error(t, err)
		assert.Nil(t, invoice)
	})
}
