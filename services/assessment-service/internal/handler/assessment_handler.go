package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type AssessmentHandler struct {
	useCase usecase.GradingUseCase
}

func NewAssessmentHandler(useCase usecase.GradingUseCase) *AssessmentHandler {
	return &AssessmentHandler{useCase: useCase}
}

func (h *AssessmentHandler) Create(c *gin.Context) {
	var req struct {
		SubjectID       string    `json:"subject_id" binding:"required"`
		TeacherID       string    `json:"teacher_id" binding:"required"`
		ClassID         string    `json:"class_id" binding:"required"`
		GradeCategoryID string    `json:"grade_category_id" binding:"required"`
		Name            string    `json:"name" binding:"required"`
		MaxScore        int       `json:"max_score" binding:"required,min=0"`
		Description     string    `json:"description"`
		Date            time.Time `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Subject ID")
		return
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Teacher ID")
		return
	}

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Class ID")
		return
	}

	gradeCategoryID, err := uuid.Parse(req.GradeCategoryID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Grade Category ID")
		return
	}

	assessment := &entity.Assessment{
		SubjectID:       subjectID,
		TeacherID:       teacherID,
		ClassID:         classID,
		GradeCategoryID: gradeCategoryID,
		Name:            req.Name,
		MaxScore:        req.MaxScore,
		Description:     req.Description,
		Date:            req.Date,
	}

	if err := h.useCase.CreateAssessment(c.Request.Context(), assessment); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, assessment)
}
