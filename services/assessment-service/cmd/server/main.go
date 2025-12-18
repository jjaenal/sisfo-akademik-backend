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
	"github.com/jjaenal/sisfo-akademik-backend/services/assessment-service/internal/infrastructure/storage"
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

	// Init Repositories
	assessmentRepo := postgres.NewAssessmentRepository(dbPool)
	gradeRepo := postgres.NewGradeRepository(dbPool)
	gradeCategoryRepo := postgres.NewGradeCategoryRepository(dbPool)
	reportCardRepo := postgres.NewReportCardRepository(dbPool)
	templateRepo := postgres.NewTemplateRepository(dbPool)

	// Init Services
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:9094/files"
	}
	fileStorage := storage.NewLocalStorage(storagePath, baseURL)

	// Init UseCases
	gradingUseCase := usecase.NewGradingUseCase(assessmentRepo, gradeRepo, 10*time.Second)
	gradeCategoryUseCase := usecase.NewGradeCategoryUseCase(gradeCategoryRepo, 10*time.Second)
	reportCardUseCase := usecase.NewReportCardUseCase(reportCardRepo, gradeRepo, assessmentRepo, gradeCategoryRepo, fileStorage)
	templateUseCase := usecase.NewTemplateUseCase(templateRepo)

	// Init Handlers
	gradeHandler := handler.NewGradeHandler(gradingUseCase)
	gradeCategoryHandler := handler.NewGradeCategoryHandler(gradeCategoryUseCase)
	reportCardHandler := handler.NewReportCardHandler(reportCardUseCase)
	templateHandler := handler.NewTemplateHandler(templateUseCase)
	assessmentHandler := handler.NewAssessmentHandler(gradingUseCase)

	// Init Gin
	r := gin.Default()
	
	// Static file serving for local storage
	r.Static("/files", storagePath)

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

	api := r.Group("/api/v1")
	{
		// Assessment Routes
		assessments := api.Group("/assessments")
		{
			assessments.POST("", assessmentHandler.Create)
		}

		// Grade Routes
		grades := api.Group("/grades")
		{
			grades.POST("", gradeHandler.InputGrade)
			grades.GET("/student/:student_id", gradeHandler.GetStudentGrades)
			grades.GET("/final/:student_id", gradeHandler.CalculateFinalScore)
			grades.POST("/:id/approve", gradeHandler.ApproveGrade)
		}

		// Grade Category Routes
		categories := api.Group("/grade-categories")
		{
			categories.POST("", gradeCategoryHandler.Create)
			categories.GET("", gradeCategoryHandler.List)
			categories.PUT("/:id", gradeCategoryHandler.Update)
			categories.DELETE("/:id", gradeCategoryHandler.Delete)
		}

		// Report Card Routes
		reportCards := api.Group("/report-cards")
		{
			reportCards.POST("/generate", reportCardHandler.Generate)
			reportCards.GET("/student/:student_id", reportCardHandler.GetByStudent)
			reportCards.GET("/:id/download", reportCardHandler.GetPDF)
		}
		
		// Template Routes
		templates := api.Group("/templates")
		{
			templates.POST("", templateHandler.Create)
			templates.GET("", templateHandler.List)
			templates.GET("/:id", templateHandler.GetByID)
			templates.PUT("/:id", templateHandler.Update)
			templates.DELETE("/:id", templateHandler.Delete)
		}
	}

	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9094" // Default port for assessment-service
	}

	log.Info(fmt.Sprintf("starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err.Error())
	}
}
