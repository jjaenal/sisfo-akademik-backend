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

func TestStudentUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	u := usecase.NewStudentUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()

	validStudent := &entity.Student{
		ID:       id,
		TenantID: tenantID,
		Name:     "John Doe",
		Status:   "active",
	}

	t.Run("Create Success", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validStudent).Return(nil)
		err := u.Create(context.Background(), validStudent)
		assert.NoError(t, err)
	})

	t.Run("Create Validation Error", func(t *testing.T) {
		invalidStudent := &entity.Student{}
		err := u.Create(context.Background(), invalidStudent)
		assert.Error(t, err)
	})

	t.Run("Create Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validStudent).Return(assert.AnError)
		err := u.Create(context.Background(), validStudent)
		assert.Error(t, err)
	})

	t.Run("GetByID Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validStudent, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validStudent, res)
	})

	t.Run("GetByID Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)
		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("List Success", func(t *testing.T) {
		list := []entity.Student{*validStudent}
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

	t.Run("Update Success", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validStudent).Return(nil)
		err := u.Update(context.Background(), validStudent)
		assert.NoError(t, err)
	})

	t.Run("Update Validation Error", func(t *testing.T) {
		invalidStudent := &entity.Student{}
		err := u.Update(context.Background(), invalidStudent)
		assert.Error(t, err)
	})

	t.Run("Update Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validStudent).Return(assert.AnError)
		err := u.Update(context.Background(), validStudent)
		assert.Error(t, err)
	})

	t.Run("Delete Success", func(t *testing.T) {
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
