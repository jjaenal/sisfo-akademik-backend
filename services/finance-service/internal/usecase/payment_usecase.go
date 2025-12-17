package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type PaymentUseCase interface {
	RecordPayment(ctx context.Context, payment *entity.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error)
	ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*entity.Payment, error)
}

type paymentUseCase struct {
	paymentRepo repository.PaymentRepository
	invoiceRepo repository.InvoiceRepository
	timeout     time.Duration
}

func NewPaymentUseCase(
	paymentRepo repository.PaymentRepository,
	invoiceRepo repository.InvoiceRepository,
	timeout time.Duration,
) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
		invoiceRepo: invoiceRepo,
		timeout: timeout,
	}
}

func (u *paymentUseCase) RecordPayment(ctx context.Context, payment *entity.Payment) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// 1. Get Invoice
	invoice, err := u.invoiceRepo.GetByID(ctx, payment.InvoiceID)
	if err != nil {
		return fmt.Errorf("failed to get invoice: %w", err)
	}

	// 2. Validate Payment
	if invoice.Status == entity.InvoiceStatusPaid {
		return errors.New("invoice is already paid")
	}
	if invoice.Status == entity.InvoiceStatusCancelled {
		return errors.New("cannot pay cancelled invoice")
	}

	remainingAmount := invoice.Amount - invoice.PaidAmount
	if payment.Amount > remainingAmount {
		return fmt.Errorf("payment amount %.2f exceeds remaining amount %.2f", payment.Amount, remainingAmount)
	}

	// 3. Create Payment
	if err := u.paymentRepo.Create(ctx, payment); err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// 4. Update Invoice Status
	invoice.PaidAmount += payment.Amount
	if invoice.PaidAmount >= invoice.Amount {
		invoice.Status = entity.InvoiceStatusPaid
	} else {
		invoice.Status = entity.InvoiceStatusPartial
	}

	if err := u.invoiceRepo.Update(ctx, invoice); err != nil {
		// Note: In a real system, we should rollback payment creation here or use transaction
		return fmt.Errorf("failed to update invoice status: %w", err)
	}

	return nil
}

func (u *paymentUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.paymentRepo.GetByID(ctx, id)
}

func (u *paymentUseCase) ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*entity.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.paymentRepo.ListByInvoiceID(ctx, invoiceID)
}
