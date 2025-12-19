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

func TestApplicationUseCase_ListApplications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewApplicationUseCase(mockRepo, nil, timeout)

	ctx := context.Background()
	filter := map[string]interface{}{"status": "submitted"}

	t.Run("Success", func(t *testing.T) {
		expectedApps := []*entity.Application{
			{ID: uuid.New(), Status: entity.ApplicationStatusSubmitted},
		}
		mockRepo.EXPECT().List(gomock.Any(), filter).Return(expectedApps, nil)

		apps, err := u.ListApplications(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, expectedApps, apps)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepo.EXPECT().List(gomock.Any(), filter).Return(nil, errors.New("db error"))

		apps, err := u.ListApplications(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, apps)
	})
}

func TestApplicationUseCase_InputTestScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	u := usecase.NewApplicationUseCase(mockRepo, nil, 2*time.Second)

	ctx := context.Background()
	id := uuid.New()
	score := 85.5

	t.Run("Success", func(t *testing.T) {
		app := &entity.Application{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(app, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, a *entity.Application) error {
			assert.Equal(t, score, *a.TestScore)
			return nil
		})

		err := u.InputTestScore(ctx, id, score)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		err := u.InputTestScore(ctx, id, score)
		assert.Error(t, err)
		assert.Equal(t, "application not found", err.Error())
	})
}

func TestApplicationUseCase_InputInterviewScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	u := usecase.NewApplicationUseCase(mockRepo, nil, 2*time.Second)

	ctx := context.Background()
	id := uuid.New()
	score := 90.0

	t.Run("Success", func(t *testing.T) {
		app := &entity.Application{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(app, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, a *entity.Application) error {
			assert.Equal(t, score, *a.InterviewScore)
			return nil
		})

		err := u.InputInterviewScore(ctx, id, score)
		assert.NoError(t, err)
	})
}

func TestApplicationUseCase_CalculateFinalScores(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	u := usecase.NewApplicationUseCase(mockRepo, nil, 2*time.Second)

	ctx := context.Background()
	periodID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		testScore := 80.0
		interviewScore := 90.0
		avgScore := 85.0
		
		// Expected: (80*0.4) + (90*0.4) + (85*0.2) = 32 + 36 + 17 = 85

		app := &entity.Application{
			ID:             uuid.New(),
			TestScore:      &testScore,
			InterviewScore: &interviewScore,
			AverageScore:   avgScore,
		}

		mockRepo.EXPECT().List(gomock.Any(), map[string]interface{}{"admission_period_id": periodID}).Return([]*entity.Application{app}, nil)
		
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, a *entity.Application) error {
			assert.NotNil(t, a.FinalScore)
			assert.InDelta(t, 85.0, *a.FinalScore, 0.01)
			return nil
		})

		err := u.CalculateFinalScores(ctx, periodID)
		assert.NoError(t, err)
	})

	t.Run("SkipMissingScores", func(t *testing.T) {
		app := &entity.Application{ID: uuid.New()} // Missing scores
		mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*entity.Application{app}, nil)
		// Expect no Update call

		err := u.CalculateFinalScores(ctx, periodID)
		assert.NoError(t, err)
	})
}

func TestApplicationUseCase_RegisterStudent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockApplicationRepository(ctrl)
	u := usecase.NewApplicationUseCase(mockRepo, nil, 2*time.Second)

	ctx := context.Background()
	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		app := &entity.Application{
			ID:     id,
			Status: entity.ApplicationStatusAccepted,
		}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(app, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, a *entity.Application) error {
			assert.Equal(t, entity.ApplicationStatusRegistered, a.Status)
			return nil
		})

		err := u.RegisterStudent(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("NotAccepted", func(t *testing.T) {
		app := &entity.Application{
			ID:     id,
			Status: entity.ApplicationStatusSubmitted,
		}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(app, nil)

		err := u.RegisterStudent(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, "application must be accepted to register", err.Error())
	})

	t.Run("AlreadyRegistered", func(t *testing.T) {
		app := &entity.Application{
			ID:     id,
			Status: entity.ApplicationStatusRegistered,
		}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(app, nil)

		err := u.RegisterStudent(ctx, id)
		assert.Error(t, err)
		assert.Equal(t, "student already registered", err.Error())
	})
}
