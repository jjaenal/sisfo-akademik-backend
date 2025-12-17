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

type StudentAttendanceHandler struct {
	useCase usecase.StudentAttendanceUseCase
}

func NewStudentAttendanceHandler(useCase usecase.StudentAttendanceUseCase) *StudentAttendanceHandler {
	return &StudentAttendanceHandler{useCase: useCase}
}

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
