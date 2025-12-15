package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

type UsersHandler struct {
	uc usecase.Users
}

func NewUsersHandler(uc usecase.Users) *UsersHandler {
	return &UsersHandler{uc: uc}
}

func (h *UsersHandler) RegisterProtected(r *gin.RouterGroup) {
	r.POST("/api/v1/users", h.create)
	r.GET("/api/v1/users/:id", h.get)
	r.GET("/api/v1/users", h.list)
	r.PUT("/api/v1/users/:id", h.update)
	r.DELETE("/api/v1/users/:id", h.delete)
}

type createReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UsersHandler) create(c *gin.Context) {
	val, _ := c.Get("claims")
	claims, _ := val.(jwtutil.Claims)
	in := createReq{}
	if err := c.BindJSON(&in); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid json")
		return
	}
	u, err := h.uc.Register(c.Request.Context(), usecase.UserRegisterInput{
		TenantID: claims.TenantID,
		Email:    in.Email,
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
		"is_active": u.IsActive,
	})
}

func (h *UsersHandler) get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid id")
		return
	}
	u, err := h.uc.Get(c.Request.Context(), id)
	if err != nil || !u.IsActive {
		httputil.Error(c.Writer, http.StatusNotFound, "5002", "Resource Not Found", "user not found")
		return
	}
	httputil.Success(c.Writer, map[string]any{
		"id":        u.ID,
		"tenant_id": u.TenantID,
		"email":     u.Email,
		"is_active": u.IsActive,
	})
}

func (h *UsersHandler) list(c *gin.Context) {
	val, _ := c.Get("claims")
	claims, _ := val.(jwtutil.Claims)
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	items, total, err := h.uc.List(c.Request.Context(), claims.TenantID, limit, offset)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	out := make([]map[string]any, 0, len(items))
	for _, u := range items {
		out = append(out, map[string]any{
			"id":        u.ID,
			"tenant_id": u.TenantID,
			"email":     u.Email,
			"is_active": u.IsActive,
		})
	}
	httputil.Success(c.Writer, map[string]any{
		"items": out,
		"total": total,
	})
}

type updateReq struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

func (h *UsersHandler) update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid id")
		return
	}
	var in updateReq
	if err := c.BindJSON(&in); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid json")
		return
	}
	u, err := h.uc.Update(c.Request.Context(), id, usecase.UserUpdateInput{
		Email:    in.Email,
		Password: in.Password,
		IsActive: in.IsActive,
	})
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", err.Error())
		return
	}
	httputil.Success(c.Writer, map[string]any{
		"id":        u.ID,
		"tenant_id": u.TenantID,
		"email":     u.Email,
		"is_active": u.IsActive,
	})
}

func (h *UsersHandler) delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid id")
		return
	}
	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	httputil.Success(c.Writer, map[string]any{"deleted": true})
}
