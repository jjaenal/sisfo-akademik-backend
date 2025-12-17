package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type GradeCategoryHandler struct {
	useCase usecase.GradeCategoryUseCase
}

func NewGradeCategoryHandler(useCase usecase.GradeCategoryUseCase) *GradeCategoryHandler {
	return &GradeCategoryHandler{useCase: useCase}
}

func (h *GradeCategoryHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Weight      int    `json:"weight" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	category := &entity.GradeCategory{
		Name:        req.Name,
		Description: req.Description,
		Weight:      req.Weight,
	}

	if err := h.useCase.Create(c.Request.Context(), category); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, category)
}

func (h *GradeCategoryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	category, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if category == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Grade category not found")
		return
	}

	httputil.Success(c.Writer, category)
}

func (h *GradeCategoryHandler) List(c *gin.Context) {
	categories, err := h.useCase.List(c.Request.Context())
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, categories)
}

func (h *GradeCategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Weight      int    `json:"weight" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	category, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if category == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Grade category not found")
		return
	}

	category.Name = req.Name
	category.Description = req.Description
	category.Weight = req.Weight

	if err := h.useCase.Update(c.Request.Context(), category); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, category)
}

func (h *GradeCategoryHandler) Delete(c *gin.Context) {
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

	httputil.Success(c.Writer, gin.H{"message": "Grade category deleted successfully"})
}
