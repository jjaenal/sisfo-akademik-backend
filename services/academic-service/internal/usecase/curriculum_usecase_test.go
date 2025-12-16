package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCurriculumUseCase_CRUD(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCurriculumRepository(ctrl)
	u := usecase.NewCurriculumUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	id := uuid.New()
	subjectID := uuid.New()

	validCurriculum := &entity.Curriculum{
		ID:       id,
		TenantID: tenantID,
		Name:     "K-13",
		Year:     2013,
		IsActive: true,
	}

	t.Run("Create", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validCurriculum).Return(nil)
		err := u.Create(context.Background(), validCurriculum)
		assert.NoError(t, err)
	})

	t.Run("Create Validation Error", func(t *testing.T) {
		invalidCurriculum := &entity.Curriculum{}
		err := u.Create(context.Background(), invalidCurriculum)
		assert.Error(t, err)
	})

	t.Run("Create Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Create(gomock.Any(), validCurriculum).Return(assert.AnError)
		err := u.Create(context.Background(), validCurriculum)
		assert.Error(t, err)
	})

	t.Run("GetByID", func(t *testing.T) {
		mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(validCurriculum, nil)
		res, err := u.GetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, validCurriculum, res)
	})

	t.Run("List", func(t *testing.T) {
		list := []entity.Curriculum{*validCurriculum}
		total := 1
		mockRepo.EXPECT().List(gomock.Any(), tenantID, 10, 0).Return(list, total, nil)
		res, count, err := u.List(context.Background(), tenantID, 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, list, res)
		assert.Equal(t, total, count)
	})

	t.Run("Update", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validCurriculum).Return(nil)
		err := u.Update(context.Background(), validCurriculum)
		assert.NoError(t, err)
	})

	t.Run("Update Validation Error", func(t *testing.T) {
		invalidCurriculum := &entity.Curriculum{}
		err := u.Update(context.Background(), invalidCurriculum)
		assert.Error(t, err)
	})

	t.Run("Update Repo Error", func(t *testing.T) {
		mockRepo.EXPECT().Update(gomock.Any(), validCurriculum).Return(assert.AnError)
		err := u.Update(context.Background(), validCurriculum)
		assert.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
		err := u.Delete(context.Background(), id)
		assert.NoError(t, err)
	})

	t.Run("AddSubject", func(t *testing.T) {
		cs := &entity.CurriculumSubject{
			CurriculumID: id,
			SubjectID:    subjectID,
			GradeLevel:   10,
			Semester:     1,
		}
		mockRepo.EXPECT().AddSubject(gomock.Any(), cs).Return(nil)
		err := u.AddSubject(context.Background(), cs)
		assert.NoError(t, err)
	})

	t.Run("AddSubject Validation Error", func(t *testing.T) {
		cs := &entity.CurriculumSubject{}
		err := u.AddSubject(context.Background(), cs)
		assert.Error(t, err)
	})

	t.Run("AddSubject Repo Error", func(t *testing.T) {
		cs := &entity.CurriculumSubject{
			CurriculumID: id,
			SubjectID:    subjectID,
			GradeLevel:   10,
			Semester:     1,
		}
		mockRepo.EXPECT().AddSubject(gomock.Any(), cs).Return(assert.AnError)
		err := u.AddSubject(context.Background(), cs)
		assert.Error(t, err)
	})

	t.Run("RemoveSubject", func(t *testing.T) {
		relationID := uuid.New()
		mockRepo.EXPECT().RemoveSubject(gomock.Any(), relationID).Return(nil)
		err := u.RemoveSubject(context.Background(), relationID)
		assert.NoError(t, err)
	})
}

func TestCurriculumUseCase_GradingRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCurriculumRepository(ctrl)
	u := usecase.NewCurriculumUseCase(mockRepo, time.Second*2)

	tenantID := "tenant-1"
	curriculumID := uuid.New()

	validRule := &entity.GradingRule{
		TenantID:     tenantID,
		CurriculumID: curriculumID,
		Grade:        "A",
		MinScore:     85.0,
		MaxScore:     100.0,
		Points:       4.0,
	}

	t.Run("AddGradingRule Success", func(t *testing.T) {
		mockRepo.EXPECT().AddGradingRule(gomock.Any(), validRule).Return(nil)
		err := u.AddGradingRule(context.Background(), validRule)
		assert.NoError(t, err)
	})

	t.Run("AddGradingRule Validation Error", func(t *testing.T) {
		invalidRule := &entity.GradingRule{
			TenantID: tenantID,
			Grade:    "", // Invalid
		}
		err := u.AddGradingRule(context.Background(), invalidRule)
		assert.Error(t, err)
	})

	t.Run("AddGradingRule Repository Error", func(t *testing.T) {
		mockRepo.EXPECT().AddGradingRule(gomock.Any(), validRule).Return(errors.New("db error"))
		err := u.AddGradingRule(context.Background(), validRule)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("GetGradingRules Success", func(t *testing.T) {
		expectedRules := []entity.GradingRule{
			{Grade: "A", MinScore: 85},
			{Grade: "B", MinScore: 70},
		}
		mockRepo.EXPECT().ListGradingRules(gomock.Any(), curriculumID).Return(expectedRules, nil)

		rules, err := u.GetGradingRules(context.Background(), curriculumID)
		assert.NoError(t, err)
		assert.Len(t, rules, 2)
		assert.Equal(t, "A", rules[0].Grade)
	})

	t.Run("GetGradingRules Repository Error", func(t *testing.T) {
		mockRepo.EXPECT().ListGradingRules(gomock.Any(), curriculumID).Return(nil, errors.New("db error"))

		rules, err := u.GetGradingRules(context.Background(), curriculumID)
		assert.Error(t, err)
		assert.Nil(t, rules)
	})
}
