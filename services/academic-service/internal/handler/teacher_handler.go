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

type TeacherHandler struct {
	useCase usecase.TeacherUseCase
}

func NewTeacherHandler(useCase usecase.TeacherUseCase) *TeacherHandler {
	return &TeacherHandler{useCase: useCase}
}

func (h *TeacherHandler) Create(c *gin.Context) {
	var req struct {
		TenantID   string     `json:"tenant_id" binding:"required"`
		UserID     *uuid.UUID `json:"user_id"`
		NIP        string     `json:"nip"`
		Name       string     `json:"name" binding:"required"`
		Gender     string     `json:"gender"`
		TitleFront string     `json:"title_front"`
		TitleBack  string     `json:"title_back"`
		Phone      string     `json:"phone"`
		Email      string     `json:"email"`
		Status     string     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	teacher := &entity.Teacher{
		TenantID:   req.TenantID,
		UserID:     req.UserID,
		NIP:        req.NIP,
		Name:       req.Name,
		Gender:     req.Gender,
		TitleFront: req.TitleFront,
		TitleBack:  req.TitleBack,
		Phone:      req.Phone,
		Email:      req.Email,
		Status:     req.Status,
	}

	if err := h.useCase.Create(c.Request.Context(), teacher); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, teacher)
}

func (h *TeacherHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	teacher, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if teacher == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Teacher not found")
		return
	}

	httputil.Success(c.Writer, teacher)
}

func (h *TeacherHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Tenant ID is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	teachers, total, err := h.useCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]interface{}{
		"teachers": teachers,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *TeacherHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		UserID     *uuid.UUID `json:"user_id"`
		NIP        string     `json:"nip"`
		Name       string     `json:"name" binding:"required"`
		Gender     string     `json:"gender"`
		TitleFront string     `json:"title_front"`
		TitleBack  string     `json:"title_back"`
		Phone      string     `json:"phone"`
		Email      string     `json:"email"`
		Status     string     `json:"status"`
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
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Teacher not found")
		return
	}

	existing.UserID = req.UserID
	existing.NIP = req.NIP
	existing.Name = req.Name
	existing.Gender = req.Gender
	existing.TitleFront = req.TitleFront
	existing.TitleBack = req.TitleBack
	existing.Phone = req.Phone
	existing.Email = req.Email
	existing.Status = req.Status

	if err := h.useCase.Update(c.Request.Context(), existing); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

func (h *TeacherHandler) Delete(c *gin.Context) {
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

	httputil.Success(c.Writer, map[string]string{"message": "Teacher deleted successfully"})
}
