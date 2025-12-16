package handler_test

import (
	"bytes"
	"encoding/json"
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

func TestClassHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockClassUseCase(ctrl)
	h := handler.NewClassHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "Class 10A",
			"level":     10,
			"major":     "Science",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/classes", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			// Missing required name
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/classes", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestClassHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockClassUseCase(ctrl)
	h := handler.NewClassHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		class := &entity.Class{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(class, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/classes/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestClassHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockClassUseCase(ctrl)
	h := handler.NewClassHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.EXPECT().List(gomock.Any(), "tenant-1", 10, 0).Return([]entity.Class{}, 0, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/classes?tenant_id=tenant-1", nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
