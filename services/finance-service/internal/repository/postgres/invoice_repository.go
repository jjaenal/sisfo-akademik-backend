package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type invoiceRepository struct {
	db *pgxpool.Pool
}

func NewInvoiceRepository(db *pgxpool.Pool) repository.InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *entity.Invoice) error {
	query := `
		INSERT INTO invoices (
			id, tenant_id, student_id, billing_config_id, invoice_number, 
			amount, status, due_date, paid_amount, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, $10, $11
		)
	`
	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}
	now := time.Now()
	if invoice.CreatedAt.IsZero() {
		invoice.CreatedAt = now
	}
	if invoice.UpdatedAt.IsZero() {
		invoice.UpdatedAt = now
	}

	_, err := r.db.Exec(ctx, query,
		invoice.ID, invoice.TenantID, invoice.StudentID, invoice.BillingConfigID, invoice.InvoiceNumber,
		invoice.Amount, invoice.Status, invoice.DueDate, invoice.PaidAmount, invoice.CreatedAt, invoice.UpdatedAt,
	)
	return err
}

func (r *invoiceRepository) Update(ctx context.Context, invoice *entity.Invoice) error {
	query := `
		UPDATE invoices
		SET status = $1, paid_amount = $2, updated_at = $3
		WHERE id = $4
	`
	invoice.UpdatedAt = time.Now()
	tag, err := r.db.Exec(ctx, query,
		invoice.Status, invoice.PaidAmount, invoice.UpdatedAt, invoice.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("invoice not found")
	}
	return nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error) {
	query := `
		SELECT id, tenant_id, student_id, billing_config_id, invoice_number, 
		       amount, status, due_date, paid_amount, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`
	var invoice entity.Invoice
	err := r.db.QueryRow(ctx, query, id).Scan(
		&invoice.ID, &invoice.TenantID, &invoice.StudentID, &invoice.BillingConfigID, &invoice.InvoiceNumber,
		&invoice.Amount, &invoice.Status, &invoice.DueDate, &invoice.PaidAmount, &invoice.CreatedAt, &invoice.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*entity.Invoice, error) {
	query := `
		SELECT id, tenant_id, student_id, billing_config_id, invoice_number, 
		       amount, status, due_date, paid_amount, created_at, updated_at
		FROM invoices
		WHERE invoice_number = $1
	`
	var invoice entity.Invoice
	err := r.db.QueryRow(ctx, query, invoiceNumber).Scan(
		&invoice.ID, &invoice.TenantID, &invoice.StudentID, &invoice.BillingConfigID, &invoice.InvoiceNumber,
		&invoice.Amount, &invoice.Status, &invoice.DueDate, &invoice.PaidAmount, &invoice.CreatedAt, &invoice.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) List(ctx context.Context, tenantID uuid.UUID, studentID uuid.UUID, status entity.InvoiceStatus) ([]*entity.Invoice, error) {
	query := `
		SELECT id, tenant_id, student_id, billing_config_id, invoice_number, 
		       amount, status, due_date, paid_amount, created_at, updated_at
		FROM invoices
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argCount := 1

	if studentID != uuid.Nil {
		argCount++
		query += fmt.Sprintf(" AND student_id = $%d", argCount)
		args = append(args, studentID)
	}

	if status != "" {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []*entity.Invoice
	for rows.Next() {
		var invoice entity.Invoice
		if err := rows.Scan(
			&invoice.ID, &invoice.TenantID, &invoice.StudentID, &invoice.BillingConfigID, &invoice.InvoiceNumber,
			&invoice.Amount, &invoice.Status, &invoice.DueDate, &invoice.PaidAmount, &invoice.CreatedAt, &invoice.UpdatedAt,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, &invoice)
	}
	return invoices, nil
}

func (r *invoiceRepository) Exists(ctx context.Context, studentID, billingConfigID uuid.UUID, month, year int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM invoices 
			WHERE student_id = $1 
			AND billing_config_id = $2 
			AND EXTRACT(MONTH FROM created_at) = $3 
			AND EXTRACT(YEAR FROM created_at) = $4
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, studentID, billingConfigID, month, year).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *invoiceRepository) UpdateOverdueStatus(ctx context.Context) (int64, error) {
	query := `
		UPDATE invoices
		SET status = $1, updated_at = NOW()
		WHERE status = $2 AND due_date < NOW()
	`
	tag, err := r.db.Exec(ctx, query, entity.InvoiceStatusOverdue, entity.InvoiceStatusUnpaid)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
