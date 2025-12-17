package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationUseCase_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockTemplateRepo := mocks.NewMockNotificationTemplateRepository(ctrl)
	mockEmailService := mocks.NewMockEmailService(ctrl)
	mockWAService := mocks.NewMockWhatsAppService(ctrl)
	
	// RabbitMQ client is nil for this test, assuming it handles nil gracefully (it seems so in the code)
	// or we can just ignore it as we are testing core logic.
	// The usecase checks `if u.rabbitClient != nil`.

	timeout := 2 * time.Second
	u := usecase.NewNotificationUseCase(
		mockNotifRepo,
		mockTemplateRepo,
		mockEmailService,
		mockWAService,
		nil, // rabbitClient
		timeout,
	)

	t.Run("send email success without template", func(t *testing.T) {
		req := &usecase.SendNotificationRequest{
			Channel:   entity.NotificationChannelEmail,
			Recipient: "test@example.com",
			Subject:   "Test Subject",
			Body:      "Test Body",
		}

		mockNotifRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, n *entity.Notification) error {
			assert.Equal(t, req.Recipient, n.Recipient)
			assert.Equal(t, req.Subject, n.Subject)
			assert.Equal(t, req.Body, n.Body)
			assert.Equal(t, entity.NotificationStatusPending, n.Status)
			return nil
		})

		// Expect email service call (async)
		mockEmailService.EXPECT().SendEmail([]string{req.Recipient}, req.Subject, req.Body).Return(nil)

		// Expect update status to Sent (async)
		mockNotifRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, n *entity.Notification) error {
			assert.Equal(t, entity.NotificationStatusSent, n.Status)
			return nil
		})

		err := u.Send(context.Background(), req)
		assert.NoError(t, err)

		// Wait for async processing
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("send whatsapp success without template", func(t *testing.T) {
		req := &usecase.SendNotificationRequest{
			Channel:   entity.NotificationChannelWhatsApp,
			Recipient: "+628123456789",
			Body:      "Test WA Body",
		}

		mockNotifRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		// Expect WA service call (async)
		mockWAService.EXPECT().SendWhatsApp(req.Recipient, req.Body).Return(nil)

		// Expect update status to Sent (async)
		mockNotifRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		err := u.Send(context.Background(), req)
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)
	})

	t.Run("send with template success", func(t *testing.T) {
		templateName := "welcome_email"
		template := &entity.NotificationTemplate{
			ID:              uuid.New(),
			Name:            templateName,
			Channel:         entity.NotificationChannelEmail,
			SubjectTemplate: "Welcome {{name}}",
			BodyTemplate:    "Hi {{name}}, welcome!",
		}

		req := &usecase.SendNotificationRequest{
			TemplateName: templateName,
			Recipient:    "user@example.com",
			Data: map[string]string{
				"name": "John",
			},
		}

		mockTemplateRepo.EXPECT().GetByName(gomock.Any(), templateName).Return(template, nil)

		// Note: The current implementation doesn't replace variables yet (TODO in code).
		// So we expect the template raw strings.

		mockNotifRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, n *entity.Notification) error {
			assert.Equal(t, template.SubjectTemplate, n.Subject)
			assert.Equal(t, template.BodyTemplate, n.Body)
			assert.Equal(t, template.ID, *n.TemplateID)
			return nil
		})

		mockEmailService.EXPECT().SendEmail([]string{req.Recipient}, template.SubjectTemplate, template.BodyTemplate).Return(nil)
		mockNotifRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		err := u.Send(context.Background(), req)
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)
	})
}
