package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBillingConfigUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBillingConfigRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewBillingConfigUseCase(mockRepo, timeout)

	t.Run("success", func(t *testing.T) {
		config := &entity.BillingConfig{
			Name:      "SPP Monthly",
			Amount:    500000,
			Frequency: entity.BillingFrequencyMonthly,
		}

		mockRepo.EXPECT().Create(gomock.Any(), config).Return(nil)

		err := u.Create(context.Background(), config)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		config := &entity.BillingConfig{
			Name: "", // Invalid
		}

		err := u.Create(context.Background(), config)
		assert.Error(t, err)
		assert.Equal(t, "invalid input data", err.Error())
	})

	t.Run("repository error", func(t *testing.T) {
		config := &entity.BillingConfig{
			Name:      "SPP Monthly",
			Amount:    500000,
			Frequency: entity.BillingFrequencyMonthly,
		}

		mockRepo.EXPECT().Create(gomock.Any(), config).Return(errors.New("db error"))

		err := u.Create(context.Background(), config)
		assert.Error(t, err)
	})
}

func TestBillingConfigUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBillingConfigRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewBillingConfigUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		config := &entity.BillingConfig{
			ID:   id,
			Name: "SPP Monthly",
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(config, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, config, res)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("not found"))

		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
