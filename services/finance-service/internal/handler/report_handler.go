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

// GetDailyRevenue godoc
// @Summary      Get daily revenue
// @Description  Get daily revenue report for a date range
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        tenant_id  query     string  true  "Tenant ID"
// @Param        start_date query     string  true  "Start Date (YYYY-MM-DD)"
// @Param        end_date   query     string  true  "End Date (YYYY-MM-DD)"
// @Success      200  {object}  []map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /finance/reports/revenue/daily [get]
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

// GetMonthlyRevenue godoc
// @Summary      Get monthly revenue
// @Description  Get monthly revenue report for a year
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        tenant_id  query     string  true  "Tenant ID"
// @Param        year       query     int     true  "Year"
// @Success      200  {object}  []map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /finance/reports/revenue/monthly [get]
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

// GetOutstandingInvoices godoc
// @Summary      Get outstanding invoices
// @Description  Get list of outstanding invoices
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        tenant_id  query     string  true  "Tenant ID"
// @Success      200  {object}  []map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /finance/reports/outstanding [get]
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

// GetStudentHistory godoc
// @Summary      Get student payment history
// @Description  Get payment history for a student
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        student_id path      string  true  "Student ID"
// @Param        tenant_id  query     string  true  "Tenant ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /finance/reports/student/{student_id}/history [get]
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
