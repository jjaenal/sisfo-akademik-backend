package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/tracer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// @title           Finance Service API
// @version         1.0
// @description     Finance Service for Academic System
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9096
// @BasePath  /api/v1

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	log.Info("config loaded")

	// Tracer
	tp, err := tracer.InitTracer("finance-service", "http://jaeger:14268/api/traces")
	if err != nil {
		log.Fatal("failed to init tracer", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal("failed to shutdown tracer", zap.Error(err))
		}
	}()

	// Connect to Database
	dbPool, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer dbPool.Close()
	log.Info("database connected")

	// Init Repositories
	billingRepo := postgres.NewBillingConfigRepository(dbPool)
	invoiceRepo := postgres.NewInvoiceRepository(dbPool)
	paymentRepo := postgres.NewPaymentRepository(dbPool)
	studentRepo := postgres.NewStudentRepository(dbPool)
	reportRepo := postgres.NewReportRepository(dbPool)

	// Init UseCases
	timeout := 5 * time.Second
	billingUC := usecase.NewBillingConfigUseCase(billingRepo, timeout)
	invoiceUC := usecase.NewInvoiceUseCase(invoiceRepo, billingRepo, studentRepo, timeout)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, invoiceRepo, timeout)
	reportUC := usecase.NewReportUseCase(reportRepo, invoiceRepo, paymentRepo, timeout)

	// Start Scheduler
	go func() {
		// Run initial check after 10 seconds
		time.Sleep(10 * time.Second)
		log.Info("Starting initial invoice generation check...")
		if err := invoiceUC.GenerateAllMonthlyInvoices(context.Background()); err != nil {
			log.Error(fmt.Sprintf("Failed to generate invoices: %v", err))
		}

		log.Info("Starting initial overdue invoice check...")
		if err := invoiceUC.CheckOverdueInvoices(context.Background()); err != nil {
			log.Error(fmt.Sprintf("Failed to check overdue invoices: %v", err))
		}

		// Run every 24 hours
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Info("Starting scheduled invoice generation check...")
			if err := invoiceUC.GenerateAllMonthlyInvoices(context.Background()); err != nil {
				log.Error(fmt.Sprintf("Failed to generate invoices: %v", err))
			}

			log.Info("Starting scheduled overdue invoice check...")
			if err := invoiceUC.CheckOverdueInvoices(context.Background()); err != nil {
				log.Error(fmt.Sprintf("Failed to check overdue invoices: %v", err))
			}
		}
	}()

	// Init Handlers
	billingHandler := handler.NewBillingConfigHandler(billingUC)
	invoiceHandler := handler.NewInvoiceHandler(invoiceUC)
	paymentHandler := handler.NewPaymentHandler(paymentUC)
	reportHandler := handler.NewReportHandler(reportUC)

	// Init Gin
	r := gin.Default()
	r.Use(otelgin.Middleware("finance-service"))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"service": "finance-service",
				"status":  "healthy",
				"db":      "connected",
			},
			"meta": gin.H{
				"timestamp":  time.Now().UTC(),
				"request_id": uuid.NewString(),
			},
		})
	})

	// Routes
	v1 := r.Group("/api/v1")
	finance := v1.Group("/finance")
	{
		// Billing Configs
		finance.POST("/billing-configs", billingHandler.Create)
		finance.GET("/billing-configs", billingHandler.List)
		finance.GET("/billing-configs/:id", billingHandler.GetByID)
		finance.PUT("/billing-configs/:id", billingHandler.Update)
		finance.DELETE("/billing-configs/:id", billingHandler.Delete)

		// Invoices
		finance.POST("/invoices/generate", invoiceHandler.Generate)
		finance.GET("/invoices", invoiceHandler.List)
		finance.GET("/invoices/:id", invoiceHandler.GetByID)

		// Payments
		finance.POST("/payments", paymentHandler.Record)
		finance.GET("/payments/:id", paymentHandler.GetByID)
		finance.GET("/invoices/:id/payments", paymentHandler.ListByInvoice)

		// Reports
		finance.GET("/reports/revenue/daily", reportHandler.GetDailyRevenue)
		finance.GET("/reports/revenue/monthly", reportHandler.GetMonthlyRevenue)
		finance.GET("/reports/outstanding", reportHandler.GetOutstandingInvoices)
		finance.GET("/reports/student/:student_id/history", reportHandler.GetStudentHistory)
	}

	// Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9096"
	}

	log.Info(fmt.Sprintf("starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err.Error())
	}
}
