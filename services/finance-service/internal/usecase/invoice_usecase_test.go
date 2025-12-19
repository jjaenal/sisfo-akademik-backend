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
	mockStudentRepo := mocks.NewMockStudentRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewInvoiceUseCase(mockInvoiceRepo, mockBillingConfigRepo, mockStudentRepo, timeout)

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

func TestInvoiceUseCase_GenerateAllMonthlyInvoices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	mockBillingConfigRepo := mocks.NewMockBillingConfigRepository(ctrl)
	mockStudentRepo := mocks.NewMockStudentRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewInvoiceUseCase(mockInvoiceRepo, mockBillingConfigRepo, mockStudentRepo, timeout)

	t.Run("success", func(t *testing.T) {
		tenantID := uuid.New()
		studentID := uuid.New()
		configID := uuid.New()

		configs := []*entity.BillingConfig{
			{ID: configID, TenantID: tenantID, Amount: 100000},
		}
		students := []*entity.Student{
			{ID: studentID, TenantID: tenantID},
		}

		// 1. List active monthly configs
		mockBillingConfigRepo.EXPECT().ListAllActiveMonthly(gomock.Any()).Return(configs, nil)

		// 2. Get active students for tenant
		mockStudentRepo.EXPECT().GetActive(gomock.Any(), tenantID).Return(students, nil)

		// 3. Check existence (assume not exists)
		mockInvoiceRepo.EXPECT().Exists(gomock.Any(), studentID, configID, gomock.Any(), gomock.Any()).Return(false, nil)

		// 4. Generate call internals
		// Generate calls GetByID on configRepo
		mockBillingConfigRepo.EXPECT().GetByID(gomock.Any(), configID).Return(configs[0], nil)
		// Generate calls Create on invoiceRepo
		mockInvoiceRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		err := u.GenerateAllMonthlyInvoices(context.Background())
		assert.NoError(t, err)
	})

	t.Run("no active configs", func(t *testing.T) {
		mockBillingConfigRepo.EXPECT().ListAllActiveMonthly(gomock.Any()).Return([]*entity.BillingConfig{}, nil)

		err := u.GenerateAllMonthlyInvoices(context.Background())
		assert.NoError(t, err)
	})

	t.Run("list configs error", func(t *testing.T) {
		mockBillingConfigRepo.EXPECT().ListAllActiveMonthly(gomock.Any()).Return(nil, errors.New("db error"))

		err := u.GenerateAllMonthlyInvoices(context.Background())
		assert.Error(t, err)
	})
}

func TestInvoiceUseCase_CheckOverdueInvoices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewInvoiceUseCase(mockInvoiceRepo, nil, nil, timeout)

	t.Run("success", func(t *testing.T) {
		mockInvoiceRepo.EXPECT().UpdateOverdueStatus(gomock.Any()).Return(int64(5), nil)

		err := u.CheckOverdueInvoices(context.Background())
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockInvoiceRepo.EXPECT().UpdateOverdueStatus(gomock.Any()).Return(int64(0), errors.New("db error"))

		err := u.CheckOverdueInvoices(context.Background())
		assert.Error(t, err)
	})
}

func TestInvoiceUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewInvoiceUseCase(mockInvoiceRepo, nil, nil, timeout)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		expected := &entity.Invoice{ID: id}
		mockInvoiceRepo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error", func(t *testing.T) {
		id := uuid.New()
		mockInvoiceRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("db error"))

		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestInvoiceUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceRepo := mocks.NewMockInvoiceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewInvoiceUseCase(mockInvoiceRepo, nil, nil, timeout)

	t.Run("success", func(t *testing.T) {
		tenantID := uuid.New()
		studentID := uuid.New()
		status := entity.InvoiceStatusUnpaid
		expected := []*entity.Invoice{{ID: uuid.New()}}

		mockInvoiceRepo.EXPECT().List(gomock.Any(), tenantID, studentID, status).Return(expected, nil)

		res, err := u.List(context.Background(), tenantID, studentID, status)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error", func(t *testing.T) {
		tenantID := uuid.New()
		studentID := uuid.New()
		status := entity.InvoiceStatusUnpaid
		
		mockInvoiceRepo.EXPECT().List(gomock.Any(), tenantID, studentID, status).Return(nil, errors.New("db error"))

		res, err := u.List(context.Background(), tenantID, studentID, status)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
