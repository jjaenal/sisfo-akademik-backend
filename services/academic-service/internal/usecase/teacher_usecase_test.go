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

func TestTeacherUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTeacherRepository(ctrl)
	u := usecase.NewTeacherUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()

	validTeacher := &entity.Teacher{
		ID:       id,
		TenantID: tenantID,
		Name:     "Budi Santoso",
		NIP:      "198001012005011001",
		Status:   "active",
	}

	t.Run("Create Success", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validTeacher).Return(nil)
		err := u.Create(context.Background(), validTeacher)
		assert.NoError(t, err)
	})

	t.Run("Create Validation Error", func(t *testing.T) {
		invalidTeacher := &entity.Teacher{}
		err := u.Create(context.Background(), invalidTeacher)
		assert.Error(t, err)
	})

	t.Run("Create Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validTeacher).Return(assert.AnError)
		err := u.Create(context.Background(), validTeacher)
		assert.Error(t, err)
	})

	t.Run("GetByID Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validTeacher, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validTeacher, res)
	})

	t.Run("GetByID Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)
		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("List Success", func(t *testing.T) {
		list := []entity.Teacher{*validTeacher}
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
		mockRepo.EXPECT().Update(gomock.Any(), validTeacher).Return(nil)
		err := u.Update(context.Background(), validTeacher)
		assert.NoError(t, err)
	})

	t.Run("Update Validation Error", func(t *testing.T) {
		invalidTeacher := &entity.Teacher{}
		err := u.Update(context.Background(), invalidTeacher)
		assert.Error(t, err)
	})

	t.Run("Update Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validTeacher).Return(assert.AnError)
		err := u.Update(context.Background(), validTeacher)
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
