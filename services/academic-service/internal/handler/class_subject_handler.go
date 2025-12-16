package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type ClassSubjectHandler struct {
	useCase usecase.ClassSubjectUseCase
}

func NewClassSubjectHandler(useCase usecase.ClassSubjectUseCase) *ClassSubjectHandler {
	return &ClassSubjectHandler{useCase: useCase}
}

func (h *ClassSubjectHandler) AddSubjectToClass(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Class ID", "Class ID must be a valid UUID")
		return
	}

	var req struct {
		TenantID  string  `json:"tenant_id" binding:"required"`
		SubjectID string  `json:"subject_id" binding:"required"`
		TeacherID *string `json:"teacher_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Subject ID", "Subject ID must be a valid UUID")
		return
	}

	var teacherID *uuid.UUID
	if req.TeacherID != nil && *req.TeacherID != "" {
		tid, err := uuid.Parse(*req.TeacherID)
		if err != nil {
			httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Teacher ID", "Teacher ID must be a valid UUID")
			return
		}
		teacherID = &tid
	}

	classSubject := &entity.ClassSubject{
		TenantID:  req.TenantID,
		ClassID:   classID,
		SubjectID: subjectID,
		TeacherID: teacherID,
	}

	if err := h.useCase.Create(c.Request.Context(), classSubject); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, classSubject)
}

func (h *ClassSubjectHandler) AssignTeacher(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Class ID", "Class ID must be a valid UUID")
		return
	}

	subjectIDStr := c.Param("subject_id")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Subject ID", "Subject ID must be a valid UUID")
		return
	}

	var req struct {
		TeacherID string `json:"teacher_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Teacher ID", "Teacher ID must be a valid UUID")
		return
	}

	// First find the ClassSubject record
	cs, err := h.useCase.GetByClassAndSubject(c.Request.Context(), classID, subjectID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if cs == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Subject not assigned to this class")
		return
	}

	if err := h.useCase.AssignTeacher(c.Request.Context(), cs.ID, teacherID); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	// Re-fetch to return updated data
	updated, err := h.useCase.GetByID(c.Request.Context(), cs.ID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, updated)
}

func (h *ClassSubjectHandler) ListByClass(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Class ID", "Class ID must be a valid UUID")
		return
	}

	subjects, err := h.useCase.ListByClass(c.Request.Context(), classID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, subjects)
}

func (h *ClassSubjectHandler) RemoveSubject(c *gin.Context) {
	classIDStr := c.Param("id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Class ID", "Class ID must be a valid UUID")
		return
	}

	subjectIDStr := c.Param("subject_id")
	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Subject ID", "Subject ID must be a valid UUID")
		return
	}

	// Find record
	cs, err := h.useCase.GetByClassAndSubject(c.Request.Context(), classID, subjectID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if cs == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Subject not assigned to this class")
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), cs.ID); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Subject removed from class successfully"})
}
