package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type PaymentHandler struct {
	useCase usecase.PaymentUseCase
}

func NewPaymentHandler(useCase usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

func (h *PaymentHandler) Record(c *gin.Context) {
	var req struct {
		TenantID        uuid.UUID            `json:"tenant_id" binding:"required"`
		InvoiceID       uuid.UUID            `json:"invoice_id" binding:"required"`
		Amount          float64              `json:"amount" binding:"required,gt=0"`
		PaymentMethod   entity.PaymentMethod `json:"payment_method" binding:"required"`
		ReferenceNumber string               `json:"reference_number"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid input", err.Error())
		return
	}

	payment := &entity.Payment{
		TenantID:        req.TenantID,
		InvoiceID:       req.InvoiceID,
		Amount:          req.Amount,
		PaymentMethod:   req.PaymentMethod,
		ReferenceNumber: req.ReferenceNumber,
		Status:          entity.PaymentStatusSuccess, // Assuming direct success for now
	}

	if err := h.useCase.RecordPayment(c.Request.Context(), payment); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to record payment", err.Error())
		return
	}

	httputil.Success(c.Writer, payment)
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID format", err.Error())
		return
	}

	payment, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Payment not found", err.Error())
		return
	}

	httputil.Success(c.Writer, payment)
}

func (h *PaymentHandler) ListByInvoice(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("invoice_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid invoice_id format", err.Error())
		return
	}

	payments, err := h.useCase.ListByInvoiceID(c.Request.Context(), invoiceID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to list payments", err.Error())
		return
	}

	httputil.Success(c.Writer, payments)
}
