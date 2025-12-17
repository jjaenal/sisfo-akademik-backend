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
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAdmissionPeriodHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAdmissionPeriodUseCase(ctrl)
	h := handler.NewAdmissionPeriodHandler(mockUseCase)

	now := time.Now()
	nextWeek := now.Add(7 * 24 * time.Hour)

	t.Run("success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":       "PPDB 2025",
			"start_date": now,
			"end_date":   nextWeek,
			"is_active":  true,
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ interface{}, p *entity.AdmissionPeriod) error {
			assert.Equal(t, "PPDB 2025", p.Name)
			assert.True(t, p.IsActive)
			return nil
		})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/admission-periods", bytes.NewBuffer(body))
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
		c.Request, _ = http.NewRequest(http.MethodPost, "/admission-periods", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":       "PPDB 2025",
			"start_date": now,
			"end_date":   nextWeek,
			"is_active":  true,
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/admission-periods", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAdmissionPeriodHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAdmissionPeriodUseCase(ctrl)
	h := handler.NewAdmissionPeriodHandler(mockUseCase)

	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		expected := &entity.AdmissionPeriod{ID: id, Name: "PPDB 2025"}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/admission-periods/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/admission-periods/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAdmissionPeriodHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockAdmissionPeriodUseCase(ctrl)
	h := handler.NewAdmissionPeriodHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		expected := []*entity.AdmissionPeriod{
			{ID: uuid.New(), Name: "PPDB 2025"},
			{ID: uuid.New(), Name: "PPDB 2026"},
		}
		mockUseCase.EXPECT().List(gomock.Any()).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/admission-periods", nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		mockUseCase.EXPECT().List(gomock.Any()).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/admission-periods", nil)

		h.List(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
