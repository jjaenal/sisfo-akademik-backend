package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/usecase"
)

func TestReportCardUseCase_Generate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportCardRepository(ctrl)
	mockGradeRepo := mocks.NewMockGradeRepository(ctrl)
	mockAssessmentRepo := mocks.NewMockAssessmentRepository(ctrl)
	mockCategoryRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)

	u := usecase.NewReportCardUseCase(mockReportRepo, mockGradeRepo, mockAssessmentRepo, mockCategoryRepo, mockFileStorage)

	ctx := context.Background()
	tenantID := uuid.New()
	studentID := uuid.New()
	classID := uuid.New()
	semesterID := uuid.New()

	t.Run("success_create_new", func(t *testing.T) {
		// 1. Check existing
		mockReportRepo.EXPECT().GetByStudentAndSemester(ctx, studentID, semesterID).Return(nil, nil)

		// 2. Fetch Grades
		assessmentID := uuid.New()
		categoryID := uuid.New()
		subjectID := uuid.New()
		
		grades := []*entity.Grade{
			{
				ID:           uuid.New(),
				AssessmentID: assessmentID,
				Score:        85.0,
			},
		}
		mockGradeRepo.EXPECT().List(ctx, map[string]interface{}{"student_id": studentID}).Return(grades, nil)

		// 3. Fetch Assessment
		assessment := &entity.Assessment{
			ID:              assessmentID,
			GradeCategoryID: categoryID,
			SubjectID:       subjectID,
			ClassID:         classID,
			MaxScore:        100,
		}
		mockAssessmentRepo.EXPECT().GetByID(ctx, assessmentID).Return(assessment, nil)

		// 4. Fetch Category
		category := &entity.GradeCategory{
			ID:     categoryID,
			Weight: 100,
		}
		mockCategoryRepo.EXPECT().GetByID(ctx, categoryID).Return(category, nil)

		// 5. Upload PDF
		mockFileStorage.EXPECT().Upload(ctx, gomock.Any(), gomock.Any()).Return("http://example.com/report.pdf", nil)

		// 6. Create Report Card
		mockReportRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, rc *entity.ReportCard) error {
			assert.Equal(t, studentID, rc.StudentID)
			assert.Equal(t, entity.ReportCardStatusGenerated, rc.Status)
			assert.Len(t, rc.Details, 1)
			assert.Equal(t, 85.0, rc.Details[0].FinalScore)
			assert.Equal(t, "B", rc.Details[0].GradeLetter)
			assert.Equal(t, "http://example.com/report.pdf", rc.PDFUrl)
			return nil
		})

		rc, err := u.Generate(ctx, tenantID, studentID, classID, semesterID)
		assert.NoError(t, err)
		assert.NotNil(t, rc)
	})

	t.Run("error_already_published", func(t *testing.T) {
		existing := &entity.ReportCard{
			ID:     uuid.New(),
			Status: entity.ReportCardStatusPublished,
		}
		mockReportRepo.EXPECT().GetByStudentAndSemester(ctx, studentID, semesterID).Return(existing, nil)

		rc, err := u.Generate(ctx, tenantID, studentID, classID, semesterID)
		assert.Error(t, err)
		assert.Nil(t, rc)
		assert.Equal(t, "report card already published", err.Error())
	})
}

func TestReportCardUseCase_GetPDF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReportRepo := mocks.NewMockReportCardRepository(ctrl)
	mockGradeRepo := mocks.NewMockGradeRepository(ctrl)
	mockAssessmentRepo := mocks.NewMockAssessmentRepository(ctrl)
	mockCategoryRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)

	u := usecase.NewReportCardUseCase(mockReportRepo, mockGradeRepo, mockAssessmentRepo, mockCategoryRepo, mockFileStorage)
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		now := time.Now()
		report := &entity.ReportCard{
			ID:           id,
			StudentID:    uuid.New(),
			ClassID:      uuid.New(),
			SemesterID:   uuid.New(),
			Status:       entity.ReportCardStatusGenerated,
			GeneratedAt:  &now,
			GPA:          3.5,
			TotalCredits: 20,
			Details: []entity.ReportCardDetail{
				{
					SubjectName: "Math",
					Credit:      3,
					FinalScore:  85.0,
					GradeLetter: "B",
				},
			},
		}

		mockReportRepo.EXPECT().GetByID(ctx, id).Return(report, nil)

		pdfBytes, err := u.GetPDF(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, pdfBytes)
		assert.NotEmpty(t, pdfBytes)
		// Check PDF header signature (simple check)
		assert.Equal(t, "%PDF", string(pdfBytes[:4]))
	})

	t.Run("not_found", func(t *testing.T) {
		mockReportRepo.EXPECT().GetByID(ctx, id).Return(nil, nil)

		pdfBytes, err := u.GetPDF(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, pdfBytes)
		assert.Equal(t, "report card not found", err.Error())
	})

	t.Run("repo_error", func(t *testing.T) {
		mockReportRepo.EXPECT().GetByID(ctx, id).Return(nil, errors.New("db error"))

		pdfBytes, err := u.GetPDF(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, pdfBytes)
		assert.Equal(t, "db error", err.Error())
	})
}
