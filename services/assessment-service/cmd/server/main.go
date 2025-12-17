package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/handler"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/repository/postgres"
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/usecase"
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
	repo := postgres.NewGradeCategoryRepository(dbPool)
	uc := usecase.NewGradeCategoryUseCase(repo, 5*time.Second)
	h := handler.NewGradeCategoryHandler(uc)

	assessmentRepo := postgres.NewAssessmentRepository(dbPool)
	gradeRepo := postgres.NewGradeRepository(dbPool)
	gradingUC := usecase.NewGradingUseCase(assessmentRepo, gradeRepo, 5*time.Second)

	assessmentHandler := handler.NewAssessmentHandler(gradingUC)
	gradeHandler := handler.NewGradeHandler(gradingUC)

	reportCardRepo := postgres.NewReportCardRepository(dbPool)
	reportCardUC := usecase.NewReportCardUseCase(reportCardRepo, gradeRepo, assessmentRepo, repo)
	reportCardHandler := handler.NewReportCardHandler(reportCardUC)

	// Init Gin
	r := gin.Default()

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"service": "assessment-service",
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
	categories := v1.Group("/grade-categories")
	{
		categories.POST("", h.Create)
		categories.GET("", h.List)
		categories.GET("/:id", h.GetByID)
		categories.PUT("/:id", h.Update)
		categories.DELETE("/:id", h.Delete)
	}

	assessments := v1.Group("/assessments")
	{
		assessments.POST("", assessmentHandler.Create)
	}

	grades := v1.Group("/grades")
	{
		grades.POST("", gradeHandler.InputGrade)
		grades.GET("/students/:student_id", gradeHandler.GetStudentGrades)
		grades.GET("/students/:student_id/final-score", gradeHandler.CalculateFinalScore)
	}

	reportCards := v1.Group("/report-cards")
	{
		reportCards.POST("/generate", reportCardHandler.Generate)
		reportCards.GET("/student/:studentID/semester/:semesterID", reportCardHandler.GetByStudent)
		reportCards.GET("/:id", reportCardHandler.GetByID)
		reportCards.GET("/:id/pdf", reportCardHandler.GetPDF)
		reportCards.POST("/:id/publish", reportCardHandler.Publish)
	}

	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9092"
	}

	log.Info(fmt.Sprintf("starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err.Error())
	}
}
