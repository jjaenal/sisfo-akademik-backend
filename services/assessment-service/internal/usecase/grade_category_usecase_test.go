package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGradeCategoryUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradeCategoryUseCase(mockRepo, timeout)

	t.Run("success", func(t *testing.T) {
		category := &entity.GradeCategory{
			Name:   "Tugas",
			Weight: 20,
		}

		mockRepo.EXPECT().Create(gomock.Any(), category).Return(nil)

		err := u.Create(context.Background(), category)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		category := &entity.GradeCategory{
			Name: "", // Invalid
		}

		err := u.Create(context.Background(), category)
		assert.Error(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		category := &entity.GradeCategory{
			Name:   "Tugas",
			Weight: 20,
		}

		mockRepo.EXPECT().Create(gomock.Any(), category).Return(errors.New("db error"))

		err := u.Create(context.Background(), category)
		assert.Error(t, err)
	})
}

func TestGradeCategoryUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradeCategoryUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := &entity.GradeCategory{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("not found"))

		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGradeCategoryUseCase_GetByTenantID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradeCategoryUseCase(mockRepo, timeout)

	tenantID := "tenant-1"

	t.Run("success", func(t *testing.T) {
		expected := []*entity.GradeCategory{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}
		mockRepo.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(expected, nil)

		res, err := u.GetByTenantID(context.Background(), tenantID)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestGradeCategoryUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradeCategoryUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		category := &entity.GradeCategory{
			ID:     id,
			Name:   "Updated",
			Weight: 30,
		}

		mockRepo.EXPECT().Update(gomock.Any(), category).Return(nil)

		err := u.Update(context.Background(), category)
		assert.NoError(t, err)
	})
}

func TestGradeCategoryUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradeCategoryUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})
}
