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

func TestScheduleUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockScheduleRepository(ctrl)
	u := usecase.NewScheduleUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	classID := uuid.New()
	subjectID := uuid.New()
	teacherID := uuid.New()
	id := uuid.New()

	validSchedule := &entity.Schedule{
		ID:        id,
		TenantID:  tenantID,
		ClassID:   classID,
		SubjectID: subjectID,
		TeacherID: teacherID,
		DayOfWeek: 1,
		StartTime: "08:00",
		EndTime:   "09:30",
		Room:      "Room 101",
	}

	t.Run("Create Success", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return([]entity.Schedule{}, nil)
		mockRepo.EXPECT().Create(gomock.Any(), validSchedule).Return(nil)
		err := u.Create(context.Background(), validSchedule)
		assert.NoError(t, err)
	})

	t.Run("Create Validation Error", func(t *testing.T) {
		invalidSchedule := &entity.Schedule{}
		err := u.Create(context.Background(), invalidSchedule)
		assert.Error(t, err)
	})

	t.Run("Create CheckConflicts Error", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return(nil, assert.AnError)
		err := u.Create(context.Background(), validSchedule)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("Create Conflict", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return([]entity.Schedule{{ID: uuid.New()}}, nil)
		err := u.Create(context.Background(), validSchedule)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflict detected")
	})

	t.Run("Create Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return([]entity.Schedule{}, nil)
		mockRepo.EXPECT().Create(gomock.Any(), validSchedule).Return(assert.AnError)
		err := u.Create(context.Background(), validSchedule)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validSchedule, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validSchedule, res)
	})

	t.Run("List", func(t *testing.T) {
		list := []entity.Schedule{*validSchedule}
		total := 1
		mockRepo.EXPECT().List(gomock.Any(), tenantID, 10, 0).Return(list, total, nil)
		res, count, err := u.List(context.Background(), tenantID, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
		assert.Equal(t, total, count)
	})

	t.Run("Update Success", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return([]entity.Schedule{}, nil)
		mockRepo.EXPECT().Update(gomock.Any(), validSchedule).Return(nil)
		err := u.Update(context.Background(), validSchedule)
		assert.NoError(t, err)
	})

	t.Run("Update Validation Error", func(t *testing.T) {
		invalidSchedule := &entity.Schedule{}
		err := u.Update(context.Background(), invalidSchedule)
		assert.Error(t, err)
	})

	t.Run("Update CheckConflicts Error", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return(nil, assert.AnError)
		err := u.Update(context.Background(), validSchedule)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("Update Conflict", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return([]entity.Schedule{{ID: uuid.New()}}, nil)
		err := u.Update(context.Background(), validSchedule)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflict detected")
	})

	t.Run("Update Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().CheckConflicts(gomock.Any(), validSchedule).Return([]entity.Schedule{}, nil)
		mockRepo.EXPECT().Update(gomock.Any(), validSchedule).Return(assert.AnError)
		err := u.Update(context.Background(), validSchedule)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("ListByClass", func(t *testing.T) {
		list := []entity.Schedule{*validSchedule}
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return(list, nil)
		res, err := u.ListByClass(context.Background(), classID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})
}
