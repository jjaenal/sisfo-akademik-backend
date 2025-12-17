package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAdmissionPeriodUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdmissionPeriodRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewAdmissionPeriodUseCase(mockRepo, timeout)

	now := time.Now()
	nextWeek := now.Add(7 * 24 * time.Hour)

	t.Run("success", func(t *testing.T) {
		period := &entity.AdmissionPeriod{
			Name:      "PPDB 2025",
			StartDate: now,
			EndDate:   nextWeek,
			IsActive:  true,
		}

		mockRepo.EXPECT().Create(gomock.Any(), period).Return(nil)

		err := u.Create(context.Background(), period)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		period := &entity.AdmissionPeriod{
			Name: "", // Invalid
		}

		err := u.Create(context.Background(), period)
		assert.Error(t, err)
	})
}

func TestAdmissionPeriodUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdmissionPeriodRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewAdmissionPeriodUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := &entity.AdmissionPeriod{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestAdmissionPeriodUseCase_GetActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdmissionPeriodRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewAdmissionPeriodUseCase(mockRepo, timeout)

	t.Run("success", func(t *testing.T) {
		expected := &entity.AdmissionPeriod{IsActive: true}
		mockRepo.EXPECT().GetActive(gomock.Any()).Return(expected, nil)

		res, err := u.GetActive(context.Background())
		assert.NoError(t, err)
		assert.True(t, res.IsActive)
	})
}

func TestAdmissionPeriodUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdmissionPeriodRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewAdmissionPeriodUseCase(mockRepo, timeout)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.AdmissionPeriod{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}
		mockRepo.EXPECT().List(gomock.Any()).Return(expected, nil)

		res, err := u.List(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, len(expected), len(res))
	})
}

func TestAdmissionPeriodUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdmissionPeriodRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewAdmissionPeriodUseCase(mockRepo, timeout)

	id := uuid.New()
	now := time.Now()
	nextWeek := now.Add(7 * 24 * time.Hour)

	t.Run("success", func(t *testing.T) {
		period := &entity.AdmissionPeriod{
			ID:        id,
			Name:      "Updated",
			StartDate: now,
			EndDate:   nextWeek,
		}

		mockRepo.EXPECT().Update(gomock.Any(), period).Return(nil)

		err := u.Update(context.Background(), period)
		assert.NoError(t, err)
	})
}

func TestAdmissionPeriodUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAdmissionPeriodRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewAdmissionPeriodUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("db error"))

		err := u.Delete(context.Background(), id)
		assert.Error(t, err)
	})
}
