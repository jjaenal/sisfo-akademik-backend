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
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNotificationHandler_Send(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationUseCase(ctrl)
	h := handler.NewNotificationHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		req := usecase.SendNotificationRequest{
			Channel:   entity.NotificationChannelEmail,
			Recipient: "test@example.com",
			Subject:   "Subject",
			Body:      "Body",
		}
		
		mockUseCase.EXPECT().Send(gomock.Any(), &req).Return(nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/notifications/send", bytes.NewBuffer(body))

		h.Send(c)

		assert.Equal(t, http.StatusAccepted, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/notifications/send", bytes.NewBufferString("invalid json"))

		h.Send(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		req := usecase.SendNotificationRequest{
			Channel:   entity.NotificationChannelEmail,
			Recipient: "test@example.com",
			Subject:   "Subject",
			Body:      "Body",
		}

		mockUseCase.EXPECT().Send(gomock.Any(), &req).Return(errors.New("internal error"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/notifications/send", bytes.NewBuffer(body))

		h.Send(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNotificationHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationUseCase(ctrl)
	h := handler.NewNotificationHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		notification := &entity.Notification{ID: id}

		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(notification, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/notifications/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/notifications/invalid-uuid", nil)
		c.Params = gin.Params{{Key: "id", Value: "invalid-uuid"}}

		h.GetByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New()
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/notifications/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		id := uuid.New()
		mockUseCase.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/notifications/"+id.String(), nil)
		c.Params = gin.Params{{Key: "id", Value: id.String()}}

		h.GetByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNotificationHandler_ListByRecipient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationUseCase(ctrl)
	h := handler.NewNotificationHandler(mockUseCase)

	t.Run("success", func(t *testing.T) {
		recipient := "user@example.com"
		notifications := []*entity.Notification{{Recipient: recipient}}

		mockUseCase.EXPECT().ListByRecipient(gomock.Any(), recipient).Return(notifications, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/notifications/recipient?recipient="+recipient, nil)

		h.ListByRecipient(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing recipient", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/notifications/recipient", nil)

		h.ListByRecipient(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		recipient := "user@example.com"
		mockUseCase.EXPECT().ListByRecipient(gomock.Any(), recipient).Return(nil, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/notifications/recipient?recipient="+recipient, nil)

		h.ListByRecipient(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
