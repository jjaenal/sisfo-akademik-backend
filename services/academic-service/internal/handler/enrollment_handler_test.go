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

func TestEnrollmentHandler_Enroll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEnrollmentUseCase(ctrl)
	h := handler.NewEnrollmentHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id":  "tenant-1",
			"class_id":   uuid.New().String(),
			"student_id": uuid.New().String(),
			"status":     "active",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Enroll(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/enrollments", bytes.NewBuffer(body))

		h.Enroll(c)

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
		c.Request, _ = http.NewRequest(http.MethodPost, "/enrollments", bytes.NewBuffer(body))

		h.Enroll(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestEnrollmentHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEnrollmentUseCase(ctrl)
	h := handler.NewEnrollmentHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		enrollment := &entity.Enrollment{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(enrollment, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/enrollments/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestEnrollmentHandler_UpdateStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockEnrollmentUseCase(ctrl)
	h := handler.NewEnrollmentHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"status": "completed",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().UpdateStatus(gomock.Any(), id, "completed").Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPatch, "/enrollments/"+id.String()+"/status", bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.UpdateStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"status": "completed",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().UpdateStatus(gomock.Any(), id, "completed").Return(errors.New("usecase error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPatch, "/enrollments/"+id.String()+"/status", bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.UpdateStatus(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
