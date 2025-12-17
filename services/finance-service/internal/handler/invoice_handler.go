package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type InvoiceHandler struct {
	useCase usecase.InvoiceUseCase
}

func NewInvoiceHandler(useCase usecase.InvoiceUseCase) *InvoiceHandler {
	return &InvoiceHandler{useCase: useCase}
}

func (h *InvoiceHandler) Generate(c *gin.Context) {
	var req struct {
		TenantID        uuid.UUID `json:"tenant_id" binding:"required"`
		StudentID       uuid.UUID `json:"student_id" binding:"required"`
		BillingConfigID uuid.UUID `json:"billing_config_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid input", err.Error())
		return
	}

	invoice, err := h.useCase.Generate(c.Request.Context(), req.TenantID, req.StudentID, req.BillingConfigID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to generate invoice", err.Error())
		return
	}

	httputil.Success(c.Writer, invoice)
}

func (h *InvoiceHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID format", err.Error())
		return
	}

	invoice, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Invoice not found", err.Error())
		return
	}

	httputil.Success(c.Writer, invoice)
}

func (h *InvoiceHandler) List(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid tenant_id format", err.Error())
		return
	}

	var studentID uuid.UUID
	if sID := c.Query("student_id"); sID != "" {
		parsed, err := uuid.Parse(sID)
		if err != nil {
			httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid student_id format", err.Error())
			return
		}
		studentID = parsed
	}

	status := entity.InvoiceStatus(c.Query("status"))

	invoices, err := h.useCase.List(c.Request.Context(), tenantID, studentID, status)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to list invoices", err.Error())
		return
	}

	httputil.Success(c.Writer, invoices)
}
