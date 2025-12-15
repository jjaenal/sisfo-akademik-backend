package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func testDBRolesH(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DB_URL")
	if url == "" {
		url = "postgres://dev:dev@localhost:55432/devdb?sslmode=disable"
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Skip("no db available:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Skip("db connect failed:", err)
	}
	return db
}

func ensureMigrationsRolesH(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, _ = db.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	b1, err := os.ReadFile("../../migrations/001_users_roles_permissions.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b1))
	}
	b2, err := os.ReadFile("../../migrations/002_audit_logs.up.sql")
	if err == nil {
		_, _ = db.Exec(ctx, string(b2))
	}
}

func TestRolesHandler_AssignListUnassign(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBRolesH(t)
	ensureMigrationsRolesH(t, db)
	usersRepo := repository.NewUsersRepo(db)
	rolesRepo := repository.NewRolesRepo(db)
	rolesUC := usecase.NewRoles(usersRepo, rolesRepo)
	h := NewRolesHandler(rolesUC)
	usersUC := usecase.NewUsers(usersRepo)
	usersH := NewUsersHandler(usersUC)

	r := gin.New()
	tenant := "t-" + uuid.NewString()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: uuid.New()})
	})
	usersH.RegisterProtected(protected)
	h.RegisterProtected(protected)

	// create user
	body := map[string]string{"email": "u2@test.local", "password": "password123"}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create status=%d body=%s", w.Code, w.Body.String())
	}
	var created struct {
		Data map[string]any `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	idStr, _ := created.Data["id"].(string)
	if idStr == "" {
		t.Fatalf("missing id in response: %s", w.Body.String())
	}

	// assign role
	assignBody := map[string]string{"role_name": "teacher"}
	b, _ = json.Marshal(assignBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/"+idStr+"/roles", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("assign status=%d body=%s", w.Code, w.Body.String())
	}

	// assign with empty role_name -> expect 400
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/"+idStr+"/roles", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}

	// list roles
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/"+idStr+"/roles", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list roles status=%d body=%s", w.Code, w.Body.String())
	}

	// extract role_id from list
	var listResp struct {
		Data map[string][]map[string]any `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &listResp)
	items := listResp.Data["items"]
	if len(items) == 0 {
		t.Fatalf("empty roles list")
	}
	roleID, _ := items[0]["id"].(string)
	if roleID == "" {
		t.Fatalf("missing role id")
	}

	// unassign
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/"+idStr+"/roles/"+roleID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("unassign status=%d body=%s", w.Code, w.Body.String())
	}

	// invalid unassign id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/"+idStr+"/roles/not-a-uuid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}
