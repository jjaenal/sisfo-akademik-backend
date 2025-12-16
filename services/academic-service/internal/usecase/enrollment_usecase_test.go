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

func TestEnrollmentUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEnrollmentRepository(ctrl)
	mockClassRepo := mocks.NewMockClassRepository(ctrl)
	u := usecase.NewEnrollmentUseCase(mockRepo, mockClassRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()
	classID := uuid.New()
	studentID := uuid.New()

	validEnrollment := &entity.Enrollment{
		ID:        id,
		TenantID:  tenantID,
		ClassID:   classID,
		StudentID: studentID,
		Status:    "active",
	}

	validClass := &entity.Class{
		ID:       classID,
		TenantID: tenantID,
		Name:     "Class 1A",
		Capacity: 30,
	}

	t.Run("Enroll Success", func(t *testing.T) {
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(validClass, nil)
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return([]entity.Enrollment{}, nil)
		mockRepo.EXPECT().Enroll(gomock.Any(), validEnrollment).Return(nil)
		err := u.Enroll(context.Background(), validEnrollment)
		assert.NoError(t, err)
	})

	t.Run("Enroll Validation Error", func(t *testing.T) {
		invalidEnrollment := &entity.Enrollment{}
		err := u.Enroll(context.Background(), invalidEnrollment)
		assert.Error(t, err)
	})

	t.Run("Enroll Class Not Found", func(t *testing.T) {
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(nil, nil)
		err := u.Enroll(context.Background(), validEnrollment)
		assert.Error(t, err)
		assert.Equal(t, "class not found", err.Error())
	})

	t.Run("Enroll Capacity Exceeded", func(t *testing.T) {
		smallClass := &entity.Class{
			ID:       classID,
			TenantID: tenantID,
			Capacity: 1,
		}
		existingEnrollment := entity.Enrollment{ID: uuid.New(), ClassID: classID}

		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(smallClass, nil)
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return([]entity.Enrollment{existingEnrollment}, nil)
		
		err := u.Enroll(context.Background(), validEnrollment)
		assert.Error(t, err)
		assert.Equal(t, "class capacity exceeded", err.Error())
	})

	t.Run("Enroll Repo Error", func(t *testing.T) {
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(validClass, nil)
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return([]entity.Enrollment{}, nil)
		mockRepo.EXPECT().Enroll(gomock.Any(), validEnrollment).Return(assert.AnError)
		err := u.Enroll(context.Background(), validEnrollment)
		assert.Error(t, err)
	})

	t.Run("Unenroll Success", func(t *testing.T) {
		mockRepo.EXPECT().Unenroll(gomock.Any(), id).Return(nil)
		err := u.Unenroll(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("Unenroll Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Unenroll(gomock.Any(), id).Return(assert.AnError)
		err := u.Unenroll(context.Background(), id)
		assert.Error(t, err)
	})

	t.Run("BulkEnroll Success", func(t *testing.T) {
		studentIDs := []uuid.UUID{uuid.New(), uuid.New()}
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(validClass, nil)
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return([]entity.Enrollment{}, nil)
		mockRepo.EXPECT().BulkEnroll(gomock.Any(), gomock.Any()).Return(nil)

		err := u.BulkEnroll(context.Background(), classID, studentIDs)
		assert.NoError(t, err)
	})

	t.Run("BulkEnroll Validation Error (Empty)", func(t *testing.T) {
		err := u.BulkEnroll(context.Background(), classID, []uuid.UUID{})
		assert.Error(t, err)
		assert.Equal(t, "student IDs are required", err.Error())
	})

	t.Run("BulkEnroll Class Not Found", func(t *testing.T) {
		studentIDs := []uuid.UUID{uuid.New()}
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(nil, nil)
		
		err := u.BulkEnroll(context.Background(), classID, studentIDs)
		assert.Error(t, err)
		assert.Equal(t, "class not found", err.Error())
	})

	t.Run("BulkEnroll Capacity Exceeded", func(t *testing.T) {
		smallClass := &entity.Class{
			ID:       classID,
			TenantID: tenantID,
			Capacity: 1,
		}
		studentIDs := []uuid.UUID{uuid.New(), uuid.New()} // 2 students
		
		mockClassRepo.EXPECT().GetByID(gomock.Any(), classID).Return(smallClass, nil)
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return([]entity.Enrollment{}, nil) // 0 existing

		err := u.BulkEnroll(context.Background(), classID, studentIDs)
		assert.Error(t, err)
		assert.Equal(t, "class capacity exceeded", err.Error())
	})

	t.Run("GetByID Success", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validEnrollment, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validEnrollment, res)
	})

	t.Run("GetByID Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, assert.AnError)
		res, err := u.GetByID(context.Background(), id)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("ListByClass Success", func(t *testing.T) {
		list := []entity.Enrollment{*validEnrollment}
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return(list, nil)
		res, err := u.ListByClass(context.Background(), classID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})

	t.Run("ListByClass Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().ListByClass(gomock.Any(), classID).Return(nil, assert.AnError)
		res, err := u.ListByClass(context.Background(), classID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("ListByStudent Success", func(t *testing.T) {
		list := []entity.Enrollment{*validEnrollment}
		mockRepo.EXPECT().ListByStudent(gomock.Any(), studentID).Return(list, nil)
		res, err := u.ListByStudent(context.Background(), studentID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})

	t.Run("ListByStudent Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().ListByStudent(gomock.Any(), studentID).Return(nil, assert.AnError)
		res, err := u.ListByStudent(context.Background(), studentID)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("UpdateStatus Success", func(t *testing.T) {
		status := "moved"
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), id, status).Return(nil)
		err := u.UpdateStatus(context.Background(), id, status)
		assert.NoError(t, err)
	})

	t.Run("UpdateStatus Validation Error", func(t *testing.T) {
		status := ""
		err := u.UpdateStatus(context.Background(), id, status)
		assert.Error(t, err)
		assert.Equal(t, "status is required", err.Error())
	})

	t.Run("UpdateStatus Repo Error", func(t *testing.T) {
		status := "moved"
		mockRepo.EXPECT().UpdateStatus(gomock.Any(), id, status).Return(assert.AnError)
		err := u.UpdateStatus(context.Background(), id, status)
		assert.Error(t, err)
	})
}
