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
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestWebhookHandler_HandleWebhook_Email(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationUseCase(ctrl)
	h := handler.NewWebhookHandler(mockUseCase)

	notificationID := uuid.New()

	tests := []struct {
		name           string
		provider       string
		payload        interface{}
		setupMock      func()
		expectedStatus int
	}{
		{
			name:     "Success - Email Delivered",
			provider: "email",
			payload: map[string]interface{}{
				"event":           "delivered",
				"notification_id": notificationID,
				"reason":          "",
			},
			setupMock: func() {
				mockUseCase.EXPECT().
					UpdateStatus(gomock.Any(), notificationID, entity.NotificationStatusSent, "").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Success - Email Failed",
			provider: "email",
			payload: map[string]interface{}{
				"event":           "failed",
				"notification_id": notificationID,
				"reason":          "mailbox full",
			},
			setupMock: func() {
				mockUseCase.EXPECT().
					UpdateStatus(gomock.Any(), notificationID, entity.NotificationStatusFailed, "mailbox full").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Internal Server Error",
			provider: "email",
			payload: map[string]interface{}{
				"event":           "delivered",
				"notification_id": notificationID,
			},
			setupMock: func() {
				mockUseCase.EXPECT().
					UpdateStatus(gomock.Any(), notificationID, entity.NotificationStatusSent, "").
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:     "Bad Request - Invalid Payload",
			provider: "email",
			payload:  "invalid-json",
			setupMock: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.payload)
			if s, ok := tt.payload.(string); ok {
				body = []byte(s)
			}
			
			req, _ := http.NewRequest(http.MethodPost, "/webhooks/"+tt.provider, bytes.NewBuffer(body))
			c.Request = req
			c.Params = gin.Params{{Key: "provider", Value: tt.provider}}

			h.HandleWebhook(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestWebhookHandler_HandleWebhook_WhatsApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationUseCase(ctrl)
	h := handler.NewWebhookHandler(mockUseCase)

	notificationID := uuid.New()

	tests := []struct {
		name           string
		provider       string
		payload        interface{}
		setupMock      func()
		expectedStatus int
	}{
		{
			name:     "Success - WhatsApp Delivered",
			provider: "whatsapp",
			payload: map[string]interface{}{
				"status":          "delivered",
				"notification_id": notificationID,
				"error":           "",
			},
			setupMock: func() {
				mockUseCase.EXPECT().
					UpdateStatus(gomock.Any(), notificationID, entity.NotificationStatusSent, "").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Success - WhatsApp Failed",
			provider: "whatsapp",
			payload: map[string]interface{}{
				"status":          "failed",
				"notification_id": notificationID,
				"error":           "invalid number",
			},
			setupMock: func() {
				mockUseCase.EXPECT().
					UpdateStatus(gomock.Any(), notificationID, entity.NotificationStatusFailed, "invalid number").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Internal Server Error",
			provider: "whatsapp",
			payload: map[string]interface{}{
				"status":          "delivered",
				"notification_id": notificationID,
			},
			setupMock: func() {
				mockUseCase.EXPECT().
					UpdateStatus(gomock.Any(), notificationID, entity.NotificationStatusSent, "").
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/webhooks/"+tt.provider, bytes.NewBuffer(body))
			c.Request = req
			c.Params = gin.Params{{Key: "provider", Value: tt.provider}}

			h.HandleWebhook(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestWebhookHandler_HandleWebhook_UnknownProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUseCase := mocks.NewMockNotificationUseCase(ctrl)
	h := handler.NewWebhookHandler(mockUseCase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodPost, "/webhooks/unknown", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "provider", Value: "unknown"}}

	h.HandleWebhook(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
