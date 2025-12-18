package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/repository"
	"github.com/jjaenal/sisfo-akademik-backend/services/file-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/middleware"
	"go.uber.org/zap"
)

// @title           File Service API
// @version         1.0
// @description     File Service for Academic System
// @host      localhost:9098
// @BasePath  /api/v1
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
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Storage Setup
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./uploads"
	}
	baseURL := os.Getenv("APP_BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%d/uploads", cfg.HTTPPort)
	}

	storage, err := repository.NewLocalStorageProvider(storagePath, baseURL)
	if err != nil {
		log.Fatal("failed to init storage", zap.Error(err))
	}

	// Layers
	repo := repository.NewPostgresFileRepository(db)
	uc := usecase.NewFileUseCase(repo, storage)
	h := handler.NewFileHandler(uc)

	// Router
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	// Adapt shared net/http middleware to gin
	r.Use(toGinLogging(log))
	r.Use(toGinCORS(cfg.CORSAllowedOrigins))

	// Serve static files
	r.Static("/uploads", storagePath)

	// Health Check
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "up",
			"service": "file-service",
		})
	})

	// Routes
	// Protected routes need Auth middleware
	protected := r.Group("/")
	protected.Use(toGinAuthWith(cfg.JWTAccessSecret, cfg.JWTIssuer, cfg.JWTAudience))
	
	h.RegisterRoutes(protected)

	if err := r.Run(fmt.Sprintf(":%d", cfg.HTTPPort)); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}
}

// --- gin adapters for shared net/http middleware ---
func toGinLogging(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { c.Next() })
		h := middleware.Logging(l, next)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func toGinCORS(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { c.Next() })
		h := middleware.CORS(allowed, next)
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func toGinAuthWith(secret, issuer, audience string) gin.HandlerFunc {
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Propagate claims from request context into gin context for handlers
			if val := r.Context().Value(middleware.ClaimsKey); val != nil {
				c.Set("claims", val)
			}
			c.Next()
		})
		h := middleware.AuthWith(secret, issuer, audience, next)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
