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

func TestSubjectUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSubjectRepository(ctrl)
	u := usecase.NewSubjectUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()

	validSubject := &entity.Subject{
		ID:          id,
		TenantID:    tenantID,
		Code:        "MATH101",
		Name:        "Mathematics",
		Description: "Intro to Math",
		CreditUnits: 3,
	}

	t.Run("Create", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validSubject).Return(nil)
		err := u.Create(context.Background(), validSubject)
		assert.NoError(t, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validSubject, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validSubject, res)
	})

	t.Run("List", func(t *testing.T) {
		list := []entity.Subject{*validSubject}
		total := 1
		mockRepo.EXPECT().List(gomock.Any(), tenantID, 10, 0).Return(list, total, nil)
		res, count, err := u.List(context.Background(), tenantID, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
		assert.Equal(t, total, count)
	})

	t.Run("Update", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validSubject).Return(nil)
		err := u.Update(context.Background(), validSubject)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})
}
