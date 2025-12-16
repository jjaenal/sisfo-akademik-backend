package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type SchoolHandler struct {
	useCase usecase.SchoolUseCase
}

func NewSchoolHandler(useCase usecase.SchoolUseCase) *SchoolHandler {
	return &SchoolHandler{useCase: useCase}
}

// Create handles school creation
func (h *SchoolHandler) Create(c *gin.Context) {
	var req struct {
		TenantID      string `json:"tenant_id" binding:"required"`
		Name          string `json:"name" binding:"required"`
		Address       string `json:"address"`
		Phone         string `json:"phone"`
		Email         string `json:"email"`
		Website       string `json:"website"`
		LogoURL       string `json:"logo_url"`
		Accreditation string `json:"accreditation"`
		Headmaster    string `json:"headmaster"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	school := &entity.School{
		TenantID:      req.TenantID,
		Name:          req.Name,
		Address:       req.Address,
		Phone:         req.Phone,
		Email:         req.Email,
		Website:       req.Website,
		LogoURL:       req.LogoURL,
		Accreditation: req.Accreditation,
		Headmaster:    req.Headmaster,
	}

	// TODO: Get UserID from context for CreatedBy
	// For now, assuming middleware sets user_id in context
	// if userID, ok := c.Get("user_id"); ok {
	// 	uid := uuid.MustParse(userID.(string))
	// 	school.CreatedBy = &uid
	// 	school.UpdatedBy = &uid
	// }

	if err := h.useCase.Create(c.Request.Context(), school); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, school)
}

// GetByID retrieves a school by ID
func (h *SchoolHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	school, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if school == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "School not found")
		return
	}

	httputil.Success(c.Writer, school)
}

// GetByTenantID retrieves a school by Tenant ID
func (h *SchoolHandler) GetByTenantID(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Tenant ID is required")
		return
	}

	school, err := h.useCase.GetByTenantID(c.Request.Context(), tenantID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if school == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "School not found for this tenant")
		return
	}

	httputil.Success(c.Writer, school)
}

// Update handles school updates
func (h *SchoolHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Name          string `json:"name" binding:"required"`
		Address       string `json:"address"`
		Phone         string `json:"phone"`
		Email         string `json:"email"`
		Website       string `json:"website"`
		LogoURL       string `json:"logo_url"`
		Accreditation string `json:"accreditation"`
		Headmaster    string `json:"headmaster"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	// First verify exists
	existing, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if existing == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "School not found")
		return
	}

	// Update fields
	existing.Name = req.Name
	existing.Address = req.Address
	existing.Phone = req.Phone
	existing.Email = req.Email
	existing.Website = req.Website
	existing.LogoURL = req.LogoURL
	existing.Accreditation = req.Accreditation
	existing.Headmaster = req.Headmaster

	// TODO: Get UserID for UpdatedBy
	// if userID, ok := c.Get("user_id"); ok {
	// 	uid := uuid.MustParse(userID.(string))
	// 	existing.UpdatedBy = &uid
	// }

	if err := h.useCase.Update(c.Request.Context(), existing); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

// Delete handles school deletion
func (h *SchoolHandler) Delete(c *gin.Context) {
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

	httputil.Success(c.Writer, map[string]string{"message": "School deleted successfully"})
}
