package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type TeacherAttendanceHandler struct {
	useCase usecase.TeacherAttendanceUseCase
}

func NewTeacherAttendanceHandler(useCase usecase.TeacherAttendanceUseCase) *TeacherAttendanceHandler {
	return &TeacherAttendanceHandler{useCase: useCase}
}

func (h *TeacherAttendanceHandler) CheckIn(c *gin.Context) {
	var req struct {
		TenantID       string                       `json:"tenant_id" binding:"required"`
		TeacherID      string                       `json:"teacher_id" binding:"required"`
		SemesterID     string                       `json:"semester_id" binding:"required"`
		AttendanceDate time.Time                    `json:"attendance_date" binding:"required"`
		CheckInTime    time.Time                    `json:"check_in_time" binding:"required"`
		Status         entity.TeacherAttendanceStatus `json:"status" binding:"required"`
		Notes          string                       `json:"notes"`
		LocationLatitude  *float64                  `json:"location_latitude"`
		LocationLongitude *float64                  `json:"location_longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Teacher ID")
		return
	}

	semesterID, err := uuid.Parse(req.SemesterID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Semester ID")
		return
	}

	attendance := &entity.TeacherAttendance{
		TenantID:          req.TenantID,
		TeacherID:         teacherID,
		SemesterID:        semesterID,
		AttendanceDate:    req.AttendanceDate,
		CheckInTime:       &req.CheckInTime,
		Status:            req.Status,
		Notes:             req.Notes,
		LocationLatitude:  req.LocationLatitude,
		LocationLongitude: req.LocationLongitude,
	}

	if err := h.useCase.CheckIn(c.Request.Context(), attendance); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, attendance)
}

func (h *TeacherAttendanceHandler) CheckOut(c *gin.Context) {
	var req struct {
		TeacherID    string    `json:"teacher_id" binding:"required"`
		Date         time.Time `json:"date" binding:"required"`
		CheckOutTime time.Time `json:"check_out_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Teacher ID")
		return
	}

	if err := h.useCase.CheckOut(c.Request.Context(), teacherID, req.Date, req.CheckOutTime); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Check-out successful"})
}

func (h *TeacherAttendanceHandler) GetByTeacherAndDate(c *gin.Context) {
	teacherIDStr := c.Query("teacher_id")
	dateStr := c.Query("date")

	if teacherIDStr == "" || dateStr == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Teacher ID and Date are required")
		return
	}

	teacherID, err := uuid.Parse(teacherIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Teacher ID")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Date Format (YYYY-MM-DD)")
		return
	}

	attendance, err := h.useCase.GetByTeacherAndDate(c.Request.Context(), teacherID, date)
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

func (h *TeacherAttendanceHandler) List(c *gin.Context) {
	filter := make(map[string]interface{})
	
	if teacherID := c.Query("teacher_id"); teacherID != "" {
		if id, err := uuid.Parse(teacherID); err == nil {
			filter["teacher_id"] = id
		}
	}
	
	if semesterID := c.Query("semester_id"); semesterID != "" {
		if id, err := uuid.Parse(semesterID); err == nil {
			filter["semester_id"] = id
		}
	}

	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	attendances, err := h.useCase.List(c.Request.Context(), filter)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, attendances)
}
