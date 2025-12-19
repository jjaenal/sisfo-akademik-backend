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

		expectedSubject := "Welcome John"
		expectedBody := "Hi John, welcome!"

		mockNotifRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, n *entity.Notification) error {
			assert.Equal(t, expectedSubject, n.Subject)
			assert.Equal(t, expectedBody, n.Body)
			assert.Equal(t, template.ID, *n.TemplateID)
			return nil
		})

		mockEmailService.EXPECT().SendEmail([]string{req.Recipient}, expectedSubject, expectedBody).Return(nil)
		mockNotifRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		err := u.Send(context.Background(), req)
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)
	})
}

func TestNotificationUseCase_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	mockTemplateRepo := mocks.NewMockNotificationTemplateRepository(ctrl)
	mockEmailService := mocks.NewMockEmailService(ctrl)
	mockWAService := mocks.NewMockWhatsAppService(ctrl)

	timeout := 2 * time.Second
	u := usecase.NewNotificationUseCase(
		mockNotifRepo,
		mockTemplateRepo,
		mockEmailService,
		mockWAService,
		nil,
		timeout,
	)

	t.Run("process email notification success", func(t *testing.T) {
		id := uuid.New()
		notification := &entity.Notification{
			ID:        id,
			Channel:   entity.NotificationChannelEmail,
			Recipient: "test@example.com",
			Subject:   "Subject",
			Body:      "Body",
			Status:    entity.NotificationStatusPending,
		}

		mockNotifRepo.EXPECT().GetByID(gomock.Any(), id).Return(notification, nil)

		mockEmailService.EXPECT().SendEmail([]string{notification.Recipient}, notification.Subject, notification.Body).Return(nil)

		mockNotifRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, n *entity.Notification) error {
			assert.Equal(t, entity.NotificationStatusSent, n.Status)
			return nil
		})

		err := u.Process(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("process notification not found", func(t *testing.T) {
		id := uuid.New()
		mockNotifRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		err := u.Process(context.Background(), id)
		assert.Error(t, err)
		assert.Equal(t, "notification not found", err.Error())
	})
}

func TestNotificationUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationUseCase(mockNotifRepo, nil, nil, nil, nil, timeout)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		expected := &entity.Notification{ID: id}
		mockNotifRepo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error", func(t *testing.T) {
		id := uuid.New()
		mockNotifRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)

		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestNotificationUseCase_ListByRecipient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationUseCase(mockNotifRepo, nil, nil, nil, nil, timeout)

	t.Run("success", func(t *testing.T) {
		recipient := "user@example.com"
		expected := []*entity.Notification{{Recipient: recipient}}
		mockNotifRepo.EXPECT().ListByRecipient(gomock.Any(), recipient).Return(expected, nil)

		res, err := u.ListByRecipient(context.Background(), recipient)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error", func(t *testing.T) {
		recipient := "user@example.com"
		mockNotifRepo.EXPECT().ListByRecipient(gomock.Any(), recipient).Return(nil, assert.AnError)

		res, err := u.ListByRecipient(context.Background(), recipient)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestNotificationUseCase_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNotifRepo := mocks.NewMockNotificationRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationUseCase(mockNotifRepo, nil, nil, nil, nil, timeout)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		notif := &entity.Notification{ID: id, Status: entity.NotificationStatusPending}
		
		mockNotifRepo.EXPECT().GetByID(gomock.Any(), id).Return(notif, nil)
		mockNotifRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, n *entity.Notification) error {
			assert.Equal(t, entity.NotificationStatusSent, n.Status)
			assert.NotNil(t, n.SentAt)
			return nil
		})

		err := u.UpdateStatus(context.Background(), id, entity.NotificationStatusSent, "")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New()
		mockNotifRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		err := u.UpdateStatus(context.Background(), id, entity.NotificationStatusSent, "")
		assert.Error(t, err)
		assert.Equal(t, "notification not found", err.Error())
	})

	t.Run("get error", func(t *testing.T) {
		id := uuid.New()
		mockNotifRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)

		err := u.UpdateStatus(context.Background(), id, entity.NotificationStatusSent, "")
		assert.Error(t, err)
	})
}
