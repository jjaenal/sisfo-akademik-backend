package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func TestAuditHandler_List_Search_Export(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBAuth(t)
	ensureMigrationsAuth(t, db)
	auditRepo := repository.NewAuditRepo(db)
	usersRepo := repository.NewUsersRepo(db)
	tenant := "t-" + uuid.NewString()
	u, err := usersRepo.Create(context.Background(), repository.CreateUserParams{
		TenantID: tenant, Email: "audit@test.local", Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("seed user err: %v", err)
	}
	rid := uuid.New()
	if err := auditRepo.Log(context.Background(), tenant, &u.ID, "auth.login", "user", &rid, map[string]any{"success": true}); err != nil {
		t.Fatalf("log err: %v", err)
	}
	if err := auditRepo.Log(context.Background(), tenant, &u.ID, "auth.refresh", "user", &rid, map[string]any{"success": true}); err != nil {
		t.Fatalf("log2 err: %v", err)
	}
	h := NewAuditHandler(auditRepo)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: u.ID})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs?limit=10", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", w.Code, w.Body.String())
	}
	var resp struct {
		Data struct {
			Items []repository.AuditLog `json:"items"`
			Total int                   `json:"total"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data.Items) < 2 || resp.Data.Total < 2 {
		t.Fatalf("expected at least 2 items, got %d total=%d", len(resp.Data.Items), resp.Data.Total)
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs/search?q=auth.login", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("search status=%d body=%s", w.Code, w.Body.String())
	}
	var sresp struct {
		Data struct {
			Items []repository.AuditLog `json:"items"`
			Total int                   `json:"total"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &sresp)
	if sresp.Data.Total < 1 {
		t.Fatalf("expected search results >=1")
	}
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs/export", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("export status=%d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "text/csv" {
		t.Fatalf("expected text/csv got %s", ct)
	}
}
