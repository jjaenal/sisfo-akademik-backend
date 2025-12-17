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

func TestStudentAttendanceUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, timeout)

	studentID := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		attendance := &entity.StudentAttendance{
			StudentID:      studentID,
			ClassID:        classID,
			SemesterID:     semesterID,
			AttendanceDate: now,
			Status:         entity.AttendanceStatusPresent,
			Notes:          "Present",
		}

		mockRepo.EXPECT().Create(gomock.Any(), attendance).Return(nil)

		err := u.Create(context.Background(), attendance)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		attendance := &entity.StudentAttendance{
			StudentID: uuid.Nil, // Invalid
		}

		// Expect Create NOT to be called on repo
		err := u.Create(context.Background(), attendance)
		assert.Error(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		attendance := &entity.StudentAttendance{
			StudentID:      studentID,
			ClassID:        classID,
			SemesterID:     semesterID,
			AttendanceDate: now,
			Status:         entity.AttendanceStatusPresent,
		}

		mockRepo.EXPECT().Create(gomock.Any(), attendance).Return(errors.New("db error"))

		err := u.Create(context.Background(), attendance)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestStudentAttendanceUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, timeout)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedAttendance := &entity.StudentAttendance{ID: id}
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expectedAttendance, nil)

		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, expectedAttendance, res)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("not found"))

		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestStudentAttendanceUseCase_GetByClassAndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, timeout)

	classID := uuid.New()
	date := time.Now()

	t.Run("success", func(t *testing.T) {
		expectedList := []*entity.StudentAttendance{
			{ID: uuid.New()},
			{ID: uuid.New()},
		}
		mockRepo.EXPECT().GetByClassAndDate(gomock.Any(), classID, date).Return(expectedList, nil)

		res, err := u.GetByClassAndDate(context.Background(), classID, date)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedList), len(res))
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetByClassAndDate(gomock.Any(), classID, date).Return(nil, errors.New("db error"))

		res, err := u.GetByClassAndDate(context.Background(), classID, date)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestStudentAttendanceUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, timeout)

	id := uuid.New()
	studentID := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		attendance := &entity.StudentAttendance{
			ID:             id,
			StudentID:      studentID,
			ClassID:        classID,
			SemesterID:     semesterID,
			AttendanceDate: now,
			Status:         entity.AttendanceStatusPresent,
		}

		mockRepo.EXPECT().Update(gomock.Any(), attendance).Return(nil)

		err := u.Update(context.Background(), attendance)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		attendance := &entity.StudentAttendance{
			ID:        id,
			StudentID: uuid.Nil, // Invalid
		}

		err := u.Update(context.Background(), attendance)
		assert.Error(t, err)
	})
}
