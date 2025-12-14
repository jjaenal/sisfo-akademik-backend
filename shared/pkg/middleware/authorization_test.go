package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"

	"github.com/google/uuid"
)

type fakeAuthz struct {
	allow bool
	err   error
	last  struct {
		userID    uuid.UUID
		tenantID  string
		perm      string
	}
}

func (f *fakeAuthz) Allow(subjectID uuid.UUID, tenantID string, permission string) (bool, error) {
	f.last.userID = subjectID
	f.last.tenantID = tenantID
	f.last.perm = permission
	return f.allow, f.err
}

func TestAuthorization_Unauthenticated(t *testing.T) {
	authz := &fakeAuthz{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := Authorization(authz, "user:read", next)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthorization_Forbidden(t *testing.T) {
	authz := &fakeAuthz{allow: false}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })
	h := Authorization(authz, "user:read", next)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	claims := jwtutil.Claims{TenantID: "t1", UserID: uuid.New()}
	req = req.WithContext(withClaims(req.Context(), claims))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

func TestAuthorization_Allowed(t *testing.T) {
	authz := &fakeAuthz{allow: true}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })
	h := Authorization(authz, "user:read", next)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	claims := jwtutil.Claims{TenantID: "t1", UserID: uuid.New()}
	req = req.WithContext(withClaims(req.Context(), claims))
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	if authz.last.perm != "user:read" {
		t.Fatalf("unexpected permission: %s", authz.last.perm)
	}
}

func withClaims(ctx context.Context, c jwtutil.Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, c)
}
