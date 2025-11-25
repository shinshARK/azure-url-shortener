package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/models"
	"github.com/shinshark/azure-url-shortener/services/link-management-service/internal/service"
)

type LinkHandler struct {
	Service *service.LinkService
}

func NewLinkHandler(svc *service.LinkService) *LinkHandler {
	return &LinkHandler{Service: svc}
}

func (h *LinkHandler) CreateLink(c *gin.Context) {
	var req models.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get User Info from Context (set by AuthMiddleware)
	var userID *int
	role := "Guest"

	if id, exists := c.Get("userID"); exists {
		uid := id.(int)
		userID = &uid
	}
	if r, exists := c.Get("role"); exists {
		role = r.(string)
	}

	link, err := h.Service.CreateLink(&req, userID, role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, link)
}

func (h *LinkHandler) GetMyLinks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	links, err := h.Service.GetUserLinks(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, links)
}

func (h *LinkHandler) DeleteLink(c *gin.Context) {
	shortCode := c.Param("code")
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	
	role := c.GetString("role")

	err := h.Service.DeleteLink(shortCode, userID.(int), role)
	if err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this link"})
			return
		}
		if err.Error() == "link not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted"})
}
