package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGradeHandler_InputGrade(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockGradingUseCase(ctrl)
	h := handler.NewGradeHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		tenantID := uuid.New()
		reqBody := map[string]interface{}{
			"tenant_id":     tenantID.String(),
			"assessment_id": uuid.New().String(),
			"student_id":    uuid.New().String(),
			"score":         85.0,
			"notes":         "Good job",
			"graded_by":     uuid.New().String(),
			"status":        "draft",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().InputGrade(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, grade *entity.Grade) error {
			assert.Equal(t, tenantID.String(), grade.TenantID)
			assert.Equal(t, 85.0, grade.Score)
			assert.Equal(t, "Good job", grade.Notes)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/grades", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.InputGrade(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid score", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":     uuid.New().String(),
			"assessment_id": uuid.New().String(),
			"student_id":    uuid.New().String(),
			// Missing score or invalid
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/grades", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.InputGrade(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGradeHandler_ApproveGrade(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockGradingUseCase(ctrl)
	h := handler.NewGradeHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		gradeID := uuid.New()
		approverID := uuid.New()
		
		reqBody := map[string]interface{}{
			"approved_by": approverID.String(),
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().ApproveGrade(gomock.Any(), gradeID, approverID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/grades/"+gradeID.String()+"/approve", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "id", Value: gradeID.String()}}

		h.ApproveGrade(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		gradeID := uuid.New()
		approverID := uuid.New()
		
		reqBody := map[string]interface{}{
			"approved_by": approverID.String(),
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().ApproveGrade(gomock.Any(), gradeID, approverID).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/grades/"+gradeID.String()+"/approve", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = []gin.Param{{Key: "id", Value: gradeID.String()}}

		h.ApproveGrade(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
