package handler

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

type AuditHandler struct {
	repo *repository.AuditRepo
}

func NewAuditHandler(repo *repository.AuditRepo) *AuditHandler {
	return &AuditHandler{repo: repo}
}

func (h *AuditHandler) RegisterProtected(g *gin.RouterGroup) {
	g.GET("/api/v1/audit-logs", h.list)
	g.GET("/api/v1/audit-logs/search", h.search)
	g.GET("/api/v1/audit-logs/export", h.export)
}

func (h *AuditHandler) list(c *gin.Context) {
	val, ok := c.Get("claims")
	if !ok {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", nil)
		return
	}
	claims, _ := val.(jwtutil.Claims)
	lim, _ := strconv.Atoi(c.Query("limit"))
	off, _ := strconv.Atoi(c.Query("offset"))
	p := repository.ListParams{}
	if s := c.Query("user_id"); s != "" {
		if id, err := uuid.Parse(s); err == nil {
			p.UserID = &id
		}
	}
	if s := c.Query("action"); s != "" {
		p.Action = s
	}
	if s := c.Query("resource_type"); s != "" {
		p.ResourceType = s
	}
	if s := c.Query("resource_id"); s != "" {
		if id, err := uuid.Parse(s); err == nil {
			p.ResourceID = &id
		}
	}
	if s := c.Query("start"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			p.Start = &t
		}
	}
	if s := c.Query("end"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			p.End = &t
		}
	}
	items, total, err := h.repo.List(c.Request.Context(), claims.TenantID, p, lim, off)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	httputil.Success(c.Writer, map[string]any{"items": items, "total": total})
}

func (h *AuditHandler) search(c *gin.Context) {
	val, ok := c.Get("claims")
	if !ok {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", nil)
		return
	}
	claims, _ := val.(jwtutil.Claims)
	q := c.Query("q")
	if q == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "missing q")
		return
	}
	lim, _ := strconv.Atoi(c.Query("limit"))
	off, _ := strconv.Atoi(c.Query("offset"))
	items, total, err := h.repo.Search(c.Request.Context(), claims.TenantID, q, lim, off)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	httputil.Success(c.Writer, map[string]any{"items": items, "total": total})
}

func (h *AuditHandler) export(c *gin.Context) {
	val, ok := c.Get("claims")
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	claims, _ := val.(jwtutil.Claims)
	p := repository.ListParams{}
	lim, _ := strconv.Atoi(c.Query("limit"))
	off, _ := strconv.Atoi(c.Query("offset"))
	items, _, err := h.repo.List(c.Request.Context(), claims.TenantID, p, lim, off)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=\"audit_logs.csv\"")
	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"id", "tenant_id", "user_id", "action", "resource_type", "resource_id", "created_at"})
	for _, it := range items {
		uid := ""
		if it.UserID != nil {
			uid = it.UserID.String()
		}
		rid := ""
		if it.ResourceID != nil {
			rid = it.ResourceID.String()
		}
		_ = w.Write([]string{
			it.ID.String(),
			it.TenantID,
			uid,
			it.Action,
			it.ResourceType,
			rid,
			it.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	w.Flush()
}

