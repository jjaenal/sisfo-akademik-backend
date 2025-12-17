package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
)

type ReportRepository interface {
	GetDailyRevenue(ctx context.Context, tenantID uuid.UUID, startDate, endDate time.Time) ([]*entity.DailyRevenue, error)
	GetMonthlyRevenue(ctx context.Context, tenantID uuid.UUID, year int) ([]*entity.MonthlyRevenue, error)
	GetOutstandingInvoices(ctx context.Context, tenantID uuid.UUID) ([]*entity.OutstandingReport, error)
}
