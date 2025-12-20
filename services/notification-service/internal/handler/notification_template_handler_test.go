package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationTemplateHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationTemplateUseCase(ctrl)
	h := handler.NewNotificationTemplateHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		req := entity.NotificationTemplate{
			Name:         "Welcome",
			Channel:      entity.NotificationChannelEmail,
			BodyTemplate: "Hello",
		}
		
		mockUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		body, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPost, "/templates", bytes.NewBuffer(body))

		h.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		c.Request, _ = http.NewRequest(http.MethodPost, "/templates", bytes.NewBufferString("invalid json"))

		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestNotificationTemplateHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationTemplateUseCase(ctrl)
	h := handler.NewNotificationTemplateHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		template := &entity.NotificationTemplate{
			ID:   id,
			Name: "Welcome",
		}

		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(template, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/templates/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New()
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/templates/"+id.String(), nil)

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestNotificationTemplateHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationTemplateUseCase(ctrl)
	h := handler.NewNotificationTemplateHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		req := entity.NotificationTemplate{
			Name:         "Welcome Updated",
			Channel:      entity.NotificationChannelEmail,
			BodyTemplate: "Hello Updated",
		}

		mockUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		
		body, _ := json.Marshal(req)
		c.Request, _ = http.NewRequest(http.MethodPut, "/templates/"+id.String(), bytes.NewBuffer(body))

		h.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("bad request - invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}
		
		c.Request, _ = http.NewRequest(http.MethodPut, "/templates/invalid", nil)

		h.Update(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}



func TestNotificationTemplateHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationTemplateUseCase(ctrl)
	h := handler.NewNotificationTemplateHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		mockUseCase.EXPECT().Delete(gomock.Any(), id).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest(http.MethodDelete, "/templates/"+id.String(), nil)

		h.Delete(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestNotificationTemplateHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationTemplateUseCase(ctrl)
	h := handler.NewNotificationTemplateHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		templates := []*entity.NotificationTemplate{
			{Name: "T1"},
			{Name: "T2"},
		}

		mockUseCase.EXPECT().List(gomock.Any()).Return(templates, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/templates", nil)

		h.List(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
