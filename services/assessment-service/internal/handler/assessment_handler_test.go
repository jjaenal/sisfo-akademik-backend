package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAssessmentHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockGradingUseCase(ctrl)
	h := handler.NewAssessmentHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":         "tenant-123",
			"subject_id":        uuid.New().String(),
			"teacher_id":        uuid.New().String(),
			"class_id":          uuid.New().String(),
			"grade_category_id": uuid.New().String(),
			"name":              "Midterm Exam",
			"max_score":         100.0,
			"description":       "Midterm examination",
			"date":              time.Now().Format(time.RFC3339),
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().CreateAssessment(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, assessment *entity.Assessment) error {
			assert.Equal(t, "tenant-123", assessment.TenantID)
			assert.Equal(t, "Midterm Exam", assessment.Name)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/assessments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - missing fields", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-123",
			// Missing name
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/assessments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request - invalid UUID", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":         "tenant-123",
			"subject_id":        "invalid-uuid",
			"teacher_id":        uuid.New().String(),
			"class_id":          uuid.New().String(),
			"grade_category_id": uuid.New().String(),
			"name":              "Midterm Exam",
			"max_score":         100.0,
			"date":              time.Now().Format(time.RFC3339),
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/assessments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":         "tenant-123",
			"subject_id":        uuid.New().String(),
			"teacher_id":        uuid.New().String(),
			"class_id":          uuid.New().String(),
			"grade_category_id": uuid.New().String(),
			"name":              "Midterm Exam",
			"max_score":         100.0,
			"description":       "Midterm examination",
			"date":              time.Now().Format(time.RFC3339),
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().CreateAssessment(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/assessments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
