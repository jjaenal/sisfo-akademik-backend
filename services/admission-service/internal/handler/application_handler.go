package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/entity"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/domain/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type ApplicationHandler struct {
	useCase usecase.ApplicationUseCase
}

func NewApplicationHandler(useCase usecase.ApplicationUseCase) *ApplicationHandler {
	return &ApplicationHandler{useCase: useCase}
}

func (h *ApplicationHandler) Submit(c *gin.Context) {
	var req struct {
		AdmissionPeriodID string  `json:"admission_period_id" binding:"required"`
		FirstName         string  `json:"first_name" binding:"required"`
		LastName          string  `json:"last_name" binding:"required"`
		Email             string  `json:"email" binding:"required,email"`
		PhoneNumber       string  `json:"phone_number" binding:"required"`
		PreviousSchool    string  `json:"previous_school" binding:"required"`
		AverageScore      float64 `json:"average_score" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	periodID, err := uuid.Parse(req.AdmissionPeriodID)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Invalid Admission Period ID")
		return
	}

	application := &entity.Application{
		AdmissionPeriodID: periodID,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		PhoneNumber:       req.PhoneNumber,
		PreviousSchool:    req.PreviousSchool,
		AverageScore:      req.AverageScore,
	}

	submittedApp, err := h.useCase.SubmitApplication(c.Request.Context(), application)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, submittedApp)
}

func (h *ApplicationHandler) GetStatus(c *gin.Context) {
	regNumber := c.Query("registration_number")
	if regNumber == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "Registration number is required")
		return
	}

	app, err := h.useCase.GetApplicationStatus(c.Request.Context(), regNumber)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}
	if app == nil {
		httputil.Error(c.Writer, http.StatusNotFound, "4004", "Not Found", "Application not found")
		return
	}

	httputil.Success(c.Writer, app)
}

func (h *ApplicationHandler) Verify(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	if err := h.useCase.VerifyApplication(c.Request.Context(), id, entity.ApplicationStatus(req.Status)); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Application status updated"})
}

func (h *ApplicationHandler) List(c *gin.Context) {
	filter := make(map[string]interface{})
	
	if periodID := c.Query("admission_period_id"); periodID != "" {
		if id, err := uuid.Parse(periodID); err == nil {
			filter["admission_period_id"] = id
		}
	}
	
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	apps, err := h.useCase.ListApplications(c.Request.Context(), filter)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, apps)
}

func (h *ApplicationHandler) InputTestScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Score float64 `json:"score" binding:"required,min=0,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	if err := h.useCase.InputTestScore(c.Request.Context(), id, req.Score); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Test score updated"})
}

func (h *ApplicationHandler) InputInterviewScore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	var req struct {
		Score float64 `json:"score" binding:"required,min=0,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}

	if err := h.useCase.InputInterviewScore(c.Request.Context(), id, req.Score); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Interview score updated"})
}

func (h *ApplicationHandler) CalculateFinalScores(c *gin.Context) {
	idStr := c.Param("id") // periodID
	periodID, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.CalculateFinalScores(c.Request.Context(), periodID); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Final scores calculated"})
}

func (h *ApplicationHandler) Register(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid ID", "ID must be a valid UUID")
		return
	}

	if err := h.useCase.RegisterStudent(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "5001", "Internal Server Error", err.Error())
		return
	}

	httputil.Success(c.Writer, gin.H{"message": "Student registered successfully"})
}
