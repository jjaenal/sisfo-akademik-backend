package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

// Mock Services
type MockEmailService struct {
	SendEmailFunc func(to []string, subject, body string) error
}

func (m *MockEmailService) SendEmail(to []string, subject, body string) error {
	if m.SendEmailFunc != nil {
		return m.SendEmailFunc(to, subject, body)
	}
	return nil
}

type MockWhatsAppService struct {
	SendWhatsAppFunc func(to, message string) error
}

func (m *MockWhatsAppService) SendWhatsApp(to, message string) error {
	if m.SendWhatsAppFunc != nil {
		return m.SendWhatsAppFunc(to, message)
	}
	return nil
}

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
		"../../migrations/001_create_notification_tables.up.sql",
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
	templateRepo := postgres.NewNotificationTemplateRepository(db)
	notifRepo := postgres.NewNotificationRepository(db)

	// Init Mocks
	emailService := &MockEmailService{}
	waService := &MockWhatsAppService{}

	// Init UseCases
	timeout := 5 * time.Second
	templateUC := usecase.NewNotificationTemplateUseCase(templateRepo, timeout)
	notifUC := usecase.NewNotificationUseCase(notifRepo, templateRepo, emailService, waService, nil, timeout)

	// Init Handlers
	templateHandler := handler.NewNotificationTemplateHandler(templateUC)
	notifHandler := handler.NewNotificationHandler(notifUC)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	v1 := r.Group("/api/v1")
	notifications := v1.Group("/notifications")
	{
		// Templates
		notifications.POST("/templates", templateHandler.Create)
		notifications.GET("/templates", templateHandler.List)
		notifications.GET("/templates/:id", templateHandler.GetByID)
		
		// Notifications
		notifications.POST("/send", notifHandler.Send)
		notifications.GET("/:id", notifHandler.GetByID)
		notifications.GET("/recipient", notifHandler.ListByRecipient)
	}

	return r, db
}

func TestIntegration_Notification(t *testing.T) {
	r, db := setupServer(t)
	defer db.Close()

	t.Run("Create Template", func(t *testing.T) {
		body := map[string]interface{}{
			"name":             "welcome_email",
			"channel":          "EMAIL",
			"subject_template": "Welcome, {{name}}!",
			"body_template":    "Hello {{name}}, welcome to our platform.",
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/notifications/templates", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Send Email Notification", func(t *testing.T) {
		body := map[string]interface{}{
			"channel":   "EMAIL",
			"recipient": "user@example.com",
			"subject":   "Test Email",
			"body":      "This is a test email.",
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/notifications/send", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusAccepted {
			t.Logf("Response Body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusAccepted, w.Code) 

		// Verify notification status via GET by Recipient
		time.Sleep(100 * time.Millisecond) // Allow async process to run
		
		wGet := httptest.NewRecorder()
		reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/notifications/recipient?recipient=user@example.com", nil)
		r.ServeHTTP(wGet, reqGet)
		
		assert.Equal(t, http.StatusOK, wGet.Code)
		
		var respGet struct {
			Data []entity.Notification `json:"data"`
		}
		err := json.Unmarshal(wGet.Body.Bytes(), &respGet)
		assert.NoError(t, err)
		assert.NotEmpty(t, respGet.Data)
		assert.Equal(t, "user@example.com", respGet.Data[0].Recipient)
		// Since we use mock service which does nothing but return nil, status might be SENT if logic sets it to SENT after calling service
		// Check usecase logic. Usually it sets to SENT if service returns nil.
		assert.Equal(t, entity.NotificationStatusSent, respGet.Data[0].Status)
	})

	t.Run("Send WhatsApp Notification", func(t *testing.T) {
		body := map[string]interface{}{
			"channel":   "WHATSAPP",
			"recipient": "08123456789",
			"body":      "This is a test WA message.",
		}
		b, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/notifications/send", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusAccepted, w.Code)
	})
}
