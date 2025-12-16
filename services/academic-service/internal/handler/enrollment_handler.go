package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type EnrollmentHandler struct {
	useCase usecase.EnrollmentUseCase
}

func NewEnrollmentHandler(useCase usecase.EnrollmentUseCase) *EnrollmentHandler {
	return &EnrollmentHandler{useCase: useCase}
}

func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	var req struct {
		TenantID  string    `json:"tenant_id" binding:"required"`
		ClassID   uuid.UUID `json:"class_id" binding:"required"`
		StudentID uuid.UUID `json:"student_id" binding:"required"`
		Status    string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	enrollment := &entity.Enrollment{
		TenantID:  req.TenantID,
		ClassID:   req.ClassID,
		StudentID: req.StudentID,
		Status:    req.Status,
	}

	if err := h.useCase.Enroll(c.Request.Context(), enrollment); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, enrollment)
}

func (h *EnrollmentHandler) Unenroll(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.Unenroll(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Student unenrolled successfully"})
}

func (h *EnrollmentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	enrollment, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if enrollment == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Enrollment not found")
		return
	}

	httputil.Success(c.Writer, enrollment)
}

func (h *EnrollmentHandler) ListByClass(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "Class ID must be a valid UUID")
		return
	}

	enrollments, err := h.useCase.ListByClass(c.Request.Context(), classID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, enrollments)
}

func (h *EnrollmentHandler) ListByStudent(c *gin.Context) {
	studentIDStr := c.Param("id")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "Student ID must be a valid UUID")
		return
	}

	enrollments, err := h.useCase.ListByStudent(c.Request.Context(), studentID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, enrollments)
}

func (h *EnrollmentHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	if err := h.useCase.UpdateStatus(c.Request.Context(), id, req.Status); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Enrollment status updated successfully"})
}
