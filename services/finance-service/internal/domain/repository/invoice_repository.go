package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
)

type InvoiceRepository interface {
	Create(ctx context.Context, invoice *entity.Invoice) error
	Update(ctx context.Context, invoice *entity.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error)
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*entity.Invoice, error)
	List(ctx context.Context, tenantID, studentID uuid.UUID, status entity.InvoiceStatus) ([]*entity.Invoice, error)
	Exists(ctx context.Context, studentID, billingConfigID uuid.UUID, month, year int) (bool, error)
	UpdateOverdueStatus(ctx context.Context) (int64, error)
}
