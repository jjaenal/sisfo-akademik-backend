package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		url = "postgres://dev:dev@localhost:55432/devdb?sslmode=disable"
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Skip("no db available:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Skip("db connect failed:", err)
	}
	return db
}

func ensureMigrations(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	
	files := []string{
		"../../migrations/001_create_finance_tables.up.sql",
		"../../migrations/002_create_students_table.up.sql",
	}

	for _, file := range files {
		b, err := os.ReadFile(file)
		if err == nil {
			_, _ = db.Exec(ctx, string(b))
		}
	}
}

func setupServer(t *testing.T) (*gin.Engine, *pgxpool.Pool) {
	gin.SetMode(gin.TestMode)
	db := testDB(t)
	ensureMigrations(t, db)

	// Init Repositories
	billingRepo := postgres.NewBillingConfigRepository(db)
	invoiceRepo := postgres.NewInvoiceRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	studentRepo := postgres.NewStudentRepository(db)
	reportRepo := postgres.NewReportRepository(db)

	// Init UseCases
	timeout := 5 * time.Second
	billingUC := usecase.NewBillingConfigUseCase(billingRepo, timeout)
	invoiceUC := usecase.NewInvoiceUseCase(invoiceRepo, billingRepo, studentRepo, timeout)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, invoiceRepo, timeout)
	reportUC := usecase.NewReportUseCase(reportRepo, invoiceRepo, paymentRepo, timeout)

	// Init Handlers
	billingHandler := handler.NewBillingConfigHandler(billingUC)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUC)
	paymentHandler := handler.NewPaymentHandler(paymentUC)
	reportHandler := handler.NewReportHandler(reportUC)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	finance := v1.Group("/finance")
	{
		// Billing Configs
		finance.POST("/billing-configs", billingHandler.Create)
		finance.GET("/billing-configs", billingHandler.List)
		finance.GET("/billing-configs/:id", billingHandler.GetByID)
		
		// Invoices
		finance.POST("/invoices/generate", invoiceHandler.Generate)
		finance.GET("/invoices", invoiceHandler.List)
		finance.GET("/invoices/:id", invoiceHandler.GetByID)

		// Payments
		finance.POST("/payments", paymentHandler.Record)
		finance.GET("/payments/:id", paymentHandler.GetByID)

		// Reports
		finance.GET("/reports/revenue/daily", reportHandler.GetDailyRevenue)
		finance.GET("/reports/revenue/monthly", reportHandler.GetMonthlyRevenue)
		finance.GET("/reports/outstanding", reportHandler.GetOutstandingInvoices)
		finance.GET("/reports/student/:student_id/history", reportHandler.GetStudentHistory)
	}

	return r, db
}

func TestIntegration_Finance(t *testing.T) {
	r, db := setupServer(t)
	defer db.Close()

	tenantID := uuid.New()
	studentID := uuid.New()
	
	// Setup: Create a student
	ctx := context.Background()
	_, err := db.Exec(ctx, `
		INSERT INTO students (id, tenant_id, name, status, created_at, updated_at)
		VALUES ($1, $2, 'John Doe', 'ACTIVE', NOW(), NOW())
	`, studentID, tenantID)
	if err != nil {
		t.Fatalf("Failed to create student: %v", err)
	}

	var billingConfigID string
	var invoiceID string

	t.Run("Create Billing Config", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_id": tenantID.String(),
			"name":      "SPP Bulanan",
			"amount":    500000,
			"frequency": "MONTHLY",
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/finance/billing-configs", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		billingConfigID = resp.Data.ID
		assert.NotEmpty(t, billingConfigID)
	})

	t.Run("Generate Invoice", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_id":         tenantID.String(),
			"student_id":        studentID.String(),
			"billing_config_id": billingConfigID,
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/finance/invoices/generate", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		invoiceID = resp.Data.ID
		assert.NotEmpty(t, invoiceID)
	})

	t.Run("Record Payment", func(t *testing.T) {
		body := map[string]interface{}{
			"tenant_id":        tenantID.String(),
			"invoice_id":       invoiceID,
			"amount":           500000,
			"payment_method":   entity.PaymentMethodTransfer,
			"reference_number": "REF123456",
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/finance/payments", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "SUCCESS", resp.Data.Status)
	})

	t.Run("Verify Invoice Paid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/finance/invoices/%s", invoiceID), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp struct {
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "PAID", resp.Data.Status)
	})
}
