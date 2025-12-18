package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type TemplateHandler struct {
	useCase usecase.TemplateUseCase
}

func NewTemplateHandler(useCase usecase.TemplateUseCase) *TemplateHandler {
	return &TemplateHandler{useCase: useCase}
}

func (h *TemplateHandler) Create(c *gin.Context) {
	var template entity.ReportCardTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid request body", err.Error())
		return
	}

	// Assume tenant ID comes from middleware/context (mocked for now)
	// In real implementation: tenantID := c.GetString("tenant_id")
	if template.TenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4002", "tenant_id is required", nil)
		return
	}

	if err := h.useCase.Create(c.Request.Context(), &template); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to create template", err.Error())
		return
	}

	httputil.Success(c.Writer, template)
}

func (h *TemplateHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid template ID", err.Error())
		return
	}

	template, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to fetch template", err.Error())
		return
	}
	if template == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Template not found", nil)
		return
	}

	httputil.Success(c.Writer, template)
}

func (h *TemplateHandler) List(c *gin.Context) {
	tenantIDStr := c.Query("tenant_id")
	if tenantIDStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4002", "tenant_id is required", nil)
		return
	}

	templates, err := h.useCase.GetByTenantID(c.Request.Context(), tenantIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to fetch templates", err.Error())
		return
	}

	httputil.Success(c.Writer, templates)
}

func (h *TemplateHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid template ID", err.Error())
		return
	}

	var template entity.ReportCardTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid request body", err.Error())
		return
	}
	template.ID = id

	if err := h.useCase.Update(c.Request.Context(), &template); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to update template", err.Error())
		return
	}

	httputil.Success(c.Writer, template)
}

func (h *TemplateHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid template ID", err.Error())
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to delete template", err.Error())
		return
	}

	httputil.Success(c.Writer, nil)
}
