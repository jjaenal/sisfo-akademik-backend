package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type ReportCardHandler struct {
	useCase usecase.ReportCardUseCase
}

func NewReportCardHandler(useCase usecase.ReportCardUseCase) *ReportCardHandler {
	return &ReportCardHandler{useCase: useCase}
}

// Generate godoc
// @Summary      Generate a report card
// @Description  Generate a report card for a student
// @Tags         report-cards
// @Accept       json
// @Produce      json
// @Param        request body object true "Generate Request"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /report-cards/generate [post]
func (h *ReportCardHandler) Generate(c *gin.Context) {
	var req struct {
		TenantID   string `json:"tenant_id" binding:"required"`
		StudentID  string `json:"student_id" binding:"required"`
		ClassID    string `json:"class_id" binding:"required"`
		SemesterID string `json:"semester_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	tenantID := req.TenantID
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Student ID")
		return
	}
	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Class ID")
		return
	}
	semesterID, err := uuid.Parse(req.SemesterID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Semester ID")
		return
	}

	rc, err := h.useCase.Generate(c.Request.Context(), tenantID, studentID, classID, semesterID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, rc)
}

// GetPDF godoc
// @Summary      Download report card PDF
// @Description  Download the generated report card PDF
// @Tags         report-cards
// @Produce      application/pdf
// @Param        id   path      string  true  "Report Card ID"
// @Success      200  {file}    file
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /report-cards/{id}/download [get]
func (h *ReportCardHandler) GetPDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid ID")
		return
	}

	pdfBytes, err := h.useCase.GetPDF(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	c.Header("Content-Disposition", "attachment; filename=report_card.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetByStudent godoc
// @Summary      Get report card by student
// @Description  Get report card for a student in a semester
// @Tags         report-cards
// @Accept       json
// @Produce      json
// @Param        student_id   path      string  true  "Student ID"
// @Param        semester_id  query     string  true  "Semester ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      404  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /report-cards/student/{student_id} [get]
func (h *ReportCardHandler) GetByStudent(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	semesterIDStr := c.Query("semester_id")

	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Student ID")
		return
	}
	semesterID, err := uuid.Parse(semesterIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Semester ID")
		return
	}

	rc, err := h.useCase.GetByStudent(c.Request.Context(), studentID, semesterID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if rc == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Report Card not found")
		return
	}

	httputil.Success(c.Writer, rc)
}

func (h *ReportCardHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid ID")
		return
	}

	rc, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if rc == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Report Card not found")
		return
	}

	httputil.Success(c.Writer, rc)
}

func (h *ReportCardHandler) Publish(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid ID")
		return
	}

	err = h.useCase.Publish(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Report Card published"})
}
