package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGradingUseCase_CreateAssessment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := mocks.NewMockAssessmentRepository(ctrl)
	mockGradeRepo := mocks.NewMockGradeRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradingUseCase(mockAssessmentRepo, mockGradeRepo, timeout)

	t.Run("success", func(t *testing.T) {
		assessment := &entity.Assessment{
			TeacherID:       uuid.New(),
			SubjectID:       uuid.New(),
			ClassID:         uuid.New(),
			GradeCategoryID: uuid.New(),
			Name:            "Midterm Exam",
			MaxScore:        100,
			Date:            time.Now(),
		}

		mockAssessmentRepo.EXPECT().Create(gomock.Any(), assessment).Return(nil)

		err := u.CreateAssessment(context.Background(), assessment)
		assert.NoError(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		assessment := &entity.Assessment{
			Name: "", // Invalid
		}

		err := u.CreateAssessment(context.Background(), assessment)
		assert.Error(t, err)
		assert.Equal(t, "validation error", err.Error())
	})
}

func TestGradingUseCase_InputGrade(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := mocks.NewMockAssessmentRepository(ctrl)
	mockGradeRepo := mocks.NewMockGradeRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradingUseCase(mockAssessmentRepo, mockGradeRepo, timeout)

	assessmentID := uuid.New()
	studentID := uuid.New()

	t.Run("create new grade", func(t *testing.T) {
		grade := &entity.Grade{
			AssessmentID: assessmentID,
			StudentID:    studentID,
			Score:        85.5,
		}

		mockGradeRepo.EXPECT().GetByStudentAndAssessment(gomock.Any(), studentID, assessmentID).Return(nil, errors.New("not found"))
		mockGradeRepo.EXPECT().Create(gomock.Any(), grade).Return(nil)

		err := u.InputGrade(context.Background(), grade)
		assert.NoError(t, err)
	})

	t.Run("update existing grade", func(t *testing.T) {
		grade := &entity.Grade{
			AssessmentID: assessmentID,
			StudentID:    studentID,
			Score:        90.0,
		}
		existingID := uuid.New()
		existingGrade := &entity.Grade{
			ID:           existingID,
			AssessmentID: assessmentID,
			StudentID:    studentID,
			Score:        85.5,
		}

		mockGradeRepo.EXPECT().GetByStudentAndAssessment(gomock.Any(), studentID, assessmentID).Return(existingGrade, nil)
		mockGradeRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, g *entity.Grade) error {
			assert.Equal(t, existingID, g.ID)
			assert.Equal(t, 90.0, g.Score)
			return nil
		})

		err := u.InputGrade(context.Background(), grade)
		assert.NoError(t, err)
	})
}

func TestGradingUseCase_CalculateFinalScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssessmentRepo := mocks.NewMockAssessmentRepository(ctrl)
	mockGradeRepo := mocks.NewMockGradeRepository(ctrl)
	timeout := 2 * time.Second
	u := usecase.NewGradingUseCase(mockAssessmentRepo, mockGradeRepo, timeout)

	studentID := uuid.New()
	subjectID := uuid.New()
	semesterID := uuid.New()

	t.Run("success", func(t *testing.T) {
		assessment1 := &entity.Assessment{ID: uuid.New(), SubjectID: subjectID}
		assessment2 := &entity.Assessment{ID: uuid.New(), SubjectID: subjectID}
		assessments := []*entity.Assessment{assessment1, assessment2}

		mockAssessmentRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(assessments, nil)

		mockGradeRepo.EXPECT().GetByStudentAndAssessment(gomock.Any(), studentID, assessment1.ID).Return(&entity.Grade{Score: 80}, nil)
		mockGradeRepo.EXPECT().GetByStudentAndAssessment(gomock.Any(), studentID, assessment2.ID).Return(&entity.Grade{Score: 90}, nil)

		score, err := u.CalculateFinalScore(context.Background(), studentID, subjectID, semesterID)
		assert.NoError(t, err)
		assert.Equal(t, 85.0, score) // (80 + 90) / 2
	})
}
