package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/academic-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type StudentHandler struct {
	useCase usecase.StudentUseCase
}

func NewStudentHandler(useCase usecase.StudentUseCase) *StudentHandler {
	return &StudentHandler{useCase: useCase}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req struct {
		TenantID      string     `json:"tenant_id" binding:"required"`
		UserID        *uuid.UUID `json:"user_id"`
		NIS           string     `json:"nis"`
		NISN          string     `json:"nisn"`
		Name          string     `json:"name" binding:"required"`
		Gender        string     `json:"gender"`
		BirthPlace    string     `json:"birth_place"`
		BirthDate     *time.Time `json:"birth_date"`
		Address       string     `json:"address"`
		Phone         string     `json:"phone"`
		Email         string     `json:"email"`
		ParentName    string     `json:"parent_name"`
		ParentPhone   string     `json:"parent_phone"`
		AdmissionDate *time.Time `json:"admission_date"`
		Status        string     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	student := &entity.Student{
		TenantID:      req.TenantID,
		UserID:        req.UserID,
		NIS:           req.NIS,
		NISN:          req.NISN,
		Name:          req.Name,
		Gender:        req.Gender,
		BirthPlace:    req.BirthPlace,
		BirthDate:     req.BirthDate,
		Address:       req.Address,
		Phone:         req.Phone,
		Email:         req.Email,
		ParentName:    req.ParentName,
		ParentPhone:   req.ParentPhone,
		AdmissionDate: req.AdmissionDate,
		Status:        req.Status,
	}

	if err := h.useCase.Create(c.Request.Context(), student); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, student)
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	student, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if student == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Student not found")
		return
	}

	httputil.Success(c.Writer, student)
}

func (h *StudentHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Tenant ID is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	students, total, err := h.useCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]interface{}{
		"students": students,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *StudentHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		UserID        *uuid.UUID `json:"user_id"`
		NIS           string     `json:"nis"`
		NISN          string     `json:"nisn"`
		Name          string     `json:"name" binding:"required"`
		Gender        string     `json:"gender"`
		BirthPlace    string     `json:"birth_place"`
		BirthDate     *time.Time `json:"birth_date"`
		Address       string     `json:"address"`
		Phone         string     `json:"phone"`
		Email         string     `json:"email"`
		ParentName    string     `json:"parent_name"`
		ParentPhone   string     `json:"parent_phone"`
		AdmissionDate *time.Time `json:"admission_date"`
		Status        string     `json:"status"`
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
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Student not found")
		return
	}

	existing.UserID = req.UserID
	existing.NIS = req.NIS
	existing.NISN = req.NISN
	existing.Name = req.Name
	existing.Gender = req.Gender
	existing.BirthPlace = req.BirthPlace
	existing.BirthDate = req.BirthDate
	existing.Address = req.Address
	existing.Phone = req.Phone
	existing.Email = req.Email
	existing.ParentName = req.ParentName
	existing.ParentPhone = req.ParentPhone
	existing.AdmissionDate = req.AdmissionDate
	existing.Status = req.Status

	if err := h.useCase.Update(c.Request.Context(), existing); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

func (h *StudentHandler) Delete(c *gin.Context) {
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

	httputil.Success(c.Writer, map[string]string{"message": "Student deleted successfully"})
}
