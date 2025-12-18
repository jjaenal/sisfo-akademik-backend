package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/infrastructure/client"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/tracer"

	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// @title           Attendance Service API
// @version         1.0
// @description     Service for managing student and teacher attendance
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9093
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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
	tp, err := tracer.InitTracer("attendance-service", "http://jaeger:14268/api/traces")
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

	// Init Layers
	repo := postgres.NewStudentAttendanceRepository(dbPool)
	schoolService := client.NewSchoolService(cfg.AcademicServiceURL, 5*time.Second)
	uc := usecase.NewStudentAttendanceUseCase(repo, schoolService, 5*time.Second)
	h := handler.NewStudentAttendanceHandler(uc)

	teacherRepo := postgres.NewTeacherAttendanceRepository(dbPool)
	teacherUC := usecase.NewTeacherAttendanceUseCase(teacherRepo, 5*time.Second)
	teacherHandler := handler.NewTeacherAttendanceHandler(teacherUC)

	// Init Gin
	r := gin.Default()
	r.Use(otelgin.Middleware("attendance-service"))

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"service": "attendance-service",
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
	v1 := r.Group("/api/v1")
	attendance := v1.Group("/attendance")
	{
		attendance.POST("/students", h.Create)
		attendance.POST("/students/bulk", h.BulkCreate)
		attendance.GET("/students", h.GetByClassAndDate)
		attendance.GET("/students/:id", h.GetByID)
		attendance.GET("/students/:id/summary", h.GetSummary)
		attendance.PUT("/students/:id", h.Update)

		attendance.GET("/reports/daily", h.GetDailyReport)
		attendance.GET("/reports/monthly", h.GetMonthlyReport)
		attendance.GET("/reports/class/:class_id", h.GetClassReport)

		attendance.POST("/teachers/checkin", teacherHandler.CheckIn)
		attendance.PUT("/teachers/checkout", teacherHandler.CheckOut)
		attendance.GET("/teachers", teacherHandler.GetByTeacherAndDate)
	}

	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9093"
	}

	log.Info(fmt.Sprintf("starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err.Error())
	}
}
