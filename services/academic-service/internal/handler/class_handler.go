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

type ClassHandler struct {
	useCase usecase.ClassUseCase
}

func NewClassHandler(useCase usecase.ClassUseCase) *ClassHandler {
	return &ClassHandler{useCase: useCase}
}

func (h *ClassHandler) Create(c *gin.Context) {
	var req struct {
		TenantID          string     `json:"tenant_id" binding:"required"`
		SchoolID          *uuid.UUID `json:"school_id"`
		AcademicYearID    *uuid.UUID `json:"academic_year_id"`
		Name              string     `json:"name" binding:"required"`
		Level             int        `json:"level"`
		Major             string     `json:"major"`
		HomeroomTeacherID *uuid.UUID `json:"homeroom_teacher_id"`
		Capacity          int        `json:"capacity"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	class := &entity.Class{
		TenantID:          req.TenantID,
		SchoolID:          req.SchoolID,
		AcademicYearID:    req.AcademicYearID,
		Name:              req.Name,
		Level:             req.Level,
		Major:             req.Major,
		HomeroomTeacherID: req.HomeroomTeacherID,
		Capacity:          req.Capacity,
	}

	if err := h.useCase.Create(c.Request.Context(), class); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, class)
}

func (h *ClassHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	class, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if class == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Class not found")
		return
	}

	httputil.Success(c.Writer, class)
}

func (h *ClassHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Tenant ID is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	classes, total, err := h.useCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]interface{}{
		"classes": classes,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *ClassHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		SchoolID          *uuid.UUID `json:"school_id"`
		AcademicYearID    *uuid.UUID `json:"academic_year_id"`
		Name              string     `json:"name" binding:"required"`
		Level             int        `json:"level"`
		Major             string     `json:"major"`
		HomeroomTeacherID *uuid.UUID `json:"homeroom_teacher_id"`
		Capacity          int        `json:"capacity"`
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
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Class not found")
		return
	}

	existing.SchoolID = req.SchoolID
	existing.AcademicYearID = req.AcademicYearID
	existing.Name = req.Name
	existing.Level = req.Level
	existing.Major = req.Major
	existing.HomeroomTeacherID = req.HomeroomTeacherID
	existing.Capacity = req.Capacity

	if err := h.useCase.Update(c.Request.Context(), existing); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

func (h *ClassHandler) Delete(c *gin.Context) {
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

	httputil.Success(c.Writer, map[string]string{"message": "Class deleted successfully"})
}
