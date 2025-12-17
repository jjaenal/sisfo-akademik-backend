package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type ReportUseCase interface {
	GetDailyRevenue(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]*entity.DailyRevenue, error)
	GetMonthlyRevenue(ctx context.Context, tenantID uuid.UUID, year int) ([]*entity.MonthlyRevenue, error)
	GetOutstandingInvoices(ctx context.Context, tenantID uuid.UUID) ([]*entity.OutstandingReport, error)
	GetStudentHistory(ctx context.Context, tenantID, studentID uuid.UUID) (*entity.StudentHistory, error)
}

type reportUseCase struct {
	reportRepo  repository.ReportRepository
	invoiceRepo repository.InvoiceRepository
	paymentRepo repository.PaymentRepository
	timeout     time.Duration
}

func NewReportUseCase(
	reportRepo repository.ReportRepository,
	invoiceRepo repository.InvoiceRepository,
	paymentRepo repository.PaymentRepository,
	timeout time.Duration,
) ReportUseCase {
	return &reportUseCase{
		reportRepo:  reportRepo,
		invoiceRepo: invoiceRepo,
		paymentRepo: paymentRepo,
		timeout:     timeout,
	}
}

func (u *reportUseCase) GetDailyRevenue(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]*entity.DailyRevenue, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.reportRepo.GetDailyRevenue(ctx, tenantID, startDate, endDate)
}

func (u *reportUseCase) GetMonthlyRevenue(ctx context.Context, tenantID uuid.UUID, year int) ([]*entity.MonthlyRevenue, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.reportRepo.GetMonthlyRevenue(ctx, tenantID, year)
}

func (u *reportUseCase) GetOutstandingInvoices(ctx context.Context, tenantID uuid.UUID) ([]*entity.OutstandingReport, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.reportRepo.GetOutstandingInvoices(ctx, tenantID)
}

func (u *reportUseCase) GetStudentHistory(ctx context.Context, tenantID, studentID uuid.UUID) (*entity.StudentHistory, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	invoices, err := u.invoiceRepo.List(ctx, tenantID, studentID, "")
	if err != nil {
		return nil, err
	}

	var allPayments []*entity.Payment
	for _, inv := range invoices {
		payments, err := u.paymentRepo.ListByInvoiceID(ctx, inv.ID)
		if err != nil {
			return nil, err
		}
		allPayments = append(allPayments, payments...)
	}

	return &entity.StudentHistory{
		Invoices: invoices,
		Payments: allPayments,
	}, nil
}
