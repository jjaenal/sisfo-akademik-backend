package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/repository"
)

type InvoiceUseCase interface {
	Generate(ctx context.Context, tenantID, studentID, billingConfigID uuid.UUID) (*entity.Invoice, error)
	GenerateAllMonthlyInvoices(ctx context.Context) error
	CheckOverdueInvoices(ctx context.Context) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error)
	List(ctx context.Context, tenantID, studentID uuid.UUID, status entity.InvoiceStatus) ([]*entity.Invoice, error)
}

type invoiceUseCase struct {
	invoiceRepo       repository.InvoiceRepository
	billingConfigRepo repository.BillingConfigRepository
	studentRepo       repository.StudentRepository
	timeout           time.Duration
}

func NewInvoiceUseCase(
	invoiceRepo repository.InvoiceRepository,
	billingConfigRepo repository.BillingConfigRepository,
	studentRepo repository.StudentRepository,
	timeout time.Duration,
) InvoiceUseCase {
	return &invoiceUseCase{
		invoiceRepo:       invoiceRepo,
		billingConfigRepo: billingConfigRepo,
		studentRepo:       studentRepo,
		timeout:           timeout,
	}
}

func (u *invoiceUseCase) Generate(ctx context.Context, tenantID, studentID, billingConfigID uuid.UUID) (*entity.Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	// 1. Get Billing Config
	config, err := u.billingConfigRepo.GetByID(ctx, billingConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing config: %w", err)
	}

	// 2. Create Invoice
	now := time.Now()
	invoiceNumber := fmt.Sprintf("INV-%s-%d", now.Format("20060102"), rand.Intn(10000))
	
	// Default due date 7 days from now
	dueDate := now.AddDate(0, 0, 7)

	invoice := &entity.Invoice{
		ID:              uuid.New(),
		TenantID:        tenantID,
		StudentID:       studentID,
		BillingConfigID: billingConfigID,
		InvoiceNumber:   invoiceNumber,
		Amount:          config.Amount,
		Status:          entity.InvoiceStatusUnpaid,
		DueDate:         dueDate,
		PaidAmount:      0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (u *invoiceUseCase) GenerateAllMonthlyInvoices(ctx context.Context) error {
	// 1. Get all active monthly billing configs for all tenants
	configs, err := u.billingConfigRepo.ListAllActiveMonthly(ctx)
	if err != nil {
		return fmt.Errorf("failed to list all active monthly configs: %w", err)
	}

	if len(configs) == 0 {
		return nil
	}

	// Group configs by TenantID to minimize student queries
	configsByTenant := make(map[uuid.UUID][]*entity.BillingConfig)
	for _, config := range configs {
		configsByTenant[config.TenantID] = append(configsByTenant[config.TenantID], config)
	}

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	// 2. Process per tenant
	for tenantID, tenantConfigs := range configsByTenant {
		// Get all active students for this tenant
		students, err := u.studentRepo.GetActive(ctx, tenantID)
		if err != nil {
			fmt.Printf("failed to get active students for tenant %s: %v\n", tenantID, err)
			continue
		}

		if len(students) == 0 {
			continue
		}

		// 3. Generate invoices for each student and each config
		for _, student := range students {
			for _, config := range tenantConfigs {
				// Check if invoice exists
				exists, err := u.invoiceRepo.Exists(ctx, student.ID, config.ID, month, year)
				if err != nil {
					fmt.Printf("failed to check invoice existence for student %s config %s: %v\n", student.ID, config.ID, err)
					continue
				}

				if !exists {
					_, err := u.Generate(ctx, tenantID, student.ID, config.ID)
					if err != nil {
						fmt.Printf("failed to generate invoice for student %s config %s: %v\n", student.ID, config.ID, err)
					}
				}
			}
		}
	}

	return nil
}

func (u *invoiceUseCase) CheckOverdueInvoices(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	count, err := u.invoiceRepo.UpdateOverdueStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to update overdue status: %w", err)
	}

	if count > 0 {
		fmt.Printf("Updated %d invoices to OVERDUE status\n", count)
	}
	return nil
}

func (u *invoiceUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.invoiceRepo.GetByID(ctx, id)
}

func (u *invoiceUseCase) List(ctx context.Context, tenantID, studentID uuid.UUID, status entity.InvoiceStatus) ([]*entity.Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.invoiceRepo.List(ctx, tenantID, studentID, status)
}
