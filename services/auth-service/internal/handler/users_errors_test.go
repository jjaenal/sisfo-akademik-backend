package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
)

func TestUsersHandler_InvalidID_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &stubUsers{}
	h := NewUsersHandler(uc)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t-test", UserID: uuid.New()})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/not-a-uuid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestUsersHandler_InvalidJSON_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := &stubUsers{}
	h := NewUsersHandler(uc)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t-test", UserID: uuid.New()})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUsersHandler_MissingEmail_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := usecase.NewUsers(repository.NewUsersRepo(nil))
	h := NewUsersHandler(uc)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t-test", UserID: uuid.New()})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(`{"password":"password123!"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUsersHandler_MissingPassword_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := usecase.NewUsers(repository.NewUsersRepo(nil))
	h := NewUsersHandler(uc)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t-test", UserID: uuid.New()})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(`{"email":"a@b.c"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestUsersHandler_WeakPassword_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uc := usecase.NewUsers(repository.NewUsersRepo(nil))
	h := NewUsersHandler(uc)
	r := gin.New()
	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("claims", jwtutil.Claims{TenantID: "t-test", UserID: uuid.New()})
	})
	h.RegisterProtected(protected)
	w := httptest.NewRecorder()
	reqBody := `{"email":"weak@test.local","password":"weakpass"}`
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d body=%s", w.Code, w.Body.String())
	}
}

type stubUsers struct{}

func (*stubUsers) Register(_ context.Context, _ usecase.UserRegisterInput) (*repository.User, error) {
	return nil, nil
}
func (*stubUsers) Get(_ context.Context, _ uuid.UUID) (*repository.User, error) { return nil, nil }
func (*stubUsers) List(_ context.Context, _ string, _ int, _ int) ([]repository.User, int, error) {
	return nil, 0, nil
}
func (*stubUsers) Update(_ context.Context, _ uuid.UUID, _ usecase.UserUpdateInput) (*repository.User, error) {
	return nil, nil
}
func (*stubUsers) Delete(_ context.Context, _ uuid.UUID) error { return nil }
