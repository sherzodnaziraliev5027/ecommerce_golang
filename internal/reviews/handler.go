package reviews

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateReview(c *gin.Context) {

	var req struct {
		ProductID string `json:"product_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid body"})
		return
	}

	productUUID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product_id"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userIDStr.(string))

	err = h.service.CreateReview(userUUID, productUUID, req.Rating, req.Comment)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "review created"})
}

func (h *Handler) GetProductReviews(c *gin.Context) {

	productID := c.Param("id")

	productUUID, err := uuid.Parse(productID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product_id"})
		return
	}

	avg, reviews, err := h.service.GetProductReviewsWithAverage(productUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"average_rating": avg,
		"reviews":        reviews,
	})
}
