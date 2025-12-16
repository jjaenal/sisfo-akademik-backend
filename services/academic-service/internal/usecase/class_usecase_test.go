package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClassUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClassRepository(ctrl)
	u := usecase.NewClassUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	academicYearID := uuid.New()
	id := uuid.New()

	validClass := &entity.Class{
		ID:             id,
		TenantID:       tenantID,
		AcademicYearID: &academicYearID,
		Name:           "Class 10A",
		Level:          10,
		Major:          "Science",
		Capacity:       30,
	}

	t.Run("Create", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validClass).Return(nil)
		err := u.Create(context.Background(), validClass)
		assert.NoError(t, err)
	})

	t.Run("Create Validation Error", func(t *testing.T) {
		invalidClass := &entity.Class{}
		err := u.Create(context.Background(), invalidClass)
		assert.Error(t, err)
	})

	t.Run("Create Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validClass).Return(assert.AnError)
		err := u.Create(context.Background(), validClass)
		assert.Error(t, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validClass, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validClass, res)
	})

	t.Run("GetByID Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)
		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("List", func(t *testing.T) {
		list := []entity.Class{*validClass}
		total := 1
		mockRepo.EXPECT().List(gomock.Any(), tenantID, 10, 0).Return(list, total, nil)
		res, count, err := u.List(context.Background(), tenantID, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
		assert.Equal(t, total, count)
	})

	t.Run("List Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().List(gomock.Any(), tenantID, 10, 0).Return(nil, 0, assert.AnError)
		res, count, err := u.List(context.Background(), tenantID, 10, 0)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, 0, count)
	})

	t.Run("Update", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validClass).Return(nil)
		err := u.Update(context.Background(), validClass)
		assert.NoError(t, err)
	})

	t.Run("Update Validation Error", func(t *testing.T) {
		invalidClass := &entity.Class{}
		err := u.Update(context.Background(), invalidClass)
		assert.Error(t, err)
	})

	t.Run("Update Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validClass).Return(assert.AnError)
		err := u.Update(context.Background(), validClass)
		assert.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("Delete Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(assert.AnError)
		err := u.Delete(context.Background(), id)
		assert.Error(t, err)
	})
}
