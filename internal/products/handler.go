package products

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

// 🔥 Create product with variations
func (h *Handler) CreateProduct(c *gin.Context) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Variations  []struct {
			Price float64 `json:"price"`
			Stock int     `json:"stock"`
		} `json:"variations"`
	}

	// 1. Parse request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	// 2. Validate & parse category_id
	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid category_id"})
		return
	}

	// 3. Convert variations
	var variations []ProductVariation
	for _, v := range req.Variations {
		variations = append(variations, ProductVariation{
			Price: v.Price,
			Stock: v.Stock,
		})
	}

	// 4. Call service
	err = h.service.CreateProduct(
		req.Name,
		req.Description,
		categoryUUID,
		variations,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "product created",
	})
}

// 🔥 Get all products (with optional category filter)
func (h *Handler) GetAllProducts(c *gin.Context) {
	// 1. Read query param
	categoryIDStr := c.Query("category_id")

	// 2. Convert to UUID (if provided)
	var categoryUUID *uuid.UUID

	if categoryIDStr != "" {
		parsed, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid category_id"})
			return
		}
		categoryUUID = &parsed
	}

	// 3. Call service
	products, err := h.service.GetAllProducts(categoryUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 4. Return response
	c.JSON(200, products)
}

func (h *Handler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")

	// 1) parse UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product id"})
		return
	}

	// 2) call service
	product, err := h.service.GetProductByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 3) handle not found
	if product == nil {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}

	// 4) return
	c.JSON(200, product)
}
