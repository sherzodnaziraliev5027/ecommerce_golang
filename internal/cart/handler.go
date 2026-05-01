package cart

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

// 🔥 Add item to cart
func (h *Handler) AddToCart(c *gin.Context) {
	var req struct {
		VariationID string `json:"variation_id"`
		Quantity    int    `json:"quantity"`
	}

	// 1. Parse request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	// 2. Parse variation_id
	variationUUID, err := uuid.Parse(req.VariationID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid variation_id"})
		return
	}

	// 🔥 3. Get user ID from context (VERY IMPORTANT)
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

	// 4. Call service
	err = h.service.AddToCart(userUUID, variationUUID, req.Quantity)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "added to cart",
	})
}

// 🔥 Get user cart
func (h *Handler) GetCart(c *gin.Context) {

	// 1. Get user ID from context
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
	items, err := h.service.GetUserCart(userUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, items)
}

func (h *Handler) UpdateCart(c *gin.Context) {
	var req struct {
		VariationID string `json:"variation_id"`
		Quantity    int    `json:"quantity"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	variationUUID, err := uuid.Parse(req.VariationID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid variation_id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userIDStr.(string))

	err = h.service.UpdateQuantity(userUUID, variationUUID, req.Quantity)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "cart updated"})
}

func (h *Handler) RemoveFromCart(c *gin.Context) {
	variationID := c.Param("variation_id")

	variationUUID, err := uuid.Parse(variationID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid variation_id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userIDStr.(string))

	err = h.service.RemoveFromCart(userUUID, variationUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "item removed"})
}
