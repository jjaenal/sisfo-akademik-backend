package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type PaymentUseCaseMock struct {
	mock.Mock
}

func (m *PaymentUseCaseMock) RecordPayment(ctx context.Context, payment *entity.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *PaymentUseCaseMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Payment), args.Error(1)
}

func (m *PaymentUseCaseMock) ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*entity.Payment, error) {
	args := m.Called(ctx, invoiceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Payment), args.Error(1)
}
