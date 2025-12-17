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
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
	"go.uber.org/zap"
)

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

	msgs, err := rb.Consume("events", "notification-service.events", []string{"auth.password_reset.requested"})
	if err != nil {
		logr.Error(fmt.Sprintf("rabbit consume error: %v", err))
		return
	}

	logr.Info("RabbitMQ consumer started")
	for m := range msgs {
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
}
