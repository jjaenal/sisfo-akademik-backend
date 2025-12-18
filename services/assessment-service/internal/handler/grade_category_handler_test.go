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

func TestGradeCategoryHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockGradeCategoryUseCase(ctrl)
	h := handler.NewGradeCategoryHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":   "tenant-123",
			"name":        "Quiz",
			"description": "Daily Quiz",
			"weight":      10,
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, cat *entity.GradeCategory) error {
			assert.Equal(t, "tenant-123", cat.TenantID)
			assert.Equal(t, "Quiz", cat.Name)
			assert.Equal(t, "Daily Quiz", cat.Description)
			assert.Equal(t, 10.0, cat.Weight)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/grade-categories", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "", // Required
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/grade-categories", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":   "tenant-123",
			"name":        "Quiz",
			"description": "Daily Quiz",
			"weight":      10,
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/grade-categories", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGradeCategoryHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockGradeCategoryUseCase(ctrl)
	h := handler.NewGradeCategoryHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := &entity.GradeCategory{ID: id, Name: "Quiz"}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/grade-categories/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/grade-categories/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGradeCategoryHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockGradeCategoryUseCase(ctrl)
	h := handler.NewGradeCategoryHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.GradeCategory{
			{ID: uuid.New(), Name: "Quiz"},
			{ID: uuid.New(), Name: "Exam"},
		}
		tenantID := "tenant-1"
		mockUseCase.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/grade-categories?tenant_id="+tenantID, nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		tenantID := "tenant-1"
		mockUseCase.EXPECT().GetByTenantID(gomock.Any(), tenantID).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/grade-categories?tenant_id="+tenantID, nil)

		h.List(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
