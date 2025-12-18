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

// InputGrade godoc
// @Summary      Input a grade
// @Description  Input a grade for a student in an assessment
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        request body object true "Grade Input"
// @Success      200  {object}  entity.Grade
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grades [post]
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

// GetStudentGrades godoc
// @Summary      Get student grades
// @Description  Get grades for a student by class and semester
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        student_id   path      string  true  "Student ID"
// @Param        class_id     query     string  true  "Class ID"
// @Param        semester_id  query     string  true  "Semester ID"
// @Success      200  {array}   entity.Grade
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grades/student/{student_id} [get]
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

// CalculateFinalScore godoc
// @Summary      Calculate final score
// @Description  Calculate final score for a student in a subject
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        student_id   path      string  true  "Student ID"
// @Param        subject_id   query     string  true  "Subject ID"
// @Param        semester_id  query     string  true  "Semester ID"
// @Param        class_id     query     string  true  "Class ID"
// @Success      200  {object}  map[string]float64
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grades/calculate/{student_id} [get]
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

// ApproveGrade godoc
// @Summary      Approve a grade
// @Description  Approve a grade
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        id           path      string  true  "Grade ID"
// @Param        request      body      object  true  "Approve Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grades/{id}/approve [put]
func (h *GradeHandler) ApproveGrade(c *gin.Context) {
	gradeIDStr := c.Param("id")
	gradeID, err := uuid.Parse(gradeIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Grade ID")
		return
	}

	var req struct {
		ApprovedBy string `json:"approved_by" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	approvedBy, err := uuid.Parse(req.ApprovedBy)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Approver ID")
		return
	}

	if err := h.useCase.ApproveGrade(c.Request.Context(), gradeID, approvedBy); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Grade approved successfully"})
}
