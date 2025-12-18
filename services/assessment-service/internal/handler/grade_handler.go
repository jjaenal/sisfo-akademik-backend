package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type GradeHandler struct {
	useCase usecase.GradingUseCase
}

func NewGradeHandler(useCase usecase.GradingUseCase) *GradeHandler {
	return &GradeHandler{useCase: useCase}
}

func (h *GradeHandler) InputGrade(c *gin.Context) {
	var req struct {
		TenantID     string  `json:"tenant_id" binding:"required"`
		AssessmentID string  `json:"assessment_id" binding:"required"`
		StudentID    string  `json:"student_id" binding:"required"`
		Score        float64 `json:"score" binding:"required,min=0"`
		Notes        string  `json:"notes"`
		GradedBy     string  `json:"graded_by" binding:"required"`
		Status       string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	assessmentID, err := uuid.Parse(req.AssessmentID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Assessment ID")
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Student ID")
		return
	}

	gradedBy, err := uuid.Parse(req.GradedBy)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Grader ID")
		return
	}

	status := entity.GradeStatusDraft
	if req.Status != "" {
		status = entity.GradeStatus(req.Status)
	}

	grade := &entity.Grade{
		TenantID:     req.TenantID,
		AssessmentID: assessmentID,
		StudentID:    studentID,
		Score:        req.Score,
		Notes:        req.Notes,
		GradedBy:     gradedBy,
		Status:       status,
	}

	if err := h.useCase.InputGrade(c.Request.Context(), grade); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, grade)
}

func (h *GradeHandler) GetStudentGrades(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	classIDStr := c.Query("class_id")
	semesterIDStr := c.Query("semester_id")

	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Student ID")
		return
	}

	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Class ID")
		return
	}

	semesterID, err := uuid.Parse(semesterIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Semester ID")
		return
	}

	grades, err := h.useCase.GetStudentGrades(c.Request.Context(), studentID, classID, semesterID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, grades)
}

func (h *GradeHandler) CalculateFinalScore(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	subjectIDStr := c.Query("subject_id")
	semesterIDStr := c.Query("semester_id")
	classIDStr := c.Query("class_id")

	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Student ID")
		return
	}

	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Subject ID")
		return
	}

	semesterID, err := uuid.Parse(semesterIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Semester ID")
		return
	}

	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Class ID")
		return
	}

	score, err := h.useCase.CalculateFinalScore(c.Request.Context(), studentID, classID, subjectID, semesterID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"final_score": score})
}
