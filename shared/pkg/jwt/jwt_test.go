package jwtutil

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateAndValidate(t *testing.T) {
	secret := "secret"
	c := Claims{UserID: uuid.New(), TenantID: "t1", Roles: []string{"r1"}}
	s, err := GenerateAccess(secret, 15*time.Minute, c)
	if err != nil || s == "" {
		t.Fatalf("GenerateAccess failed")
	}
	var out Claims
	if err := Validate(secret, s, &out); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if out.UserID != c.UserID || out.TenantID != c.TenantID {
		t.Fatalf("claims mismatch")
	}
}

func TestValidateWrongSecret(t *testing.T) {
	secret := "s1"
	c := Claims{UserID: uuid.New(), TenantID: "t1"}
	s, _ := GenerateAccess(secret, time.Minute, c)
	var out Claims
	if err := Validate("s2", s, &out); err == nil {
		t.Fatalf("expected error for wrong secret")
	}
}

func TestExpiredToken(t *testing.T) {
	secret := "s"
	c := Claims{UserID: uuid.New(), TenantID: "t"}
	s, _ := GenerateAccess(secret, time.Millisecond*10, c)
	time.Sleep(time.Millisecond * 20)
	var out Claims
	if err := Validate(secret, s, &out); err == nil {
		t.Fatalf("expected expired token error")
	}
}

func TestInvalidAlg(t *testing.T) {
	secret := "s"
	now := time.Now()
	rc := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
	// Buat token dengan algoritma berbeda (HS384) untuk memicu penolakan
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS384, rc)
	s, _ := tkn.SignedString([]byte(secret))
	var out Claims
	if err := Validate(secret, s, &out); err == nil {
		t.Fatalf("expected error for invalid alg")
	}
}

func TestGenerateRefreshAndValidate(t *testing.T) {
	secret := "s"
	sub := "user-123"
	s, err := GenerateRefresh(secret, time.Minute, sub)
	if err != nil || s == "" {
		t.Fatalf("GenerateRefresh failed")
	}
	// Validasi refresh token menggunakan Claims untuk mendapatkan RegisteredClaims.Subject
	var out Claims
	if err := Validate(secret, s, &out); err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if out.Subject != sub {
		t.Fatalf("subject mismatch")
	}
}

func TestMalformedToken(t *testing.T) {
	var out Claims
	if err := Validate("s", "not-a-jwt", &out); err == nil {
		t.Fatalf("expected parse error for malformed token")
	}
}

func TestNotBeforeFuture(t *testing.T) {
	secret := "s"
	now := time.Now()
	c := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, _ := tkn.SignedString([]byte(secret))
	var out Claims
	if err := Validate(secret, s, &out); err == nil {
		t.Fatalf("expected not-before error")
	}
}
