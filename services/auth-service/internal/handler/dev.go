package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
)

type DevHandler struct {
	repo *repository.UsersRepo
}

func NewDevHandler(repo *repository.UsersRepo) *DevHandler {
	return &DevHandler{repo: repo}
}

func (h *DevHandler) Register(r *gin.Engine) {
	r.POST("/api/v1/auth/dev/bootstrap-user", h.bootstrapUser)
}

type devBootstrapReq struct {
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *DevHandler) bootstrapUser(c *gin.Context) {
	var in devBootstrapReq
	if err := c.BindJSON(&in); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid json")
		return
	}
	in.TenantID = strings.TrimSpace(in.TenantID)
	in.Email = strings.TrimSpace(in.Email)
	in.Password = strings.TrimSpace(in.Password)
	u, err := h.repo.Create(c.Request.Context(), repository.CreateUserParams{
		TenantID: in.TenantID,
		Email:    strings.ToLower(in.Email),
		Password: in.Password,
	})
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}
	httputil.Success(c.Writer, map[string]any{
		"id":        u.ID,
		"tenant_id": u.TenantID,
		"email":     u.Email,
	})
}
