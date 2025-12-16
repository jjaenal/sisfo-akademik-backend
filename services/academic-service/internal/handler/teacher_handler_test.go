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

func TestTeacherHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTeacherUseCase(ctrl)
	h := handler.NewTeacherHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"tenant_id": "tenant-1",
			"name":      "Mr. Teacher",
		}
		body, _ := json.Marshal(reqBody)

		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/teachers", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTeacherHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTeacherUseCase(ctrl)
	h := handler.NewTeacherHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		teacher := &entity.Teacher{ID: id}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(teacher, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/teachers/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTeacherHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTeacherUseCase(ctrl)
	h := handler.NewTeacherHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "Mr. Teacher Updated",
		}
		body, _ := json.Marshal(reqBody)

		existing := &entity.Teacher{ID: id, Name: "Mr. Teacher"}
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
		mockUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPut, "/teachers/"+id.String(), bytes.NewBuffer(body))
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTeacherHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTeacherUseCase(ctrl)
	h := handler.NewTeacherHandler(mockUseCase)

	id := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockUseCase.EXPECT().Delete(gomock.Any(), id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/teachers/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTeacherHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockTeacherUseCase(ctrl)
	h := handler.NewTeacherHandler(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.EXPECT().List(gomock.Any(), "tenant-1", 10, 0).Return([]entity.Teacher{}, 0, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/teachers?tenant_id=tenant-1", nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
