package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID string    `json:"tenant_id"`
	Roles    []string  `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateAccess(secret string, ttl time.Duration, c Claims) (string, error) {
	now := time.Now()
	c.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(secret))
}

func GenerateRefresh(secret string, ttl time.Duration, subject string) (string, error) {
	now := time.Now()
	rc := jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	return token.SignedString([]byte(secret))
}

func Validate(secret string, tokenString string, out *Claims) error {
	tkn, err := jwt.ParseWithClaims(tokenString, out, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return err
	}
	if !tkn.Valid {
		return jwt.ErrTokenInvalidClaims
	}
	return nil
}
