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

func TestStudentHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentUseCase(ctrl)
	h := handler.NewStudentHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "John Doe",
			"status":    "active",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(body))

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
		c.Request, _ = http.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UseCase Error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "John Doe",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("usecase error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/students", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestStudentHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentUseCase(ctrl)
	h := handler.NewStudentHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		student := &entity.Student{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(student, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/students/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/students/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestStudentHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockStudentUseCase(ctrl)
	h := handler.NewStudentHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.EXPECT().List(gomock.Any(), "tenant-1", 10, 0).Return([]entity.Student{}, 0, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/students?tenant_id=tenant-1", nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
