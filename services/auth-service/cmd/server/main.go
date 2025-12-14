package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/middleware"
	"github.com/jjaenal/sisfo-akademik-backend/services/auth-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
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
	// Rate limit applied globally using shared config
	r.Use(middleware.RateLimit(limiter, cfg.RateLimitPerMinute))
	handler.NewHealthHandler(db, redis).Register(r)
	authRepo := repository.NewUsersRepo(db)
	auditRepo := repository.NewAuditRepo(db)
	authHandler := handler.NewAuthHandler(authRepo, cfg, redis, auditRepo)
	authHandler.Register(r)
	// Protect routes
	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTAccessSecret, cfg))
	authHandler.RegisterProtected(protected)
	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	if err := r.Run(addr); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}
