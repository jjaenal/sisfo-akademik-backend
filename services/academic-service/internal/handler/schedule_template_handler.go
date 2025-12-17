package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type ScheduleTemplateHandler struct {
	useCase usecase.ScheduleTemplateUseCase
}

func NewScheduleTemplateHandler(useCase usecase.ScheduleTemplateUseCase) *ScheduleTemplateHandler {
	return &ScheduleTemplateHandler{useCase: useCase}
}

func (h *ScheduleTemplateHandler) Create(c *gin.Context) {
	var req struct {
		TenantID    string `json:"tenant_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	template := &entity.ScheduleTemplate{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
	}

	if err := h.useCase.Create(c.Request.Context(), template); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, template)
}

func (h *ScheduleTemplateHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	template, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if template == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Template not found")
		return
	}

	httputil.Success(c.Writer, template)
}

func (h *ScheduleTemplateHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Tenant ID is required")
		return
	}

	templates, err := h.useCase.List(c.Request.Context(), tenantID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, templates)
}

func (h *ScheduleTemplateHandler) AddItem(c *gin.Context) {
	idStr := c.Param("id")
	templateID, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		SubjectID *uuid.UUID `json:"subject_id"`
		DayOfWeek int        `json:"day_of_week" binding:"required"`
		StartTime string     `json:"start_time" binding:"required"`
		EndTime   string     `json:"end_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	item := &entity.ScheduleTemplateItem{
		TemplateID: templateID,
		SubjectID:  req.SubjectID,
		DayOfWeek:  req.DayOfWeek,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}

	if err := h.useCase.AddItem(c.Request.Context(), item); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, item)
}

func (h *ScheduleTemplateHandler) RemoveItem(c *gin.Context) {
	// The route should probably be DELETE /items/:item_id
	// But usually we nest it under template: DELETE /templates/:id/items/:item_id
	// Or just DELETE /templates/items/:item_id
	
	idStr := c.Param("item_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.RemoveItem(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Item removed successfully"})
}
