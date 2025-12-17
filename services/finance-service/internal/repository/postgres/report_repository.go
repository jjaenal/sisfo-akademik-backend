package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type reportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) repository.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetDailyRevenue(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]*entity.DailyRevenue, error) {
	query := `
		SELECT TO_CHAR(transaction_date, 'YYYY-MM-DD') as date, SUM(amount) as amount
		FROM payments
		WHERE tenant_id = $1 AND transaction_date >= $2 AND transaction_date <= $3 AND status = 'COMPLETED'
		GROUP BY date
		ORDER BY date ASC
	`
	rows, err := r.db.Query(ctx, query, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*entity.DailyRevenue
	for rows.Next() {
		var report entity.DailyRevenue
		if err := rows.Scan(&report.Date, &report.Amount); err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}

func (r *reportRepository) GetMonthlyRevenue(ctx context.Context, tenantID uuid.UUID, year int) ([]*entity.MonthlyRevenue, error) {
	query := `
		SELECT TO_CHAR(transaction_date, 'YYYY-MM') as month, SUM(amount) as amount
		FROM payments
		WHERE tenant_id = $1 AND EXTRACT(YEAR FROM transaction_date) = $2 AND status = 'COMPLETED'
		GROUP BY month
		ORDER BY month ASC
	`
	rows, err := r.db.Query(ctx, query, tenantID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*entity.MonthlyRevenue
	for rows.Next() {
		var report entity.MonthlyRevenue
		if err := rows.Scan(&report.Month, &report.Amount); err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}

func (r *reportRepository) GetOutstandingInvoices(ctx context.Context, tenantID uuid.UUID) ([]*entity.OutstandingReport, error) {
	query := `
		SELECT s.id, s.name, i.invoice_number, i.amount, i.paid_amount, (i.amount - i.paid_amount) as remaining, i.due_date, i.status
		FROM invoices i
		JOIN students s ON i.student_id = s.id
		WHERE i.tenant_id = $1 AND i.status IN ('UNPAID', 'OVERDUE', 'PARTIAL')
		ORDER BY i.due_date ASC
	`
	rows, err := r.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*entity.OutstandingReport
	for rows.Next() {
		var report entity.OutstandingReport
		if err := rows.Scan(
			&report.StudentID,
			&report.StudentName,
			&report.InvoiceNumber,
			&report.Amount,
			&report.PaidAmount,
			&report.Remaining,
			&report.DueDate,
			&report.Status,
		); err != nil {
			return nil, err
		}
		reports = append(reports, &report)
	}
	return reports, nil
}
