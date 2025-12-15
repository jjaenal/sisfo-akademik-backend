package jwtutil

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAccessWith_ValidateWith_OK(t *testing.T) {
	secret := "secret"
	issuer := "sisfo-akademik"
	audience := "api"
	c := Claims{UserID: uuid.New(), TenantID: "t1", Roles: []string{"admin"}}
	s, err := GenerateAccessWith(secret, time.Minute, c, issuer, audience)
	if err != nil || s == "" {
		t.Fatalf("GenerateAccessWith failed: %v", err)
	}
	var out Claims
	if err := ValidateWith(secret, s, &out, issuer, audience); err != nil {
		t.Fatalf("ValidateWith failed: %v", err)
	}
	if out.UserID != c.UserID || out.TenantID != c.TenantID {
		t.Fatalf("claims mismatch")
	}
}

func TestValidateWith_WrongIssuer(t *testing.T) {
	secret := "s"
	issuer := "right"
	audience := "api"
	c := Claims{UserID: uuid.New(), TenantID: "t1"}
	s, _ := GenerateAccessWith(secret, time.Minute, c, issuer, audience)
	var out Claims
	if err := ValidateWith(secret, s, &out, "wrong-issuer", audience); err == nil {
		t.Fatalf("expected invalid issuer")
	}
}

func TestValidateWith_WrongAudience(t *testing.T) {
	secret := "s"
	issuer := "iss"
	audience := "api"
	c := Claims{UserID: uuid.New(), TenantID: "t1"}
	s, _ := GenerateAccessWith(secret, time.Minute, c, issuer, audience)
	var out Claims
	if err := ValidateWith(secret, s, &out, issuer, "other-aud"); err == nil {
		t.Fatalf("expected invalid audience")
	}
}
