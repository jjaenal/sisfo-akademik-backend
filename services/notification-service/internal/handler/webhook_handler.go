package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type WebhookHandler struct {
	useCase usecase.NotificationUseCase
}

func NewWebhookHandler(useCase usecase.NotificationUseCase) *WebhookHandler {
	return &WebhookHandler{useCase: useCase}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	provider := c.Param("provider")

	switch provider {
	case "email":
		// Example: Mailgun/SendGrid webhook
		h.handleEmailWebhook(c)
	case "whatsapp":
		// Example: Twilio/WAbot webhook
		h.handleWhatsAppWebhook(c)
	default:
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Provider not found", "Unknown provider")
	}
}

func (h *WebhookHandler) handleEmailWebhook(c *gin.Context) {
	// Placeholder: Parse specific provider payload
	var payload struct {
		Event          string    `json:"event"`
		NotificationID uuid.UUID `json:"notification_id"` // Assuming we pass this in custom vars
		Reason         string    `json:"reason"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid payload", err.Error())
		return
	}

	status := entity.NotificationStatusSent
	if payload.Event == "failed" || payload.Event == "bounced" {
		status = entity.NotificationStatusFailed
	}

	if err := h.useCase.UpdateStatus(c.Request.Context(), payload.NotificationID, status, payload.Reason); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5000", "Failed to update status", err.Error())
		return
	}

	httputil.Success(c.Writer, "Webhook processed")
}

func (h *WebhookHandler) handleWhatsAppWebhook(c *gin.Context) {
	// Placeholder: Parse specific provider payload
	var payload struct {
		Status         string    `json:"status"`
		NotificationID uuid.UUID `json:"notification_id"`
		Error          string    `json:"error"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid payload", err.Error())
		return
	}

	status := entity.NotificationStatusSent
	if payload.Status == "failed" {
		status = entity.NotificationStatusFailed
	}

	if err := h.useCase.UpdateStatus(c.Request.Context(), payload.NotificationID, status, payload.Error); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5000", "Failed to update status", err.Error())
		return
	}

	httputil.Success(c.Writer, "Webhook processed")
}
