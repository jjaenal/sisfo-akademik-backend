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

type ScheduleHandler struct {
	useCase usecase.ScheduleUseCase
}

func NewScheduleHandler(useCase usecase.ScheduleUseCase) *ScheduleHandler {
	return &ScheduleHandler{useCase: useCase}
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var req struct {
		TenantID  string    `json:"tenant_id" binding:"required"`
		ClassID   uuid.UUID `json:"class_id" binding:"required"`
		SubjectID uuid.UUID `json:"subject_id" binding:"required"`
		TeacherID uuid.UUID `json:"teacher_id" binding:"required"`
		DayOfWeek int       `json:"day_of_week" binding:"required"`
		StartTime string    `json:"start_time" binding:"required"`
		EndTime   string    `json:"end_time" binding:"required"`
		Room      string    `json:"room"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	schedule := &entity.Schedule{
		TenantID:  req.TenantID,
		ClassID:   req.ClassID,
		SubjectID: req.SubjectID,
		TeacherID: req.TeacherID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Room:      req.Room,
	}

	if err := h.useCase.Create(c.Request.Context(), schedule); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, schedule)
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	schedule, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if schedule == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Schedule not found")
		return
	}

	httputil.Success(c.Writer, schedule)
}

func (h *ScheduleHandler) List(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Tenant ID is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	schedules, total, err := h.useCase.List(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]interface{}{
		"schedules": schedules,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *ScheduleHandler) ListByClass(c *gin.Context) {
	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "Class ID must be a valid UUID")
		return
	}

	schedules, err := h.useCase.ListByClass(c.Request.Context(), classID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, schedules)
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		ClassID   uuid.UUID `json:"class_id" binding:"required"`
		SubjectID uuid.UUID `json:"subject_id" binding:"required"`
		TeacherID uuid.UUID `json:"teacher_id" binding:"required"`
		DayOfWeek int       `json:"day_of_week" binding:"required"`
		StartTime string    `json:"start_time" binding:"required"`
		EndTime   string    `json:"end_time" binding:"required"`
		Room      string    `json:"room"`
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
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Schedule not found")
		return
	}

	existing.ClassID = req.ClassID
	existing.SubjectID = req.SubjectID
	existing.TeacherID = req.TeacherID
	existing.DayOfWeek = req.DayOfWeek
	existing.StartTime = req.StartTime
	existing.EndTime = req.EndTime
	existing.Room = req.Room

	if err := h.useCase.Update(c.Request.Context(), existing); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
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

	httputil.Success(c.Writer, map[string]string{"message": "Schedule deleted successfully"})
}
