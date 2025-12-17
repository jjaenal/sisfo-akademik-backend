package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type BillingConfigHandler struct {
	useCase usecase.BillingConfigUseCase
}

func NewBillingConfigHandler(useCase usecase.BillingConfigUseCase) *BillingConfigHandler {
	return &BillingConfigHandler{useCase: useCase}
}

func (h *BillingConfigHandler) Create(c *gin.Context) {
	var req struct {
		TenantID  uuid.UUID               `json:"tenant_id" binding:"required"`
		Name      string                  `json:"name" binding:"required"`
		Amount    float64                 `json:"amount" binding:"required,gt=0"`
		Frequency entity.BillingFrequency `json:"frequency" binding:"required,oneof=MONTHLY ONCE YEARLY"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid input", err.Error())
		return
	}

	config := &entity.BillingConfig{
		TenantID:  req.TenantID,
		Name:      req.Name,
		Amount:    req.Amount,
		Frequency: req.Frequency,
		IsActive:  true,
	}

	if err := h.useCase.Create(c.Request.Context(), config); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to create billing config", err.Error())
		return
	}

	httputil.Success(c.Writer, config)
}

func (h *BillingConfigHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID format", err.Error())
		return
	}

	var req struct {
		Name      string                  `json:"name" binding:"required"`
		Amount    float64                 `json:"amount" binding:"required,gt=0"`
		Frequency entity.BillingFrequency `json:"frequency" binding:"required,oneof=MONTHLY ONCE YEARLY"`
		IsActive  bool                    `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid input", err.Error())
		return
	}

	config := &entity.BillingConfig{
		ID:        id,
		Name:      req.Name,
		Amount:    req.Amount,
		Frequency: req.Frequency,
		IsActive:  req.IsActive,
	}

	if err := h.useCase.Update(c.Request.Context(), config); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to update billing config", err.Error())
		return
	}

	httputil.Success(c.Writer, config)
}

func (h *BillingConfigHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID format", err.Error())
		return
	}

	config, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Billing config not found", err.Error())
		return
	}

	httputil.Success(c.Writer, config)
}

func (h *BillingConfigHandler) List(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid tenant_id format", err.Error())
		return
	}

	configs, err := h.useCase.List(c.Request.Context(), tenantID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to list billing configs", err.Error())
		return
	}

	httputil.Success(c.Writer, configs)
}

func (h *BillingConfigHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID format", err.Error())
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to delete billing config", err.Error())
		return
	}

	httputil.Success(c.Writer, nil)
}
