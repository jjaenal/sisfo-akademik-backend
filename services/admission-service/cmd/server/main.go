package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/middleware"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/admission-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/rabbit"
	redisutil "github.com/jjaenal/sisfo-akademik-backend/shared/pkg/redis"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/tracer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	log.Info("config loaded")

	// Tracer
	tp, err := tracer.InitTracer("admission-service", cfg.JaegerEndpoint)
	if err != nil {
		log.Fatal("failed to init tracer", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal("failed to shutdown tracer", zap.Error(err))
		}
	}()

	// Connect to Database
	dbPool, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer dbPool.Close()
	log.Info("database connected")

	// Connect to RabbitMQ
	rabbitClient := rabbit.New(cfg.RabbitURL)
	if rabbitClient != nil {
		defer rabbitClient.Close()
		log.Info("rabbitmq connected")
	} else {
		log.Warn("rabbitmq connection failed, event publishing will be disabled")
	}

	// Init Layers
	appRepo := postgres.NewApplicationRepository(dbPool)
	repo := postgres.NewAdmissionPeriodRepository(dbPool)
	uc := usecase.NewAdmissionPeriodUseCase(repo, appRepo, 5*time.Second)
	h := handler.NewAdmissionPeriodHandler(uc)

	appUC := usecase.NewApplicationUseCase(appRepo, rabbitClient, 5*time.Second)
	appHandler := handler.NewApplicationHandler(appUC)

	docRepo := postgres.NewApplicationDocumentRepository(dbPool)
	docUC := usecase.NewApplicationDocumentUseCase(docRepo, appRepo, "./uploads")
	docHandler := handler.NewApplicationDocumentHandler(docUC)

	// Init Gin
	r := gin.Default()
	r.Use(otelgin.Middleware("admission-service"))

	redis := redisutil.New(cfg.RedisAddr)
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.RateLimitByPolicy(redisutil.NewLimiterFromCounter(redis.Raw()), 100, 30, nil))


	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"service": "admission-service",
				"status":  "healthy",
				"db":      "connected",
			},
			"meta": gin.H{
				"timestamp":  time.Now().UTC(),
				"request_id": uuid.NewString(),
			},
		})
	})

	// Routes
	v1 := r.Group("/api/v1/admission")
	periods := v1.Group("/periods")
	{
		periods.POST("", h.Create)
		periods.GET("", h.List)
		periods.GET("/active", h.GetActive)
		periods.GET("/:id", h.GetByID)
		periods.PUT("/:id", h.Update)
		periods.DELETE("/:id", h.Delete)
		periods.POST("/:id/calculate-final-scores", appHandler.CalculateFinalScores)
		periods.POST("/:id/announce", h.AnnounceResults)
	}

	applications := v1.Group("/applications")
	{
		applications.POST("", appHandler.Submit)
		applications.GET("/status", appHandler.GetStatus)
		applications.GET("", appHandler.List)
		applications.PUT("/:id/verify", appHandler.Verify)
		applications.POST("/:id/test-score", appHandler.InputTestScore)
		applications.POST("/:id/interview-score", appHandler.InputInterviewScore)
		applications.POST("/:id/register", appHandler.Register)

		// Document routes
		applications.POST("/:id/documents", docHandler.Upload)
		applications.GET("/:id/documents", docHandler.GetByApplicationID)
	}

	documents := v1.Group("/documents")
	{
		documents.DELETE("/:id", docHandler.Delete)
	}

	// Prometheus metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9095"
	}

	log.Info(fmt.Sprintf("starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err.Error())
	}
}
