package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type paymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) repository.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	query := `
		INSERT INTO payments (
			id, tenant_id, invoice_id, amount, payment_method, 
			reference_number, transaction_date, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, $10
		)
	`
	if payment.ID == uuid.Nil {
		payment.ID = uuid.New()
	}
	now := time.Now()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = now
	}
	if payment.TransactionDate.IsZero() {
		payment.TransactionDate = now
	}

	_, err := r.db.Exec(ctx, query,
		payment.ID, payment.TenantID, payment.InvoiceID, payment.Amount, payment.PaymentMethod,
		payment.ReferenceNumber, payment.TransactionDate, payment.Status, payment.CreatedAt, payment.UpdatedAt,
	)
	return err
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Payment, error) {
	query := `
		SELECT id, tenant_id, invoice_id, amount, payment_method, 
		       reference_number, transaction_date, status, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	var payment entity.Payment
	err := r.db.QueryRow(ctx, query, id).Scan(
		&payment.ID, &payment.TenantID, &payment.InvoiceID, &payment.Amount, &payment.PaymentMethod,
		&payment.ReferenceNumber, &payment.TransactionDate, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) ListByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]*entity.Payment, error) {
	query := `
		SELECT id, tenant_id, invoice_id, amount, payment_method, 
		       reference_number, transaction_date, status, created_at, updated_at
		FROM payments
		WHERE invoice_id = $1
		ORDER BY transaction_date DESC
	`
	rows, err := r.db.Query(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		var payment entity.Payment
		if err := rows.Scan(
			&payment.ID, &payment.TenantID, &payment.InvoiceID, &payment.Amount, &payment.PaymentMethod,
			&payment.ReferenceNumber, &payment.TransactionDate, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}
	return payments, nil
}
