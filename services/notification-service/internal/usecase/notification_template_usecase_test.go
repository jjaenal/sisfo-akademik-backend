package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationTemplateUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotificationTemplateRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationTemplateUseCase(mockRepo, timeout)

	t.Run("success", func(t *testing.T) {
		template := &entity.NotificationTemplate{
			Name:         "Welcome Email",
			Channel:      entity.NotificationChannelEmail,
			BodyTemplate: "Hello {{name}}",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, tmpl *entity.NotificationTemplate) error {
			assert.NotEqual(t, uuid.Nil, tmpl.ID)
			assert.True(t, tmpl.IsActive)
			assert.False(t, tmpl.CreatedAt.IsZero())
			assert.False(t, tmpl.UpdatedAt.IsZero())
			return nil
		})

		err := u.Create(context.Background(), template)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		template := &entity.NotificationTemplate{
			Name: "", // Empty name
		}

		err := u.Create(context.Background(), template)
		assert.Error(t, err)
		assert.Equal(t, "invalid input data", err.Error())
	})

	t.Run("repository error", func(t *testing.T) {
		template := &entity.NotificationTemplate{
			Name:         "Welcome Email",
			Channel:      entity.NotificationChannelEmail,
			BodyTemplate: "Hello {{name}}",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		err := u.Create(context.Background(), template)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestNotificationTemplateUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotificationTemplateRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationTemplateUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		template := &entity.NotificationTemplate{
			ID:           id,
			Name:         "Updated Name",
			Channel:      entity.NotificationChannelEmail,
			BodyTemplate: "Updated Body",
		}

		existing := &entity.NotificationTemplate{
			ID:   id,
			Name: "Old Name",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
		mockRepo.EXPECT().Update(gomock.Any(), template).Return(nil)

		err := u.Update(context.Background(), template)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		template := &entity.NotificationTemplate{
			ID:   id,
			Name: "", // Invalid
		}

		err := u.Update(context.Background(), template)
		assert.Error(t, err)
		assert.Equal(t, "invalid input data", err.Error())
	})

	t.Run("not found", func(t *testing.T) {
		template := &entity.NotificationTemplate{
			ID:           id,
			Name:         "Updated Name",
			Channel:      entity.NotificationChannelEmail,
			BodyTemplate: "Updated Body",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		err := u.Update(context.Background(), template)
		assert.Error(t, err)
		assert.Equal(t, "template not found", err.Error())
	})
}

func TestNotificationTemplateUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotificationTemplateRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationTemplateUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := &entity.NotificationTemplate{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		result, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("db error"))

		result, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestNotificationTemplateUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotificationTemplateRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationTemplateUseCase(mockRepo, timeout)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.NotificationTemplate{{Name: "T1"}, {Name: "T2"}}
		mockRepo.EXPECT().List(gomock.Any()).Return(expected, nil)

		result, err := u.List(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestNotificationTemplateUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotificationTemplateRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewNotificationTemplateUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})
}
