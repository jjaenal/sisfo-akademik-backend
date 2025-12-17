package handler

import (
	"net/http"
	"strconv"
	"strings"

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

// Create creates a new schedule
// @Summary      Create a new schedule
// @Description  Create a single schedule entry
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        request body handler.CreateScheduleRequest true "Schedule Request"
// @Success      200  {object}  entity.Schedule
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules [post]
func (h *ScheduleHandler) Create(c *gin.Context) {
	var req CreateScheduleRequest

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

// BulkCreate creates multiple schedules
// @Summary      Bulk create schedules
// @Description  Create multiple schedule entries at once
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        request body []handler.CreateScheduleRequest true "List of Schedule Requests"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      409  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules/bulk [post]
func (h *ScheduleHandler) BulkCreate(c *gin.Context) {
	var reqs []CreateScheduleRequest

	if err := c.ShouldBindJSON(&reqs); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	var schedules []*entity.Schedule
	for _, req := range reqs {
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
		schedules = append(schedules, schedule)
	}

	if err := h.useCase.BulkCreate(c.Request.Context(), schedules); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "conflict") {
			httputil.Error(c.Writer, http.StatusConflict, "4009", "Conflict", err.Error())
			return
		}
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]interface{}{
		"count": len(schedules),
	})
}

type CreateScheduleRequest struct {
	TenantID  string    `json:"tenant_id" binding:"required"`
	ClassID   uuid.UUID `json:"class_id" binding:"required"`
	SubjectID uuid.UUID `json:"subject_id" binding:"required"`
	TeacherID uuid.UUID `json:"teacher_id" binding:"required"`
	DayOfWeek int       `json:"day_of_week" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Room      string    `json:"room"`
}

// CreateFromTemplate creates schedules from a template
// @Summary      Create schedules from template
// @Description  Create multiple schedules based on a predefined template and teacher assignments
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        request body handler.CreateFromTemplateRequest true "Template Assignment Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      404  {object}  httputil.ErrorResponse
// @Failure      409  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules/from-template [post]
func (h *ScheduleHandler) CreateFromTemplate(c *gin.Context) {
	var req CreateFromTemplateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	teacherMap := make(map[uuid.UUID]uuid.UUID)
	for _, a := range req.Assignments {
		teacherMap[a.SubjectID] = a.TeacherID
	}

	if err := h.useCase.CreateFromTemplate(c.Request.Context(), req.TemplateID, req.ClassID, teacherMap); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "conflict") {
			httputil.Error(c.Writer, http.StatusConflict, "4009", "Conflict", err.Error())
			return
		}
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", err.Error())
			return
		}
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, map[string]string{
		"message": "schedules created successfully from template",
	})
}

type CreateFromTemplateRequest struct {
	TemplateID  uuid.UUID `json:"template_id" binding:"required"`
	ClassID     uuid.UUID `json:"class_id" binding:"required"`
	Assignments []struct {
		SubjectID uuid.UUID `json:"subject_id" binding:"required"`
		TeacherID uuid.UUID `json:"teacher_id" binding:"required"`
	} `json:"assignments" binding:"required"`
}

// GetByID gets a schedule by ID
// @Summary      Get schedule by ID
// @Description  Get detailed information about a schedule
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Success      200  {object}  entity.Schedule
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      404  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules/{id} [get]
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

// List lists schedules
// @Summary      List schedules
// @Description  Get a list of schedules with pagination
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        tenant_id  query     string  true  "Tenant ID"
// @Param        limit      query     int     false "Limit" default(10)
// @Param        offset     query     int     false "Offset" default(0)
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules [get]
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

// ListByClass lists schedules by class
// @Summary      List schedules by class
// @Description  Get all schedules for a specific class
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        class_id   path      string  true  "Class ID"
// @Success      200  {object}  []entity.Schedule
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules/class/{class_id} [get]
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

// Update updates a schedule
// @Summary      Update a schedule
// @Description  Update details of an existing schedule
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        id           path      string  true  "Schedule ID"
// @Param        request body handler.UpdateScheduleRequest true "Update Schedule Request"
// @Success      200  {object}  entity.Schedule
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      404  {object}  httputil.ErrorResponse
// @Failure      409  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules/{id} [put]
func (h *ScheduleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req UpdateScheduleRequest

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
		if strings.Contains(strings.ToLower(err.Error()), "conflict") {
			httputil.Error(c.Writer, http.StatusConflict, "4009", "Conflict", err.Error())
			return
		}
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, existing)
}

type UpdateScheduleRequest struct {
	ClassID   uuid.UUID `json:"class_id" binding:"required"`
	SubjectID uuid.UUID `json:"subject_id" binding:"required"`
	TeacherID uuid.UUID `json:"teacher_id" binding:"required"`
	DayOfWeek int       `json:"day_of_week" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Room      string    `json:"room"`
}

// Delete deletes a schedule
// @Summary      Delete a schedule
// @Description  Delete a schedule by ID
// @Tags         schedules
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Schedule ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  httputil.ErrorResponse
// @Failure      500  {object}  httputil.ErrorResponse
// @Router       /schedules/{id} [delete]
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
