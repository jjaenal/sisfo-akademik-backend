package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
)

type NotificationHandler struct {
	useCase usecase.NotificationUseCase
}

func NewNotificationHandler(useCase usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{useCase: useCase}
}

// Send godoc
// @Summary Send a notification
// @Description Send a notification via Email or WhatsApp
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body usecase.SendNotificationRequest true "Notification Request"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/send [post]
func (h *NotificationHandler) Send(c *gin.Context) {
	var req usecase.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.Send(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "notification queued"})
}

// GetByID godoc
// @Summary Get notification by ID
// @Description Get notification details
// @Tags notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} entity.Notification
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/{id} [get]
func (h *NotificationHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	notification, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if notification == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notification})
}

// ListByRecipient godoc
// @Summary List notifications by recipient
// @Description Get all notifications sent to a recipient
// @Tags notifications
// @Produce json
// @Param recipient query string true "Recipient"
// @Success 200 {array} entity.Notification
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/recipient [get]
func (h *NotificationHandler) ListByRecipient(c *gin.Context) {
	recipient := c.Query("recipient")
	if recipient == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "recipient is required"})
		return
	}

	notifications, err := h.useCase.ListByRecipient(c.Request.Context(), recipient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notifications})
}
