package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestReportCardIntegration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Repositories
	mockReportRepo := mocks.NewMockReportCardRepository(ctrl)
	mockGradeRepo := mocks.NewMockGradeRepository(ctrl)
	mockAssessmentRepo := mocks.NewMockAssessmentRepository(ctrl)
	mockCategoryRepo := mocks.NewMockGradeCategoryRepository(ctrl)
	mockFileStorage := mocks.NewMockFileStorage(ctrl)

	// Real UseCase with Mock Repos
	u := usecase.NewReportCardUseCase(mockReportRepo, mockGradeRepo, mockAssessmentRepo, mockCategoryRepo, mockFileStorage)

	// Real Handler with Real UseCase
	h := handler.NewReportCardHandler(u)

	// Setup Router
	r := gin.New()
	reportCards := r.Group("/api/v1/report-cards")
	{
		reportCards.POST("/generate", h.Generate)
		// reportCards.GET("/:id", h.GetByID) // Add other routes as needed
	}

	t.Run("Generate Report Card Success", func(t *testing.T) {
		tenantID := uuid.New().String()
		studentID := uuid.New()
		classID := uuid.New()
		semesterID := uuid.New()

		reqBody := map[string]interface{}{
			"tenant_id":   tenantID,
			"student_id":  studentID.String(),
			"class_id":    classID.String(),
			"semester_id": semesterID.String(),
		}
		body, _ := json.Marshal(reqBody)

		// 1. Check existing
		mockReportRepo.EXPECT().GetByStudentAndSemester(gomock.Any(), studentID, semesterID).Return(nil, nil)

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
		mockGradeRepo.EXPECT().GetByStudentID(gomock.Any(), studentID).Return(grades, nil)

		// 3. Fetch Assessment
		assessment := &entity.Assessment{
			ID:              assessmentID,
			GradeCategoryID: categoryID,
			SubjectID:       subjectID,
			ClassID:         classID,
			MaxScore:        100,
		}
		mockAssessmentRepo.EXPECT().GetByID(gomock.Any(), assessmentID).Return(assessment, nil)

		// 4. Fetch Category
		category := &entity.GradeCategory{
			ID:     categoryID,
			Weight: 100,
		}
		mockCategoryRepo.EXPECT().GetByID(gomock.Any(), categoryID).Return(category, nil)

		// 5. Upload PDF
		mockFileStorage.EXPECT().Upload(gomock.Any(), gomock.Any(), gomock.Any()).Return("http://example.com/report.pdf", nil)

		// 6. Create Report Card
		mockReportRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, rc *entity.ReportCard) error {
			assert.Equal(t, studentID, rc.StudentID)
			assert.Equal(t, entity.ReportCardStatusGenerated, rc.Status)
			assert.Len(t, rc.Details, 1)
			assert.Equal(t, 85.0, rc.Details[0].FinalScore)
			assert.Equal(t, "http://example.com/report.pdf", rc.PDFUrl)
			return nil
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/report-cards/generate", bytes.NewBuffer(body))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp struct {
			Success bool              `json:"success"`
			Data    entity.ReportCard `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.True(t, resp.Success)
		assert.Equal(t, string(entity.ReportCardStatusGenerated), string(resp.Data.Status))
		assert.Equal(t, "http://example.com/report.pdf", resp.Data.PDFUrl)
	})
}
