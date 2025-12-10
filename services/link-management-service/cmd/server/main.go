package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/config"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/database"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/handler"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/middleware"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/repository"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/service"
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
	repo := repository.NewLinkRepository(db)
	svc := service.NewLinkService(repo, cfg.CacheEvictionUrl)
	h := handler.NewLinkHandler(svc)

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "link-management-service",
		})
	})

	// Routes
	api := r.Group("/api/links")
	api.Use(middleware.AuthMiddleware(cfg)) // Apply Auth Middleware
	{
		api.POST("", h.CreateLink)
		api.GET("", h.GetMyLinks)
		api.DELETE("/:code", h.DeleteLink)
		api.PUT("/:code", h.UpdateLink)
	}

	// Start server
	log.Printf("Link Management Service starting on port %s", cfg.Port)
	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
