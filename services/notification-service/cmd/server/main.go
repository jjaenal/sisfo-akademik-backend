package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/infrastructure/email"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/infrastructure/whatsapp"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/middleware"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/tracer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// @title Notification Service API
// @version 1.0
// @description Service for managing and sending notifications (Email, WhatsApp)
// @host localhost:9097
// @BasePath /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	logr, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	logr.Info("config loaded")

	// Tracer
	tp, err := tracer.InitTracer("notification-service", cfg.JaegerEndpoint)
	if err != nil {
		logr.Fatal("failed to init tracer", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logr.Fatal("failed to shutdown tracer", zap.Error(err))
		}
	}()

	// Connect to Database
	dbPool, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		logr.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer dbPool.Close()
	logr.Info("database connected")

	// Init Repositories
	templateRepo := postgres.NewNotificationTemplateRepository(dbPool)
	notifRepo := postgres.NewNotificationRepository(dbPool)

	// Init Services
	emailService := email.NewSMTPEmailService(cfg)
	waService := whatsapp.NewHTTPWhatsAppService(cfg)

	// Init RabbitMQ for Publishing
	rabbitClient := rabbit.New(cfg.RabbitURL)
	if rabbitClient == nil {
		logr.Warn("Failed to connect to RabbitMQ for publishing events")
	} else {
		defer rabbitClient.Close()
	}

	// Init UseCases
	timeout := 5 * time.Second
	templateUC := usecase.NewNotificationTemplateUseCase(templateRepo, timeout)
	notifUC := usecase.NewNotificationUseCase(notifRepo, templateRepo, emailService, waService, rabbitClient, timeout)

	// Init Handlers
	templateHandler := handler.NewNotificationTemplateHandler(templateUC)
	notifHandler := handler.NewNotificationHandler(notifUC)
	webhookHandler := handler.NewWebhookHandler(notifUC)

	// RabbitMQ Consumer (Background)
	go startRabbitConsumer(cfg.RabbitURL, notifUC, logr)

	// Init Gin
	r := gin.Default()
	r.Use(otelgin.Middleware("notification-service"))
	redis := redisutil.New(cfg.RedisAddr)
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimitByPolicy(limiter, 100, 30, map[string]int{}))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"service": "notification-service",
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
	notifications := v1.Group("/notifications")
	{
		// Templates
		notifications.POST("/templates", templateHandler.Create)
		notifications.GET("/templates", templateHandler.List)
		notifications.GET("/templates/:id", templateHandler.GetByID)
		notifications.PUT("/templates/:id", templateHandler.Update)
		notifications.DELETE("/templates/:id", templateHandler.Delete)

		// Notifications
		notifications.POST("/send", notifHandler.Send)
		notifications.GET("/:id", notifHandler.GetByID)
		notifications.GET("/recipient", notifHandler.ListByRecipient)
	}

	// Webhooks
	r.POST("/webhooks/:provider", webhookHandler.HandleWebhook)

	// Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9097"
	}
	logr.Info(fmt.Sprintf("Starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		logr.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func startRabbitConsumer(url string, notifUC usecase.NotificationUseCase, logr *zap.Logger) {
	if url == "" {
		url = "amqp://dev:dev@rabbitmq:5672/"
	}
	rb := rabbit.New(url)
	if rb == nil {
		logr.Error("Failed to connect to RabbitMQ")
		return
	}

	// 1. Consume Events
	msgsEvents, err := rb.Consume("events", "notification-service.events", []string{"auth.password_reset.requested"})
	if err != nil {
		logr.Error(fmt.Sprintf("rabbit consume events error: %v", err))
		return
	}

	// 2. Consume Notification Tasks
	msgsTasks, err := rb.Consume("notifications", "notification-service.tasks", []string{"notification.send"})
	if err != nil {
		logr.Error(fmt.Sprintf("rabbit consume tasks error: %v", err))
		return
	}

	logr.Info("RabbitMQ consumers started")

	// Handle Events
	go func() {
		for m := range msgsEvents {
			logr.Info(fmt.Sprintf("received event: rk=%s body=%s", m.RoutingKey, string(m.Body)))
			var ev struct {
				TenantID string `json:"tenant_id"`
				UserID   string `json:"user_id"`
				Email    string `json:"email"`
				Token    string `json:"token"`
				Type     string `json:"type"`
			}
			if err := json.Unmarshal(m.Body, &ev); err == nil && ev.Type == "password_reset" && ev.Email != "" && ev.Token != "" {
				// Use UseCase to send notification
				req := &usecase.SendNotificationRequest{
					Channel:   entity.NotificationChannelEmail,
					Recipient: ev.Email,
					Subject:   "Reset Password",
					Body:      fmt.Sprintf("Click here to reset your password: %s", ev.Token), // In real app, use template
				}
				if err := notifUC.Send(context.Background(), req); err != nil {
					log.Printf("failed to send email from event: %v", err)
				}
			}
		}
	}()

	// Handle Tasks
	for m := range msgsTasks {
		logr.Info(fmt.Sprintf("received task: rk=%s", m.RoutingKey))
		var task struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(m.Body, &task); err == nil && task.ID != "" {
			id, err := uuid.Parse(task.ID)
			if err == nil {
				if err := notifUC.Process(context.Background(), id); err != nil {
					logr.Error(fmt.Sprintf("failed to process notification %s: %v", id, err))
				}
			} else {
				logr.Error(fmt.Sprintf("invalid uuid in task: %s", task.ID))
			}
		}
	}
}
