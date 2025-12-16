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

func TestSchoolUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSchoolRepository(ctrl)
	u := usecase.NewSchoolUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()

	validSchool := &entity.School{
		ID:       id,
		TenantID: tenantID,
		Name:     "SMA Negeri 1 Jakarta",
		Address:  "Jl. Budi Utomo No. 7",
	}

	t.Run("Create Success", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validSchool).Return(nil)
		err := u.Create(context.Background(), validSchool)
		assert.NoError(t, err)
	})

	t.Run("Create Validation Error", func(t *testing.T) {
		invalidSchool := &entity.School{}
		err := u.Create(context.Background(), invalidSchool)
		assert.Error(t, err)
	})

	t.Run("Create Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validSchool).Return(assert.AnError)
		err := u.Create(context.Background(), validSchool)
		assert.Error(t, err)
	})

	t.Run("GetByID Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validSchool, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validSchool, res)
	})

	t.Run("GetByID Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)
		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("GetByTenantID Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(validSchool, nil)
		res, err := u.GetByTenantID(context.Background(), tenantID)
		assert.NoError(t, err)
		assert.Equal(t, validSchool, res)
	})

	t.Run("GetByTenantID Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(nil, assert.AnError)
		res, err := u.GetByTenantID(context.Background(), tenantID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Update Success", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validSchool).Return(nil)
		err := u.Update(context.Background(), validSchool)
		assert.NoError(t, err)
	})

	t.Run("Update Validation Error", func(t *testing.T) {
		invalidSchool := &entity.School{}
		err := u.Update(context.Background(), invalidSchool)
		assert.Error(t, err)
	})

	t.Run("Update Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validSchool).Return(assert.AnError)
		err := u.Update(context.Background(), validSchool)
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
