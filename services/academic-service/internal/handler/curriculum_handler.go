package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type CurriculumHandler struct {
	useCase usecase.CurriculumUseCase
}

func NewCurriculumHandler(useCase usecase.CurriculumUseCase) *CurriculumHandler {
	return &CurriculumHandler{useCase: useCase}
}

func (h *CurriculumHandler) Create(c *gin.Context) {
	var req struct {
		TenantID    string `json:"tenant_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Year        int    `json:"year" binding:"required"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	curriculum := &entity.Curriculum{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
		Year:        req.Year,
		IsActive:    req.IsActive,
	}

	if err := h.useCase.Create(c.Request.Context(), curriculum); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, curriculum)
}

func (h *CurriculumHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	curriculum, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if curriculum == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Curriculum not found")
		return
	}

	httputil.Success(c.Writer, curriculum)
}

func (h *CurriculumHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Tenant ID is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	curricula, total, err := h.useCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]interface{}{
		"curricula": curricula,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *CurriculumHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Year        int    `json:"year" binding:"required"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	existing, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if existing == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Curriculum not found")
		return
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Year = req.Year
	existing.IsActive = req.IsActive

	if err := h.useCase.Update(c.Request.Context(), existing); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

func (h *CurriculumHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Curriculum deleted successfully"})
}

// Curriculum Subjects Handlers

func (h *CurriculumHandler) AddSubject(c *gin.Context) {
	idStr := c.Param("id")
	curriculumID, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		TenantID   string    `json:"tenant_id" binding:"required"`
		SubjectID  uuid.UUID `json:"subject_id" binding:"required"`
		GradeLevel int       `json:"grade_level" binding:"required"`
		Semester   int       `json:"semester" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	cs := &entity.CurriculumSubject{
		TenantID:     req.TenantID,
		CurriculumID: curriculumID,
		SubjectID:    req.SubjectID,
		GradeLevel:   req.GradeLevel,
		Semester:     req.Semester,
	}

	if err := h.useCase.AddSubject(c.Request.Context(), cs); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, cs)
}

func (h *CurriculumHandler) ListSubjects(c *gin.Context) {
	idStr := c.Param("id")
	curriculumID, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	subjects, err := h.useCase.ListSubjects(c.Request.Context(), curriculumID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, subjects)
}

func (h *CurriculumHandler) RemoveSubject(c *gin.Context) {
	idStr := c.Param("subject_id") // curriculum_subject_id
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.RemoveSubject(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Subject removed from curriculum successfully"})
}
