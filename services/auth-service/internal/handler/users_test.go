package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func testDBUsers(t *testing.T) *pgxpool.Pool {
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

func ensureMigrationsUsers(t *testing.T, db *pgxpool.Pool) {
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

func TestUsersHandler_CRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testDBUsers(t)
	ensureMigrationsUsers(t, db)
	usersRepo := repository.NewUsersRepo(db)
	uc := usecase.NewUsers(usersRepo)
	h := NewUsersHandler(uc)

	r := gin.New()
	tenant := "t-" + uuid.NewString()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: tenant, UserID: uuid.New()})
	})
	h.RegisterProtected(protected)

	// create
	body := map[string]string{"email": "u1@test.local", "password": "Password123!"}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("create status=%d body=%s", w.Code, w.Body.String())
	}
	var created struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &created)
	idStr, _ := created.Data["id"].(string)
	if idStr == "" {
		t.Fatalf("missing id in response: %s", w.Body.String())
	}

	// list
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users?limit=10&offset=0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", w.Code, w.Body.String())
	}

	// update
	newEmail := map[string]string{"email": "u1b@test.local"}
	b, _ = json.Marshal(newEmail)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/"+idStr, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update status=%d body=%s", w.Code, w.Body.String())
	}

	// update with short password -> expect 400
	shortPwd := map[string]string{"password": "short"}
	b, _ = json.Marshal(shortPwd)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/"+idStr, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for short password, got %d body=%s", w.Code, w.Body.String())
	}

	// create with invalid email -> expect 400
	w = httptest.NewRecorder()
	bodyInvalidEmail := map[string]string{"email": "invalid-email", "password": "password123!"}
	b, _ = json.Marshal(bodyInvalidEmail)
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid email, got %d body=%s", w.Code, w.Body.String())
	}

	// update with invalid email -> expect 400
	w = httptest.NewRecorder()
	invalidEmailUpd := map[string]string{"email": "not-an-email"}
	b, _ = json.Marshal(invalidEmailUpd)
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/"+idStr, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid email update, got %d body=%s", w.Code, w.Body.String())
	}

	// update with weak password (no symbol) -> expect 400
	w = httptest.NewRecorder()
	weakPwd := map[string]string{"password": "password123"}
	b, _ = json.Marshal(weakPwd)
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/"+idStr, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for weak password update, got %d body=%s", w.Code, w.Body.String())
	}

	// update with empty body -> expect 200 and unchanged email
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPut, "/api/v1/users/"+idStr, bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for empty update, got %d body=%s", w.Code, w.Body.String())
	}

	// delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, "/api/v1/users/"+idStr, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("delete status=%d body=%s", w.Code, w.Body.String())
	}

	// get after delete -> expect 404
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/"+idStr, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleted user, got %d body=%s", w.Code, w.Body.String())
	}
}
