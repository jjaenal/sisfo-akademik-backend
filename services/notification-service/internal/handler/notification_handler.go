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
