package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

type RolesHandler struct {
	uc usecase.Roles
}

func NewRolesHandler(uc usecase.Roles) *RolesHandler {
	return &RolesHandler{uc: uc}
}

func (h *RolesHandler) RegisterProtected(r *gin.RouterGroup) {
	r.POST("/api/v1/users/:id/roles", h.assign)
	r.GET("/api/v1/users/:id/roles", h.list)
	r.DELETE("/api/v1/users/:id/roles/:role_id", h.unassign)
}

type assignReq struct {
	RoleName string `json:"role_name"`
}

func (h *RolesHandler) assign(c *gin.Context) {
	val, _ := c.Get("claims")
	claims, _ := val.(jwtutil.Claims)
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid user id")
		return
	}
	var in assignReq
	if err := c.BindJSON(&in); err != nil || strings.TrimSpace(in.RoleName) == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "role_name required")
		return
	}
	role, err := h.uc.AssignByName(c.Request.Context(), claims.TenantID, uid, in.RoleName)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}
	httputil.Success(c.Writer, map[string]any{
		"role": map[string]any{
			"id":             role.ID,
			"name":           role.Name,
			"tenant_id":      role.TenantID,
			"is_system_role": role.IsSystemRole,
		},
	})
}

func (h *RolesHandler) list(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid user id")
		return
	}
	roles, err := h.uc.List(c.Request.Context(), uid)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	items := make([]map[string]any, 0, len(roles))
	for _, r := range roles {
		items = append(items, map[string]any{
			"id":             r.ID,
			"name":           r.Name,
			"tenant_id":      r.TenantID,
			"is_system_role": r.IsSystemRole,
		})
	}
	httputil.Success(c.Writer, map[string]any{"items": items})
}

func (h *RolesHandler) unassign(c *gin.Context) {
	uid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid user id")
		return
	}
	rid, err := uuid.Parse(c.Param("role_id"))
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid role id")
		return
	}
	if err := h.uc.Unassign(c.Request.Context(), uid, rid); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	httputil.Success(c.Writer, map[string]any{"deleted": true})
}
