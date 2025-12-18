package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type StudentAttendanceHandler struct {
	useCase usecase.StudentAttendanceUseCase
}

func NewStudentAttendanceHandler(useCase usecase.StudentAttendanceUseCase) *StudentAttendanceHandler {
	return &StudentAttendanceHandler{useCase: useCase}
}

// Create godoc
// @Summary Create student attendance
// @Description Record attendance for a student
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param request body map[string]any true "Attendance Request"
// @Success 200 {object} entity.StudentAttendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/students [post]
func (h *StudentAttendanceHandler) Create(c *gin.Context) {
	var req struct {
		TenantID         string                  `json:"tenant_id" binding:"required"`
		StudentID        string                  `json:"student_id" binding:"required"`
		ClassID          string                  `json:"class_id" binding:"required"`
		SemesterID       string                  `json:"semester_id" binding:"required"`
		AttendanceDate   time.Time               `json:"attendance_date" binding:"required"`
		Status           entity.AttendanceStatus `json:"status" binding:"required"`
		Notes            string                  `json:"notes"`
		CheckInLatitude  *float64                `json:"check_in_latitude"`
		CheckInLongitude *float64                `json:"check_in_longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

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

	attendance := &entity.StudentAttendance{
		TenantID:         req.TenantID,
		StudentID:        studentID,
		ClassID:          classID,
		SemesterID:       semesterID,
		AttendanceDate:   req.AttendanceDate,
		Status:           req.Status,
		Notes:            req.Notes,
		CheckInLatitude:  req.CheckInLatitude,
		CheckInLongitude: req.CheckInLongitude,
	}

	if err := h.useCase.Create(c.Request.Context(), attendance); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, attendance)
}

// GetDailyReport godoc
// @Summary Get daily attendance report
// @Description Get attendance report for a specific date
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Param tenant_id query string true "Tenant ID"
// @Success 200 {object} []entity.StudentAttendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/reports/daily [get]
func (h *StudentAttendanceHandler) GetDailyReport(c *gin.Context) {
	dateStr := c.Query("date")
	tenantIDStr := c.Query("tenant_id")

	if dateStr == "" || tenantIDStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Date and Tenant ID are required")
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Tenant ID")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Date format (expected YYYY-MM-DD)")
		return
	}

	report, err := h.useCase.GetDailyReport(c.Request.Context(), tenantID, date)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, report)
}

// GetMonthlyReport godoc
// @Summary Get monthly attendance report
// @Description Get attendance report for a specific month
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param month query int true "Month (1-12)"
// @Param year query int true "Year (YYYY)"
// @Param tenant_id query string true "Tenant ID"
// @Success 200 {object} []entity.StudentAttendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/reports/monthly [get]
func (h *StudentAttendanceHandler) GetMonthlyReport(c *gin.Context) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	tenantIDStr := c.Query("tenant_id")

	if monthStr == "" || yearStr == "" || tenantIDStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Month, Year and Tenant ID are required")
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Tenant ID")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Month")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Year")
		return
	}

	report, err := h.useCase.GetMonthlyReport(c.Request.Context(), tenantID, month, year)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, report)
}

// GetClassReport godoc
// @Summary Get class attendance report
// @Description Get attendance report for a specific class and date range
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Param start_date query string true "Start Date (YYYY-MM-DD)"
// @Param end_date query string true "End Date (YYYY-MM-DD)"
// @Param tenant_id query string true "Tenant ID"
// @Success 200 {object} []entity.StudentAttendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/reports/class/{class_id} [get]
func (h *StudentAttendanceHandler) GetClassReport(c *gin.Context) {
	classIDStr := c.Param("class_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	tenantIDStr := c.Query("tenant_id")

	if classIDStr == "" || startDateStr == "" || endDateStr == "" || tenantIDStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Class ID, Start Date, End Date and Tenant ID are required")
		return
	}

	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Class ID")
		return
	}

	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Tenant ID")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Start Date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid End Date format")
		return
	}

	report, err := h.useCase.GetClassReport(c.Request.Context(), tenantID, classID, startDate, endDate)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, report)
}

// BulkCreate godoc
// @Summary Bulk create student attendance
// @Description Record attendance for multiple students
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param request body []map[string]any true "Bulk Attendance Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/students/bulk [post]
func (h *StudentAttendanceHandler) BulkCreate(c *gin.Context) {
	var reqs []struct {
		TenantID         string                  `json:"tenant_id" binding:"required"`
		StudentID        string                  `json:"student_id" binding:"required"`
		ClassID          string                  `json:"class_id" binding:"required"`
		SemesterID       string                  `json:"semester_id" binding:"required"`
		AttendanceDate   time.Time               `json:"attendance_date" binding:"required"`
		Status           entity.AttendanceStatus `json:"status" binding:"required"`
		Notes            string                  `json:"notes"`
		CheckInLatitude  *float64                `json:"check_in_latitude"`
		CheckInLongitude *float64                `json:"check_in_longitude"`
	}

	if err := c.ShouldBindJSON(&reqs); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	var attendances []*entity.StudentAttendance
	for _, req := range reqs {
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

		attendance := &entity.StudentAttendance{
			TenantID:         req.TenantID,
			StudentID:        studentID,
			ClassID:          classID,
			SemesterID:       semesterID,
			AttendanceDate:   req.AttendanceDate,
			Status:           req.Status,
			Notes:            req.Notes,
			CheckInLatitude:  req.CheckInLatitude,
			CheckInLongitude: req.CheckInLongitude,
		}
		attendances = append(attendances, attendance)
	}

	if err := h.useCase.BulkCreate(c.Request.Context(), attendances); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, attendances)
}

// GetSummary godoc
// @Summary Get student attendance summary
// @Description Get attendance summary for a student in a semester
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param semester_id query string false "Semester ID"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/students/{id}/summary [get]
func (h *StudentAttendanceHandler) GetSummary(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Student ID", "ID must be a valid UUID")
		return
	}

	semesterIDStr := c.Query("semester_id")
	var semesterID uuid.UUID
	if semesterIDStr != "" {
		semesterID, err = uuid.Parse(semesterIDStr)
		if err != nil {
			httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Semester ID", "ID must be a valid UUID")
			return
		}
	}

	summary, err := h.useCase.GetSummary(c.Request.Context(), studentID, semesterID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, summary)
}

// GetByID godoc
// @Summary Get student attendance by ID
// @Description Get detailed attendance record by ID
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param id path string true "Attendance ID"
// @Success 200 {object} entity.StudentAttendance
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/students/{id} [get]
func (h *StudentAttendanceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	attendance, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if attendance == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Attendance record not found")
		return
	}

	httputil.Success(c.Writer, attendance)
}

// GetByClassAndDate godoc
// @Summary Get student attendance by class and date
// @Description Get list of attendance records for a class on a specific date
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param class_id query string true "Class ID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} []entity.StudentAttendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/students [get]
func (h *StudentAttendanceHandler) GetByClassAndDate(c *gin.Context) {
	classIDStr := c.Query("class_id")
	dateStr := c.Query("date")

	if classIDStr == "" || dateStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Class ID and Date are required")
		return
	}

	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Class ID")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Date format (expected YYYY-MM-DD)")
		return
	}

	attendances, err := h.useCase.GetByClassAndDate(c.Request.Context(), classID, date)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, attendances)
}

// Update godoc
// @Summary Update student attendance
// @Description Update attendance status and notes
// @Tags student-attendance
// @Accept json
// @Produce json
// @Param id path string true "Attendance ID"
// @Param request body map[string]any true "Update Request"
// @Success 200 {object} entity.StudentAttendance
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /attendance/students/{id} [put]
func (h *StudentAttendanceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Status entity.AttendanceStatus `json:"status" binding:"required"`
		Notes  string                  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	attendance, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if attendance == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Attendance record not found")
		return
	}

	attendance.Status = req.Status
	attendance.Notes = req.Notes

	if err := h.useCase.Update(c.Request.Context(), attendance); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, attendance)
}
