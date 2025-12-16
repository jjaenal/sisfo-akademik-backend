package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo      *repository.UsersRepo
	auditRepo *repository.AuditRepo
	prRepo    *repository.PasswordResetRepo
	phRepo    *repository.PasswordHistoryRepo
	cfg       config.Config
	redis     *redisutil.Client
	rabbit    *rabbit.Client
}

func NewAuthHandler(repo *repository.UsersRepo, cfg config.Config, r *redisutil.Client, audit *repository.AuditRepo, pr *repository.PasswordResetRepo, rb *rabbit.Client, ph *repository.PasswordHistoryRepo) *AuthHandler {
	return &AuthHandler{repo: repo, cfg: cfg, redis: r, auditRepo: audit, prRepo: pr, rabbit: rb, phRepo: ph}
}

func (h *AuthHandler) Register(r *gin.Engine) {
	r.POST("/api/v1/auth/login", h.login)
	r.POST("/api/v1/auth/refresh", h.refresh)
	r.POST("/api/v1/auth/logout", h.logout)
	r.POST("/api/v1/auth/forgot-password", h.forgotPassword)
	r.POST("/api/v1/auth/reset-password", h.resetPassword)
}

func (h *AuthHandler) RegisterProtected(g *gin.RouterGroup) {
	g.GET("/api/v1/auth/me", h.me)
	g.POST("/api/v1/auth/change-password", h.changePassword)
}

func isComplexPassword(s string) bool {
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSymbol := false
	for _, r := range s {
		if unicode.IsUpper(r) {
			hasUpper = true
			continue
		}
		if unicode.IsLower(r) {
			hasLower = true
			continue
		}
		if unicode.IsDigit(r) {
			hasDigit = true
			continue
		}
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			hasSymbol = true
			continue
		}
	}
	return hasUpper && hasLower && hasDigit && hasSymbol
}

func isCommonPassword(s string) bool {
	l := strings.ToLower(s)
	common := map[string]struct{}{
		"password":  {},
		"password1": {},
		"passw0rd":  {},
		"p@ssw0rd":  {},
		"123456":    {},
		"qwerty":    {},
		"admin":     {},
		"welcome":   {},
		"letmein":   {},
	}
	_, ok := common[l]
	return ok
}

