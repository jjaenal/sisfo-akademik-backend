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

func TestSemesterUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSemesterRepository(ctrl)
	u := usecase.NewSemesterUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	academicYearID := uuid.New()
	id := uuid.New()

	validSemester := &entity.Semester{
		ID:             id,
		TenantID:       tenantID,
		AcademicYearID: academicYearID,
		Name:           "Semester 1",
		StartDate:      time.Now(),
		EndDate:        time.Now().AddDate(0, 6, 0),
		IsActive:       true,
	}

	t.Run("Create Active", func(t *testing.T) {
		existingSemesters := []entity.Semester{
			{ID: uuid.New(), IsActive: true, AcademicYearID: academicYearID},
		}
		mockRepo.EXPECT().ListByAcademicYear(gomock.Any(), academicYearID).Return(existingSemesters, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().Create(gomock.Any(), validSemester).Return(nil)

		err := u.Create(context.Background(), validSemester)
		assert.NoError(t, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validSemester, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validSemester, res)
	})

	t.Run("List", func(t *testing.T) {
		list := []entity.Semester{*validSemester}
		mockRepo.EXPECT().List(gomock.Any(), tenantID).Return(list, nil)
		res, err := u.List(context.Background(), tenantID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})

	t.Run("ListByAcademicYear", func(t *testing.T) {
		list := []entity.Semester{*validSemester}
		mockRepo.EXPECT().ListByAcademicYear(gomock.Any(), academicYearID).Return(list, nil)
		res, err := u.ListByAcademicYear(context.Background(), academicYearID)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
	})

	t.Run("Update", func(t *testing.T) {
		// Mock logic for update (if IsActive is true, it calls deactivateOthers)
		existingSemesters := []entity.Semester{}
		mockRepo.EXPECT().ListByAcademicYear(gomock.Any(), academicYearID).Return(existingSemesters, nil)
		mockRepo.EXPECT().Update(gomock.Any(), validSemester).Return(nil)
		
		err := u.Update(context.Background(), validSemester)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("SetActive", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validSemester, nil)
		// It calls deactivateOthers
		existingSemesters := []entity.Semester{}
		mockRepo.EXPECT().ListByAcademicYear(gomock.Any(), academicYearID).Return(existingSemesters, nil)
		mockRepo.EXPECT().Update(gomock.Any(), validSemester).Return(nil)

		err := u.SetActive(context.Background(), id)
		assert.NoError(t, err)
	})
	
	t.Run("SetActive Not Found", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)
		err := u.SetActive(context.Background(), id)
		assert.Error(t, err)
	})
}
