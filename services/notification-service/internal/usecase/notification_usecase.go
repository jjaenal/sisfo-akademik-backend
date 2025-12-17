package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/service"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
)

type NotificationUseCase interface {
	Send(ctx context.Context, req *SendNotificationRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	ListByRecipient(ctx context.Context, recipient string) ([]*entity.Notification, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.NotificationStatus, errorMessage string) error
}

type SendNotificationRequest struct {
	Channel      entity.NotificationChannel
	Recipient    string
	Subject      string
	Body         string
	TemplateName string
	Data         map[string]string
}

type notificationUseCase struct {
	notifRepo    repository.NotificationRepository
	templateRepo repository.NotificationTemplateRepository
	emailService service.EmailService
	waService    service.WhatsAppService
	rabbitClient *rabbit.Client
	timeout      time.Duration
}

func NewNotificationUseCase(
	notifRepo repository.NotificationRepository,
	templateRepo repository.NotificationTemplateRepository,
	emailService service.EmailService,
	waService service.WhatsAppService,
	rabbitClient *rabbit.Client,
	timeout time.Duration,
) NotificationUseCase {
	return &notificationUseCase{
		notifRepo:    notifRepo,
		templateRepo: templateRepo,
		emailService: emailService,
		waService:    waService,
		rabbitClient: rabbitClient,
		timeout:      timeout,
	}
}

func (u *notificationUseCase) Send(ctx context.Context, req *SendNotificationRequest) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	var template *entity.NotificationTemplate
	var err error

	// 1. Resolve Template if provided
	if req.TemplateName != "" {
		template, err = u.templateRepo.GetByName(ctx, req.TemplateName)
		if err != nil {
			return fmt.Errorf("failed to get template: %w", err)
		}
		if template == nil {
			return errors.New("template not found")
		}

		// Simple substitution logic
		req.Subject = template.SubjectTemplate
		req.Body = template.BodyTemplate
		req.Channel = template.Channel
		
		// Note: Variable replacement would happen here (e.g. using strings.Replace or text/template)
		// For now we assume the body is ready or the data map is used by the caller to pre-format
		// TODO: Implement template engine logic
	}

	// 2. Create Notification Record
	notification := &entity.Notification{
		ID:        uuid.New(),
		Channel:   req.Channel,
		Recipient: req.Recipient,
		Subject:   req.Subject,
		Body:      req.Body,
		Status:    entity.NotificationStatusPending,
		CreatedAt: time.Now(),
	}

	if template != nil {
		notification.TemplateID = &template.ID
	}

	if err := u.notifRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create notification record: %w", err)
	}

	// 3. Send Notification Async
	go u.processNotification(notification)

	return nil
}

func (u *notificationUseCase) processNotification(n *entity.Notification) {
	// Create a background context for DB operations
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var err error
	maxRetries := 3

	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			log.Printf("Retry %d/%d for notification %s", i, maxRetries, n.ID)
			time.Sleep(time.Duration(i*2) * time.Second)
		}

		if n.Channel == entity.NotificationChannelEmail {
			err = u.emailService.SendEmail([]string{n.Recipient}, n.Subject, n.Body)
		} else if n.Channel == entity.NotificationChannelWhatsApp {
			err = u.waService.SendWhatsApp(n.Recipient, n.Body)
		} else {
			// Mock other channels
			log.Printf("Simulating sending %s to %s", n.Channel, n.Recipient)
			time.Sleep(1 * time.Second)
			err = nil
		}

		if err == nil {
			break
		}
	}

	if err != nil {
		log.Printf("Failed to send notification %s after retries: %v", n.ID, err)
		n.Status = entity.NotificationStatusFailed
		n.ErrorMessage = err.Error()
	} else {
		n.Status = entity.NotificationStatusSent
		sentAt := time.Now()
		n.SentAt = &sentAt
	}

	// Update status in DB
	if updateErr := u.notifRepo.Update(ctx, n); updateErr != nil {
		log.Printf("Failed to update notification status %s: %v", n.ID, updateErr)
	}

	// Publish event
	if u.rabbitClient != nil {
		eventType := "notification.sent"
		if n.Status == entity.NotificationStatusFailed {
			eventType = "notification.failed"
		}
		
		payload := map[string]any{
			"notification_id": n.ID,
			"channel":         n.Channel,
			"recipient":       n.Recipient,
			"status":          n.Status,
			"error_message":   n.ErrorMessage,
			"timestamp":       time.Now(),
		}
		
		if err := u.rabbitClient.PublishJSON("events", "notification."+string(n.Channel), payload); err != nil {
			log.Printf("Failed to publish event %s: %v", eventType, err)
		}
	}
}

func (u *notificationUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.notifRepo.GetByID(ctx, id)
}

func (u *notificationUseCase) ListByRecipient(ctx context.Context, recipient string) ([]*entity.Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()
	return u.notifRepo.ListByRecipient(ctx, recipient)
}

func (u *notificationUseCase) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.NotificationStatus, errorMessage string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	notif, err := u.notifRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if notif == nil {
		return errors.New("notification not found")
	}

	notif.Status = status
	notif.ErrorMessage = errorMessage
	if status == entity.NotificationStatusSent {
		now := time.Now()
		notif.SentAt = &now
	}

	return u.notifRepo.Update(ctx, notif)
}