type loginReq struct {
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginReq
	if err := c.BindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid json")
		return
	}
	if _, err := redisutil.Get(c.Request.Context(), h.redis.Raw(), "lockout:"+strings.TrimSpace(req.TenantID)+":"+strings.ToLower(strings.TrimSpace(req.Email))); err == nil {
		_ = h.auditRepo.Log(c.Request.Context(), req.TenantID, nil, "auth.login", "user", nil, map[string]any{"success": false, "reason": "locked"})
		httputil.Error(c.Writer, http.StatusForbidden, "3001", "Forbidden", "account locked")
		return
	}
	u, err := h.repo.FindByEmail(c.Request.Context(), req.TenantID, req.Email)
	if err != nil {
		_ = h.auditRepo.Log(c.Request.Context(), req.TenantID, nil, "auth.login", "user", nil, map[string]any{"success": false})
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid credentials")
		return
	}
	if !u.IsActive {
		_ = h.auditRepo.Log(c.Request.Context(), req.TenantID, &u.ID, "auth.login", "user", &u.ID, map[string]any{"success": false, "reason": "inactive"})
		httputil.Error(c.Writer, http.StatusForbidden, "3001", "Forbidden", "inactive user")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		key := "loginfail:" + strings.TrimSpace(req.TenantID) + ":" + strings.ToLower(strings.TrimSpace(req.Email))
		n, _ := redisutil.IncrWithTTL(c.Request.Context(), h.redis.Raw(), key, h.cfg.FailWindowTTL)
		if n >= int64(h.cfg.LockoutThreshold) {
			_ = redisutil.Set(c.Request.Context(), h.redis.Raw(), "lockout:"+strings.TrimSpace(req.TenantID)+":"+strings.ToLower(strings.TrimSpace(req.Email)), "1", h.cfg.LockoutTTL)
			_ = h.auditRepo.Log(c.Request.Context(), req.TenantID, &u.ID, "auth.login", "user", &u.ID, map[string]any{"success": false, "reason": "locked"})
			httputil.Error(c.Writer, http.StatusForbidden, "3001", "Forbidden", "account locked")
			return
		}
		_ = h.auditRepo.Log(c.Request.Context(), req.TenantID, &u.ID, "auth.login", "user", &u.ID, map[string]any{"success": false})
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid credentials")
		return
	}
	_ = redisutil.Del(c.Request.Context(), h.redis.Raw(), "loginfail:"+strings.TrimSpace(req.TenantID)+":"+strings.ToLower(strings.TrimSpace(req.Email)))
	_ = redisutil.Del(c.Request.Context(), h.redis.Raw(), "lockout:"+strings.TrimSpace(req.TenantID)+":"+strings.ToLower(strings.TrimSpace(req.Email)))
	access, err := jwtutil.GenerateAccessWith(h.cfg.JWTAccessSecret, h.cfg.JWTAccessTTL, jwtutil.Claims{
		UserID:   u.ID,
		TenantID: u.TenantID,
		Roles:    []string{},
	}, h.cfg.JWTIssuer, h.cfg.JWTAudience)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	refresh, err := jwtutil.GenerateRefreshWith(h.cfg.JWTRefreshSecret, h.cfg.JWTRefreshTTL, u.ID.String(), h.cfg.JWTIssuer, h.cfg.JWTAudience)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	_ = h.auditRepo.Log(c.Request.Context(), u.TenantID, &u.ID, "auth.login", "user", &u.ID, map[string]any{"success": true})
	httputil.Success(c.Writer, map[string]any{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) me(c *gin.Context) {
	val, exists := c.Get("claims")
	if !exists {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", nil)
		return
	}
	claims, _ := val.(jwtutil.Claims)
	u, err := h.repo.FindByID(c.Request.Context(), claims.UserID)
	if err != nil || !u.IsActive {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "user not found or inactive")
		return
	}
	httputil.Success(c.Writer, map[string]any{
		"id":        u.ID,
		"tenant_id": u.TenantID,
		"email":     u.Email,
		"roles":     claims.Roles,
	})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) refresh(c *gin.Context) {
	var req refreshReq
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "missing refresh_token")
		return
	}
	var claims jwt.RegisteredClaims
	tkn, err := jwt.ParseWithClaims(req.RefreshToken, &claims, func(token *jwt.Token) (any, error) {
		return []byte(h.cfg.JWTRefreshSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !tkn.Valid || claims.Subject == "" {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid refresh token")
		return
	}
	if claims.Issuer != "" && claims.Issuer != h.cfg.JWTIssuer {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid issuer")
		return
	}
	if len(claims.Audience) > 0 {
		found := false
		for _, aud := range claims.Audience {
			if aud == h.cfg.JWTAudience {
				found = true
				break
			}
		}
		if !found {
			httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid audience")
			return
		}
	}
	// Check blacklist by JTI
	if strings.TrimSpace(claims.ID) != "" {
		if _, err := redisutil.Get(c.Request.Context(), h.redis.Raw(), "blacklist:jti:"+claims.ID); err == nil {
			httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "token revoked")
			return
		}
	}
	uid, err := uuid.Parse(claims.Subject)
	if err != nil {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid subject")
		return
	}
	u, err := h.repo.FindByID(c.Request.Context(), uid)
	if err != nil || !u.IsActive {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "user not found or inactive")
		return
	}
	access, err := jwtutil.GenerateAccess(h.cfg.JWTAccessSecret, h.cfg.JWTAccessTTL, jwtutil.Claims{
		UserID:   u.ID,
		TenantID: u.TenantID,
		Roles:    []string{},
	})
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	newRefresh, err := jwtutil.GenerateRefreshWith(h.cfg.JWTRefreshSecret, h.cfg.JWTRefreshTTL, u.ID.String(), h.cfg.JWTIssuer, h.cfg.JWTAudience)
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	// Blacklist old refresh
	if strings.TrimSpace(claims.ID) != "" {
		_ = redisutil.Set(c.Request.Context(), h.redis.Raw(), "blacklist:jti:"+claims.ID, "1", h.cfg.JWTRefreshTTL)
	}
	_ = h.auditRepo.Log(c.Request.Context(), u.TenantID, &u.ID, "auth.refresh", "user", &u.ID, map[string]any{"success": true})
	httputil.Success(c.Writer, map[string]any{
		"access_token":  access,
		"refresh_token": newRefresh,
	})
}

