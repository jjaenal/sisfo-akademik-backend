package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAcademicYearUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAcademicYearRepository(ctrl)
	u := usecase.NewAcademicYearUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()

	validYear := &entity.AcademicYear{
		ID:        id,
		TenantID:  tenantID,
		Name:      "2024/2025",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(1, 0, 0),
		IsActive:  true,
	}

	t.Run("Create", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validYear).Return(nil)
		err := u.Create(context.Background(), validYear)
		assert.NoError(t, err)
	})

	t.Run("Create Validation Error", func(t *testing.T) {
		invalidYear := &entity.AcademicYear{} // Missing required fields
		err := u.Create(context.Background(), invalidYear)
		assert.Error(t, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validYear, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validYear, res)
	})

	t.Run("List", func(t *testing.T) {
		list := []entity.AcademicYear{*validYear}
		mockRepo.EXPECT().List(gomock.Any(), tenantID).Return(list, nil)
		res, err := u.List(context.Background(), tenantID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})

	t.Run("Update", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validYear).Return(nil)
		err := u.Update(context.Background(), validYear)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("Delete Error", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("db error"))
		err := u.Delete(context.Background(), id)
		assert.Error(t, err)
	})
}
