package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/finance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type ReportHandler struct {
	useCase usecase.ReportUseCase
}

func NewReportHandler(useCase usecase.ReportUseCase) *ReportHandler {
	return &ReportHandler{useCase: useCase}
}

func (h *ReportHandler) GetDailyRevenue(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid tenant_id format", err.Error())
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4002", "Required parameters missing", "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid start_date format", "Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid end_date format", "Use YYYY-MM-DD")
		return
	}

	reports, err := h.useCase.GetDailyRevenue(c.Request.Context(), tenantID, startDate, endDate)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to get daily revenue", err.Error())
		return
	}

	httputil.Success(c.Writer, reports)
}

func (h *ReportHandler) GetMonthlyRevenue(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid tenant_id format", err.Error())
		return
	}

	yearStr := c.Query("year")
	if yearStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4002", "Required parameters missing", "year is required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid year format", "Must be a number")
		return
	}

	reports, err := h.useCase.GetMonthlyRevenue(c.Request.Context(), tenantID, year)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to get monthly revenue", err.Error())
		return
	}

	httputil.Success(c.Writer, reports)
}

func (h *ReportHandler) GetOutstandingInvoices(c *gin.Context) {
	tenantID, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid tenant_id format", err.Error())
		return
	}

	reports, err := h.useCase.GetOutstandingInvoices(c.Request.Context(), tenantID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to get outstanding invoices", err.Error())
		return
	}

	httputil.Success(c.Writer, reports)
}

func (h *ReportHandler) GetStudentHistory(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid student_id format", err.Error())
		return
	}

	tenantID, err := uuid.Parse(c.Query("tenant_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid tenant_id format", err.Error())
		return
	}

	history, err := h.useCase.GetStudentHistory(c.Request.Context(), tenantID, studentID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Failed to get student history", err.Error())
		return
	}

	httputil.Success(c.Writer, history)
}
