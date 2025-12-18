package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStudentUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	u := NewStudentUseCase(mockRepo, 2*time.Second)

	studentID := uuid.New()
	student := &entity.Student{
		ID:       studentID,
		Name:     "John Doe",
		TenantID: "tenant-1",
		Status:   "active",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), student).Return(nil)

		err := u.Create(context.Background(), student)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		invalidStudent := &entity.Student{} // Missing Name and TenantID
		err := u.Create(context.Background(), invalidStudent)
		assert.Error(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), student).Return(errors.New("db error"))

		err := u.Create(context.Background(), student)
		assert.Error(t, err)
	})
}

func TestStudentUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	u := NewStudentUseCase(mockRepo, 2*time.Second)

	id := uuid.New()
	student := &entity.Student{ID: id, Name: "John Doe"}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(student, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, student, res)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("not found"))

		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestStudentUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	u := NewStudentUseCase(mockRepo, 2*time.Second)

	students := []entity.Student{{Name: "John"}}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().List(gomock.Any(), "tenant-1", 10, 0).Return(students, 1, nil)

		res, count, err := u.List(context.Background(), "tenant-1", 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Equal(t, students, res)
	})
}

func TestStudentUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	u := NewStudentUseCase(mockRepo, 2*time.Second)

	student := &entity.Student{
		ID:       uuid.New(),
		Name:     "John Doe",
		TenantID: "tenant-1",
		Status:   "active",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), student).Return(nil)

		err := u.Update(context.Background(), student)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		invalidStudent := &entity.Student{}
		err := u.Update(context.Background(), invalidStudent)
		assert.Error(t, err)
	})
}

func TestStudentUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentRepository(ctrl)
	u := NewStudentUseCase(mockRepo, 2*time.Second)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)

		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})
}
