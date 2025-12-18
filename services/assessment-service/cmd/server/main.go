package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	port := os.Getenv("APP_HTTP_PORT")
	if port == "" {
		port = "9094" // Default port for assessment-service
	}

	log.Info(fmt.Sprintf("starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err.Error())
	}
}
