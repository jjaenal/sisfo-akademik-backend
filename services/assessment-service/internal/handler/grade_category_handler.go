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

// Create godoc
// @Summary      Create a grade category
// @Description  Create a new grade category
// @Tags         grade-categories
// @Accept       json
// @Produce      json
// @Param        request body object true "Grade Category Request"
// @Success      200  {object}  entity.GradeCategory
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grade-categories [post]
func (h *GradeCategoryHandler) Create(c *gin.Context) {
	var req struct {
		TenantID    string `json:"tenant_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Weight      float64 `json:"weight" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	category := &entity.GradeCategory{
		TenantID:    req.TenantID,
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

// GetByID godoc
// @Summary      Get a grade category
// @Description  Get a grade category by ID
// @Tags         grade-categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Grade Category ID"
// @Success      200  {object}  entity.GradeCategory
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      404  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grade-categories/{id} [get]
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

// List godoc
// @Summary      List grade categories
// @Description  List grade categories by tenant ID
// @Tags         grade-categories
// @Accept       json
// @Produce      json
// @Param        tenant_id query     string  true  "Tenant ID"
// @Success      200  {array}   entity.GradeCategory
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grade-categories [get]
func (h *GradeCategoryHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "tenant_id is required")
		return
	}

	categories, err := h.useCase.GetByTenantID(c.Request.Context(), tenantID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, categories)
}

// Update godoc
// @Summary      Update a grade category
// @Description  Update a grade category
// @Tags         grade-categories
// @Accept       json
// @Produce      json
// @Param        id           path      string  true  "Grade Category ID"
// @Param        request      body      object  true  "Update Request"
// @Success      200  {object}  entity.GradeCategory
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      404  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grade-categories/{id} [put]
func (h *GradeCategoryHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Weight      float64 `json:"weight" binding:"required,min=0"`
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

// Delete godoc
// @Summary      Delete a grade category
// @Description  Delete a grade category
// @Tags         grade-categories
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Grade Category ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /grade-categories/{id} [delete]
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
