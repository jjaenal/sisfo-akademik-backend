package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/notification-service/internal/usecase"
)

type NotificationTemplateHandler struct {
	useCase usecase.NotificationTemplateUseCase
}

func NewNotificationTemplateHandler(useCase usecase.NotificationTemplateUseCase) *NotificationTemplateHandler {
	return &NotificationTemplateHandler{useCase: useCase}
}

// Create godoc
// @Summary Create a notification template
// @Description Create a new notification template
// @Tags templates
// @Accept json
// @Produce json
// @Param request body entity.NotificationTemplate true "Template"
// @Success 201 {object} entity.NotificationTemplate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/templates [post]
func (h *NotificationTemplateHandler) Create(c *gin.Context) {
	var req entity.NotificationTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.Create(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": req})
}

// Update godoc
// @Summary Update a notification template
// @Description Update an existing notification template
// @Tags templates
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param request body entity.NotificationTemplate true "Template"
// @Success 200 {object} entity.NotificationTemplate
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/templates/{id} [put]
func (h *NotificationTemplateHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req entity.NotificationTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id

	if err := h.useCase.Update(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": req})
}

// GetByID godoc
// @Summary Get template by ID
// @Description Get template details
// @Tags templates
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} entity.NotificationTemplate
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/templates/{id} [get]
func (h *NotificationTemplateHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	template, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": template})
}

// List godoc
// @Summary List templates
// @Description Get all notification templates
// @Tags templates
// @Produce json
// @Success 200 {array} entity.NotificationTemplate
// @Failure 500 {object} map[string]string
// @Router /notifications/templates [get]
func (h *NotificationTemplateHandler) List(c *gin.Context) {
	templates, err := h.useCase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": templates})
}

// Delete godoc
// @Summary Delete template
// @Description Delete a notification template
// @Tags templates
// @Produce json
// @Param id path string true "Template ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /notifications/templates/{id} [delete]
func (h *NotificationTemplateHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "template deleted"})
}
