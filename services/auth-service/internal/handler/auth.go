package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/httputil"
	jwtutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/jwt"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo      *repository.UsersRepo
	auditRepo *repository.AuditRepo
	cfg       config.Config
	redis     *redisutil.Client
}

func NewAuthHandler(repo *repository.UsersRepo, cfg config.Config, r *redisutil.Client, audit *repository.AuditRepo) *AuthHandler {
	return &AuthHandler{repo: repo, cfg: cfg, redis: r, auditRepo: audit}
}

func (h *AuthHandler) Register(r *gin.Engine) {
	r.POST("/api/v1/auth/login", h.login)
	r.POST("/api/v1/auth/refresh", h.refresh)
	r.POST("/api/v1/auth/logout", h.logout)
}

func (h *AuthHandler) RegisterProtected(g *gin.RouterGroup) {
	g.GET("/api/v1/auth/me", h.me)
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
		_ = h.auditRepo.Log(c.Request.Context(), req.TenantID, &u.ID, "auth.login", "user", &u.ID, map[string]any{"success": false})
		httputil.Error(c.Writer, http.StatusUnauthorized, "2001", "Unauthorized", "invalid credentials")
		return
	}
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
