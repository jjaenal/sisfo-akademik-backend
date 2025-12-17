package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestApplicationUseCase_SubmitApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewApplicationUseCase(mockRepo, nil, timeout)

	t.Run("success", func(t *testing.T) {
		application := &entity.Application{
			AdmissionPeriodID: uuid.New(),
			FirstName:         "John",
			LastName:          "Doe",
			Email:             "john@example.com",
			PhoneNumber:       "1234567890",
			PreviousSchool:    "High School A",
			AverageScore:      85.0,
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, app *entity.Application) error {
			assert.NotEmpty(t, app.RegistrationNumber)
			assert.Equal(t, entity.ApplicationStatusSubmitted, app.Status)
			assert.NotNil(t, app.SubmissionDate)
			return nil
		})

		res, err := u.SubmitApplication(context.Background(), application)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("validation error", func(t *testing.T) {
		application := &entity.Application{
			FirstName: "", // Invalid
		}

		res, err := u.SubmitApplication(context.Background(), application)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "validation error", err.Error())
	})
}

func TestApplicationUseCase_GetApplicationStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewApplicationUseCase(mockRepo, nil, timeout)

	regNum := "REG-20250101-1234"

	t.Run("success", func(t *testing.T) {
		expectedApp := &entity.Application{RegistrationNumber: regNum}
		mockRepo.EXPECT().GetByRegistrationNumber(gomock.Any(), regNum).Return(expectedApp, nil)

		res, err := u.GetApplicationStatus(context.Background(), regNum)
		assert.NoError(t, err)
		assert.Equal(t, expectedApp, res)
	})
}

func TestApplicationUseCase_VerifyApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewApplicationUseCase(mockRepo, nil, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		existingApp := &entity.Application{ID: id, Status: entity.ApplicationStatusSubmitted}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existingApp, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, app *entity.Application) error {
			assert.Equal(t, entity.ApplicationStatusVerified, app.Status)
			return nil
		})

		err := u.VerifyApplication(context.Background(), id, entity.ApplicationStatusVerified)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		err := u.VerifyApplication(context.Background(), id, entity.ApplicationStatusVerified)
		assert.Error(t, err)
		assert.Equal(t, "application not found", err.Error())
	})
}
