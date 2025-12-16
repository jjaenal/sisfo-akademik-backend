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
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/mocks"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestScheduleHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockScheduleUseCase(ctrl)
	h := handler.NewScheduleHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":   "tenant-1",
			"class_id":    uuid.New().String(),
			"subject_id":  uuid.New().String(),
			"teacher_id":  uuid.New().String(),
			"day_of_week": 1,
			"start_time":  "08:00",
			"end_time":    "10:00",
			"room":        "Room 101",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/schedules", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			// Missing required fields
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/schedules", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":   "tenant-1",
			"class_id":    uuid.New().String(),
			"subject_id":  uuid.New().String(),
			"teacher_id":  uuid.New().String(),
			"day_of_week": 1,
			"start_time":  "08:00",
			"end_time":    "10:00",
			"room":        "Room 101",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("usecase error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/schedules", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestScheduleHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockScheduleUseCase(ctrl)
	h := handler.NewScheduleHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		schedule := &entity.Schedule{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(schedule, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/schedules/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/schedules/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestScheduleHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockScheduleUseCase(ctrl)
	h := handler.NewScheduleHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.EXPECT().List(gomock.Any(), "tenant-1", 10, 0).Return([]entity.Schedule{}, 0, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/schedules?tenant_id=tenant-1", nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing Tenant ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/schedules", nil)

		h.List(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
