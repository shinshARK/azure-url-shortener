package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/config"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/database"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/handler"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/repository"
	"github.com/shinshark/azure-url-shortener/services/auth-service/internal/service"
)

func main() {
	// Load Config
	cfg := config.LoadConfig()

	// Connect to Database
	db, err := database.Connect(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Layers
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, cfg)
	h := handler.NewAuthHandler(svc)

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "auth-service",
		})
	})

	// Auth Routes
	api := r.Group("/api/auth")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
	}

	// Start server
	log.Printf("Auth Service starting on port %s", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
