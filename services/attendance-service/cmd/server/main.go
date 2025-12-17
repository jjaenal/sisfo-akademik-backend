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
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/attendance-service/internal/usecase"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/config"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/database"
	"github.com/jjaenal/sisfo-akademik-backend/shared/pkg/logger"
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

	// Connect to Database
	dbPool, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer dbPool.Close()
	log.Info("database connected")

	// Init Layers
	repo := postgres.NewStudentAttendanceRepository(dbPool)
	uc := usecase.NewStudentAttendanceUseCase(repo, 5*time.Second)
	h := handler.NewStudentAttendanceHandler(uc)

	teacherRepo := postgres.NewTeacherAttendanceRepository(dbPool)
	teacherUC := usecase.NewTeacherAttendanceUseCase(teacherRepo, 5*time.Second)
	teacherHandler := handler.NewTeacherAttendanceHandler(teacherUC)

	// Init Gin
	r := gin.Default()

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
		attendance.GET("/students", h.GetByClassAndDate)
		attendance.GET("/students/:id", h.GetByID)
		attendance.PUT("/students/:id", h.Update)

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
