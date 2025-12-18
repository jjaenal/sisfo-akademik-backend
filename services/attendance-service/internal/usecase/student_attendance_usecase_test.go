package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/service"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStudentAttendanceUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

	studentID := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()
	now := time.Now()
	tenantID := "tenant-1"

	t.Run("success", func(t *testing.T) {
		attendance := &entity.StudentAttendance{
			StudentID:      studentID,
			TenantID:       tenantID,
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

	t.Run("gps validation success", func(t *testing.T) {
		lat := -6.2000
		lon := 106.8000
		attendance := &entity.StudentAttendance{
			StudentID:        studentID,
			TenantID:         tenantID,
			ClassID:          classID,
			SemesterID:       semesterID,
			AttendanceDate:   now,
			Status:           entity.AttendanceStatusPresent,
			CheckInLatitude:  &lat,
			CheckInLongitude: &lon,
		}

		schoolLoc := &service.SchoolLocation{
			Latitude:  -6.2000,
			Longitude: 106.8000,
			Radius:    100,
		}
		mockSchoolService.EXPECT().GetLocation(gomock.Any(), tenantID).Return(schoolLoc, nil)
		mockRepo.EXPECT().Create(gomock.Any(), attendance).Return(nil)

		err := u.Create(context.Background(), attendance)
		assert.NoError(t, err)
	})

	t.Run("gps validation failed", func(t *testing.T) {
		lat := -7.2000 // Far away
		lon := 107.8000
		attendance := &entity.StudentAttendance{
			StudentID:        studentID,
			TenantID:         tenantID,
			ClassID:          classID,
			SemesterID:       semesterID,
			AttendanceDate:   now,
			Status:           entity.AttendanceStatusPresent,
			CheckInLatitude:  &lat,
			CheckInLongitude: &lon,
		}

		schoolLoc := &service.SchoolLocation{
			Latitude:  -6.2000,
			Longitude: 106.8000,
			Radius:    100,
		}
		mockSchoolService.EXPECT().GetLocation(gomock.Any(), tenantID).Return(schoolLoc, nil)

		err := u.Create(context.Background(), attendance)
		assert.Error(t, err)
		assert.Equal(t, "location too far from school", err.Error())
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
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

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
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

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
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

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

func TestStudentAttendanceUseCase_BulkCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

	studentID1 := uuid.New()
	studentID2 := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		attendances := []*entity.StudentAttendance{
			{
				StudentID:      studentID1,
				ClassID:        classID,
				SemesterID:     semesterID,
				AttendanceDate: now,
				Status:         entity.AttendanceStatusPresent,
			},
			{
				StudentID:      studentID2,
				ClassID:        classID,
				SemesterID:     semesterID,
				AttendanceDate: now,
				Status:         entity.AttendanceStatusAbsent,
			},
		}

		mockRepo.EXPECT().BulkCreate(gomock.Any(), attendances).Return(nil)

		err := u.BulkCreate(context.Background(), attendances)
		assert.NoError(t, err)
	})

	t.Run("empty list", func(t *testing.T) {
		err := u.BulkCreate(context.Background(), []*entity.StudentAttendance{})
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		attendances := []*entity.StudentAttendance{
			{
				StudentID: uuid.Nil, // Invalid
			},
		}

		err := u.BulkCreate(context.Background(), attendances)
		assert.Error(t, err)
	})
}

func TestStudentAttendanceUseCase_GetSummary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

	studentID := uuid.New()
	semesterID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedSummary := map[string]int{
			"present": 10,
			"absent":  2,
			"sick":    1,
		}
		mockRepo.EXPECT().GetSummary(gomock.Any(), studentID, semesterID).Return(expectedSummary, nil)

		res, err := u.GetSummary(context.Background(), studentID, semesterID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSummary, res)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetSummary(gomock.Any(), studentID, semesterID).Return(nil, errors.New("db error"))

		res, err := u.GetSummary(context.Background(), studentID, semesterID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestStudentAttendanceUseCase_GetDailyReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

	tenantID := uuid.New()
	date := time.Now()

	t.Run("success", func(t *testing.T) {
		expectedReport := []*entity.StudentAttendance{
			{ID: uuid.New()},
		}
		mockRepo.EXPECT().GetByTenantAndDate(gomock.Any(), tenantID, date).Return(expectedReport, nil)

		res, err := u.GetDailyReport(context.Background(), tenantID, date)
		assert.NoError(t, err)
		assert.Equal(t, expectedReport, res)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetByTenantAndDate(gomock.Any(), tenantID, date).Return(nil, errors.New("db error"))

		res, err := u.GetDailyReport(context.Background(), tenantID, date)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestStudentAttendanceUseCase_GetMonthlyReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

	tenantID := uuid.New()
	month := 1
	year := 2023

	t.Run("success", func(t *testing.T) {
		expectedReport := []*entity.StudentAttendance{
			{ID: uuid.New()},
		}
		// Expect GetByDateRange with start and end date
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-1 * time.Nanosecond)

		mockRepo.EXPECT().GetByDateRange(gomock.Any(), tenantID, startDate, endDate, nil).Return(expectedReport, nil)

		res, err := u.GetMonthlyReport(context.Background(), tenantID, month, year)
		assert.NoError(t, err)
		assert.Equal(t, expectedReport, res)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetByDateRange(gomock.Any(), tenantID, gomock.Any(), gomock.Any(), nil).Return(nil, errors.New("db error"))

		res, err := u.GetMonthlyReport(context.Background(), tenantID, month, year)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestStudentAttendanceUseCase_GetClassReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockStudentAttendanceRepository(ctrl)
	mockSchoolService := mocks.NewMockSchoolService(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewStudentAttendanceUseCase(mockRepo, mockSchoolService, timeout)

	tenantID := uuid.New()
	classID := uuid.New()
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	t.Run("success", func(t *testing.T) {
		expectedReport := []*entity.StudentAttendance{
			{ID: uuid.New()},
		}
		mockRepo.EXPECT().GetByDateRange(gomock.Any(), tenantID, startDate, endDate, &classID).Return(expectedReport, nil)

		res, err := u.GetClassReport(context.Background(), tenantID, classID, startDate, endDate)
		assert.NoError(t, err)
		assert.Equal(t, expectedReport, res)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.EXPECT().GetByDateRange(gomock.Any(), tenantID, startDate, endDate, &classID).Return(nil, errors.New("db error"))

		res, err := u.GetClassReport(context.Background(), tenantID, classID, startDate, endDate)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
