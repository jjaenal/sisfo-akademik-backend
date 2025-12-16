package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/middleware"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	db, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		panic(err)
	}
	redis := redisutil.New(cfg.RedisAddr)
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))
	limiter := redisutil.NewLimiterFromCounter(redis.Raw())
	// Rate limit policy: READ=100/min, WRITE=30/min, AUTH prefix override=5/min
	r.Use(middleware.RateLimitByPolicy(limiter, 100, 30, map[string]int{"/api/v1/auth/": 5}))
	handler.NewHealthHandler(db, redis).Register(r)
	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	prRepo := repository.NewPasswordResetRepo(db)
	rb := rabbit.New(cfg.RabbitURL)
	phRepo := repository.NewPasswordHistoryRepo(db)
	// Audit middleware logs after response asynchronously
	r.Use(middleware.Audit(auditRepo))
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo, prRepo, rb, phRepo)
	authHandler.Register(r)
	if cfg.Env == "development" {
		handler.NewDevHandler(authRepo).Register(r)
	}
	// Protect routes
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTAccessSecret, cfg))
	authHandler.RegisterProtected(protected)
	// Users handlers
	usersUC := usecase.NewUsers(authRepo)
	usersHandler := handler.NewUsersHandler(usersUC)
	usersHandler.RegisterProtected(protected)
	// Roles handlers
	rolesRepo := repository.NewRolesRepo(db)
	rolesUC := usecase.NewRoles(authRepo, rolesRepo)
	rolesHandler := handler.NewRolesHandler(rolesUC)
	rolesHandler.RegisterProtected(protected)
	// Audit handlers
	auditHandler := handler.NewAuditHandler(auditRepo)
	auditHandler.RegisterProtected(protected)
	// Log retention job
	go func() {
		for {
			cutoff := time.Now().UTC().Add(-time.Duration(cfg.AuditRetentionDays) * 24 * time.Hour)
			_, _ = auditRepo.CleanupOlderThan(context.Background(), cutoff)
			time.Sleep(24 * time.Hour)
		}
	}()
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}