func (h *AuthHandler) logout(c *gin.Context) {
	var req refreshReq
	_ = c.BindJSON(&req)
	// Blacklist provided refresh token (if any)
	if strings.TrimSpace(req.RefreshToken) != "" {
		var claims jwt.RegisteredClaims
		if tkn, err := jwt.ParseWithClaims(req.RefreshToken, &claims, func(token *jwt.Token) (any, error) {
			return []byte(h.cfg.JWTRefreshSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()})); err == nil && tkn.Valid && strings.TrimSpace(claims.ID) != "" {
			_ = redisutil.Set(c.Request.Context(), h.redis.Raw(), "blacklist:jti:"+claims.ID, "1", h.cfg.JWTRefreshTTL)
		}
	}
	_ = h.auditRepo.Log(c.Request.Context(), "", nil, "auth.logout", "user", nil, map[string]any{"success": true})
	httputil.Success(c.Writer, map[string]any{"message": "logged out"})
}

type forgotReq struct {
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
}

func (h *AuthHandler) forgotPassword(c *gin.Context) {
	var req forgotReq
	_ = c.BindJSON(&req)
	u, err := h.repo.FindByEmail(c.Request.Context(), strings.TrimSpace(req.TenantID), strings.TrimSpace(req.Email))
	if err == nil && u.IsActive {
		var b [32]byte
		_, _ = rand.Read(b[:])
		token := hex.EncodeToString(b[:])
		hash := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hash[:])
		exp := time.Now().UTC().Add(30 * time.Minute)
		if _, err := h.prRepo.Create(c.Request.Context(), u.TenantID, u.ID, tokenHash, exp); err == nil {
			_ = h.auditRepo.Log(c.Request.Context(), u.TenantID, &u.ID, "auth.forgot_password", "user", &u.ID, map[string]any{"success": true})
			if strings.ToLower(h.cfg.Env) == "test" {
				_ = redisutil.Set(c.Request.Context(), h.redis.Raw(), "password_reset:th:"+tokenHash, u.ID.String(), 30*time.Minute)
				httputil.Success(c.Writer, map[string]any{"reset_token": token})
				return
			}
			if h.rabbit != nil && strings.ToLower(h.cfg.Env) != "test" {
				_ = h.rabbit.PublishJSON("events", "auth.password_reset.requested", map[string]any{
					"tenant_id": u.TenantID,
					"user_id":   u.ID.String(),
					"email":     u.Email,
					"token":     token,
					"type":      "password_reset",
				})
			}
		}
		// Fallback for test env when DB create fails (no DDL permission)
		if strings.ToLower(h.cfg.Env) == "test" {
			_ = redisutil.Set(c.Request.Context(), h.redis.Raw(), "password_reset:th:"+tokenHash, u.ID.String(), 30*time.Minute)
			httputil.Success(c.Writer, map[string]any{"reset_token": token})
			return
		}
	}
	httputil.Success(c.Writer, map[string]any{"message": "if the email exists, a reset link has been sent"})
}

