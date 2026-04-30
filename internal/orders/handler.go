package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// 🔥 Checkout
func (h *Handler) Checkout(c *gin.Context) {

	// 1. Get user from JWT
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid user_id"})
		return
	}

	// 2. Call service
	err = h.service.Checkout(userUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "order created successfully",
	})
}
