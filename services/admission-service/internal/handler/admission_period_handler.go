package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type AdmissionPeriodHandler struct {
	useCase usecase.AdmissionPeriodUseCase
}

func NewAdmissionPeriodHandler(useCase usecase.AdmissionPeriodUseCase) *AdmissionPeriodHandler {
	return &AdmissionPeriodHandler{useCase: useCase}
}

func (h *AdmissionPeriodHandler) Create(c *gin.Context) {
	var req struct {
		Name      string    `json:"name" binding:"required"`
		StartDate time.Time `json:"start_date" binding:"required"`
		EndDate   time.Time `json:"end_date" binding:"required"`
		IsActive  bool      `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	period := &entity.AdmissionPeriod{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  req.IsActive,
	}

	if err := h.useCase.Create(c.Request.Context(), period); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, period)
}

func (h *AdmissionPeriodHandler) AnnounceResults(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		PassingGrade float64 `json:"passing_grade" binding:"required,min=0,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	if err := h.useCase.AnnounceResults(c.Request.Context(), id, req.PassingGrade); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Admission results announced"})
}

func (h *AdmissionPeriodHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	period, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if period == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Admission period not found")
		return
	}

	httputil.Success(c.Writer, period)
}

func (h *AdmissionPeriodHandler) GetActive(c *gin.Context) {
	period, err := h.useCase.GetActive(c.Request.Context())
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if period == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "No active admission period found")
		return
	}

	httputil.Success(c.Writer, period)
}

func (h *AdmissionPeriodHandler) List(c *gin.Context) {
	periods, err := h.useCase.List(c.Request.Context())
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, periods)
}

func (h *AdmissionPeriodHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Name      string    `json:"name" binding:"required"`
		StartDate time.Time `json:"start_date" binding:"required"`
		EndDate   time.Time `json:"end_date" binding:"required"`
		IsActive  bool      `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	period, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if period == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Admission period not found")
		return
	}

	period.Name = req.Name
	period.StartDate = req.StartDate
	period.EndDate = req.EndDate
	period.IsActive = req.IsActive

	if err := h.useCase.Update(c.Request.Context(), period); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, period)
}

func (h *AdmissionPeriodHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Admission period deleted successfully"})
}
