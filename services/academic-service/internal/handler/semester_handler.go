package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type SemesterHandler struct {
	useCase usecase.SemesterUseCase
}

func NewSemesterHandler(useCase usecase.SemesterUseCase) *SemesterHandler {
	return &SemesterHandler{useCase: useCase}
}

func (h *SemesterHandler) Create(c *gin.Context) {
	var req struct {
		TenantID       string    `json:"tenant_id" binding:"required"`
		AcademicYearID string    `json:"academic_year_id" binding:"required"`
		Name           string    `json:"name" binding:"required"`
		StartDate      time.Time `json:"start_date" binding:"required"`
		EndDate        time.Time `json:"end_date" binding:"required"`
		IsActive       bool      `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Academic Year ID", "Academic Year ID must be a valid UUID")
		return
	}

	semester := &entity.Semester{
		TenantID:       req.TenantID,
		AcademicYearID: academicYearID,
		Name:           req.Name,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsActive:       req.IsActive,
	}

	if err := h.useCase.Create(c.Request.Context(), semester); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, semester)
}

func (h *SemesterHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	semester, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if semester == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Semester not found")
		return
	}

	httputil.Success(c.Writer, semester)
}

func (h *SemesterHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	academicYearIDStr := c.Query("academic_year_id")

	var semesters []entity.Semester
	var err error

	if academicYearIDStr != "" {
		academicYearID, parseErr := uuid.Parse(academicYearIDStr)
		if parseErr != nil {
			httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Academic Year ID", "Academic Year ID must be a valid UUID")
			return
		}
		semesters, err = h.useCase.ListByAcademicYear(c.Request.Context(), academicYearID)
	} else if tenantID != "" {
		semesters, err = h.useCase.List(c.Request.Context(), tenantID)
	} else {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Either tenant_id or academic_year_id is required")
		return
	}

	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, semesters)
}

func (h *SemesterHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Name      string    `json:"name" binding:"required"`
		StartDate time.Time `json:"start_date" binding:"required"`
		EndDate   time.Time `json:"end_date" binding:"required"`
		IsActive  bool      `json:"is_active"`
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
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Semester not found")
		return
	}

	existing.Name = req.Name
	existing.StartDate = req.StartDate
	existing.EndDate = req.EndDate
	existing.IsActive = req.IsActive

	if err := h.useCase.Update(c.Request.Context(), existing); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

func (h *SemesterHandler) Activate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.SetActive(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{"message": "Semester activated successfully"})
}

func (h *SemesterHandler) Delete(c *gin.Context) {
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

	httputil.Success(c.Writer, map[string]string{"message": "Semester deleted successfully"})
}
