package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTeacherAttendanceUseCase_CheckIn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTeacherAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewTeacherAttendanceUseCase(mockRepo, timeout)

	teacherID := uuid.New()
	semesterID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		attendance := &entity.TeacherAttendance{
			TeacherID:      teacherID,
			SemesterID:     semesterID,
			AttendanceDate: now,
			CheckInTime:    &now,
			Status:         entity.TeacherAttendanceStatusPresent,
			Notes:          "Present",
		}

		mockRepo.EXPECT().GetByTeacherAndDate(gomock.Any(), teacherID, now).Return(nil, errors.New("not found"))
		mockRepo.EXPECT().Create(gomock.Any(), attendance).Return(nil)

		err := u.CheckIn(context.Background(), attendance)
		assert.NoError(t, err)
	})

	t.Run("already checked in", func(t *testing.T) {
		attendance := &entity.TeacherAttendance{
			TeacherID:      teacherID,
			SemesterID:     semesterID,
			AttendanceDate: now,
			CheckInTime:    &now,
			Status:         entity.TeacherAttendanceStatusPresent,
		}

		existing := &entity.TeacherAttendance{ID: uuid.New()}
		mockRepo.EXPECT().GetByTeacherAndDate(gomock.Any(), teacherID, now).Return(existing, nil)

		err := u.CheckIn(context.Background(), attendance)
		assert.Error(t, err)
		assert.Equal(t, "already checked in for this date", err.Error())
	})

	t.Run("validation error", func(t *testing.T) {
		attendance := &entity.TeacherAttendance{
			TeacherID: uuid.Nil, // Invalid
		}

		err := u.CheckIn(context.Background(), attendance)
		assert.Error(t, err)
	})
}

func TestTeacherAttendanceUseCase_CheckOut(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTeacherAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewTeacherAttendanceUseCase(mockRepo, timeout)

	teacherID := uuid.New()
	now := time.Now()
	checkOutTime := now.Add(8 * time.Hour)

	t.Run("success", func(t *testing.T) {
		existing := &entity.TeacherAttendance{
			ID:        uuid.New(),
			TeacherID: teacherID,
		}

		mockRepo.EXPECT().GetByTeacherAndDate(gomock.Any(), teacherID, now).Return(existing, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		err := u.CheckOut(context.Background(), teacherID, now, checkOutTime)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByTeacherAndDate(gomock.Any(), teacherID, now).Return(nil, nil)

		err := u.CheckOut(context.Background(), teacherID, now, checkOutTime)
		assert.Error(t, err)
		assert.Equal(t, "attendance record not found", err.Error())
	})
}

func TestTeacherAttendanceUseCase_GetByTeacherAndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTeacherAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewTeacherAttendanceUseCase(mockRepo, timeout)

	teacherID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		expected := &entity.TeacherAttendance{ID: uuid.New(), TeacherID: teacherID}
		mockRepo.EXPECT().GetByTeacherAndDate(gomock.Any(), teacherID, now).Return(expected, nil)

		res, err := u.GetByTeacherAndDate(context.Background(), teacherID, now)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetByTeacherAndDate(gomock.Any(), teacherID, now).Return(nil, errors.New("db error"))

		res, err := u.GetByTeacherAndDate(context.Background(), teacherID, now)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTeacherAttendanceUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTeacherAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewTeacherAttendanceUseCase(mockRepo, timeout)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.TeacherAttendance{{ID: uuid.New()}}
		filter := map[string]interface{}{"status": "present"}
		mockRepo.EXPECT().List(gomock.Any(), filter).Return(expected, nil)

		res, err := u.List(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("error", func(t *testing.T) {
		filter := map[string]interface{}{"status": "present"}
		mockRepo.EXPECT().List(gomock.Any(), filter).Return(nil, errors.New("db error"))

		res, err := u.List(context.Background(), filter)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