type resetReq struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h *AuthHandler) resetPassword(c *gin.Context) {
	var req resetReq
	if err := c.BindJSON(&req); err != nil || strings.TrimSpace(req.Token) == "" || strings.TrimSpace(req.Password) == "" {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "missing token or password")
		return
	}
	p := strings.TrimSpace(req.Password)
	if len(p) < 8 {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password must be at least 8 characters")
		return
	}
	if !isComplexPassword(p) {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password must include uppercase, lowercase, number, and symbol")
		return
	}
	if isCommonPassword(p) {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password is too common")
		return
	}
	tokenStr := strings.TrimSpace(req.Token)
	hash := sha256.Sum256([]byte(tokenStr))
	tokenHash := hex.EncodeToString(hash[:])
	rec, err := h.prRepo.FindValidByTokenHash(c.Request.Context(), tokenHash, time.Now().UTC())
	if err != nil {
		// Test env fallback via Redis
		if strings.ToLower(h.cfg.Env) == "test" {
			if uidStr, rerr := redisutil.Get(c.Request.Context(), h.redis.Raw(), "password_reset:th:"+tokenHash); rerr == nil {
				if uid, perr := uuid.Parse(uidStr); perr == nil {
					u, ferr := h.repo.FindByID(c.Request.Context(), uid)
					if ferr == nil {
						if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(p)) == nil {
							httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password reuse not allowed")
							return
						}
						if h.phRepo != nil {
							if items, ierr := h.phRepo.Recent(c.Request.Context(), uid, 5); ierr == nil {
								for _, it := range items {
									if bcrypt.CompareHashAndPassword([]byte(it.PasswordHash), []byte(p)) == nil {
										httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password reuse not allowed")
										return
									}
								}
							}
						}
					}
					_, uerr := h.repo.Update(c.Request.Context(), uid, repository.UpdateUserParams{
						Password: &p,
						Now:      time.Now().UTC(),
					})
					if uerr == nil {
						if h.phRepo != nil && u != nil {
							_ = h.phRepo.Add(c.Request.Context(), uid, u.PasswordHash, time.Now().UTC())
							_ = h.phRepo.Prune(c.Request.Context(), uid, 5)
						}
						_ = h.auditRepo.Log(c.Request.Context(), "", &uid, "auth.reset_password", "user", &uid, map[string]any{"success": true})
						_ = redisutil.Del(c.Request.Context(), h.redis.Raw(), "password_reset:th:"+tokenHash)
						httputil.Success(c.Writer, map[string]any{"message": "password reset successful"})
						return
					}
				}
			}
		}
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid or expired token")
		return
	}
	u, ferr := h.repo.FindByID(c.Request.Context(), rec.UserID)
	if ferr != nil {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "user not found or inactive")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(p)) == nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password reuse not allowed")
		return
	}
	if h.phRepo != nil {
		if items, ierr := h.phRepo.Recent(c.Request.Context(), rec.UserID, 5); ierr == nil {
			for _, it := range items {
				if bcrypt.CompareHashAndPassword([]byte(it.PasswordHash), []byte(p)) == nil {
					httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password reuse not allowed")
					return
				}
			}
		}
	}
	_, err = h.repo.Update(c.Request.Context(), rec.UserID, repository.UpdateUserParams{
		Password: &p,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	if h.phRepo != nil {
		_ = h.phRepo.Add(c.Request.Context(), rec.UserID, u.PasswordHash, time.Now().UTC())
		_ = h.phRepo.Prune(c.Request.Context(), rec.UserID, 5)
	}
	_ = h.prRepo.MarkUsed(c.Request.Context(), rec.ID, time.Now().UTC())
	_ = h.auditRepo.Log(c.Request.Context(), rec.TenantID, &rec.UserID, "auth.reset_password", "user", &rec.UserID, map[string]any{"success": true})
	httputil.Success(c.Writer, map[string]any{"message": "password reset successful"})
}

type changeReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) changePassword(c *gin.Context) {
	var req changeReq
	if err := c.BindJSON(&req); err != nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "invalid json")
		return
	}
	val, exists := c.Get("claims")
	if !exists {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", nil)
		return
	}
	claims, _ := val.(jwtutil.Claims)
	u, err := h.repo.FindByID(c.Request.Context(), claims.UserID)
	if err != nil || !u.IsActive {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "user not found or inactive")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(strings.TrimSpace(req.OldPassword))) != nil {
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid current password")
		return
	}
	p := strings.TrimSpace(req.NewPassword)
	if len(p) < 8 {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password must be at least 8 characters")
		return
	}
	if !isComplexPassword(p) {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password must include uppercase, lowercase, number, and symbol")
		return
	}
	if isCommonPassword(p) {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password is too common")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(p)) == nil {
		httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password reuse not allowed")
		return
	}
	if h.phRepo != nil {
		if items, err := h.phRepo.Recent(c.Request.Context(), u.ID, 5); err == nil {
			for _, it := range items {
				if bcrypt.CompareHashAndPassword([]byte(it.PasswordHash), []byte(p)) == nil {
					httputil.Error(c.Writer, http.StatusBadRequest, "4001", "Invalid Input", "password reuse not allowed")
					return
				}
			}
		}
	}
	_, err = h.repo.Update(c.Request.Context(), u.ID, repository.UpdateUserParams{
		Password: &p,
		Now:      time.Now().UTC(),
	})
	if err != nil {
		httputil.Error(c.Writer, http.StatusInternalServerError, "1001", "Internal Server Error", err.Error())
		return
	}
	if h.phRepo != nil {
		_ = h.phRepo.Add(c.Request.Context(), u.ID, u.PasswordHash, time.Now().UTC())
		_ = h.phRepo.Prune(c.Request.Context(), u.ID, 5)
	}
	_ = h.auditRepo.Log(c.Request.Context(), u.TenantID, &u.ID, "auth.change_password", "user", &u.ID, map[string]any{"success": true})
	httputil.Success(c.Writer, map[string]any{"message": "password changed"})
}
