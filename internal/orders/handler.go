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

func (h *Handler) GetOrders(c *gin.Context) {

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

	orders, err := h.service.GetUserOrders(userUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, orders)
}

func (h *Handler) PayOrder(c *gin.Context) {

	orderIDStr := c.Param("order_id")

	orderUUID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid order_id"})
		return
	}

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

	// ✅ IMPORTANT: TWO RETURNS
	status, err := h.service.PayOrder(userUUID, orderUUID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if status == "failed" {
		c.JSON(400, gin.H{"error": "payment failed"})
		return
	}

	c.JSON(200, gin.H{
		"message": "payment successful",
		"status":  status,
	})
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {

	// 🔥 check role
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	orderIDStr := c.Param("order_id")
	status := c.Query("status")

	orderUUID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid order_id"})
		return
	}

	err = h.service.UpdateOrderStatus(orderUUID, status)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "status updated"})
}
