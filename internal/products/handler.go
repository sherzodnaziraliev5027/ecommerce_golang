package products

import (
	"log"
	"net/http"
	"strconv"

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
		log.Println("GetAllProducts error:", err)

		c.JSON(500, gin.H{"error": "internal server error"})
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

func (h *Handler) GetAllProducts(c *gin.Context) {

	categoryIDStr := c.Query("category_id")

	// ✅ filtering params
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	sort := c.Query("sort")

	// 🔥 ADDED SORT HERE
	// 🔥 default sort
	if sort == "" {
		sort = "price_asc"
	}

	if sort == "price_desc" {
		sort = "price_desc"
	}
	if sort != "" && sort != "price_asc" && sort != "price_desc" {
		c.JSON(400, gin.H{"error": "invalid sort value"})
		return
	}

	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid min_price"})
			return
		}
	}

	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid max_price"})
			return
		}
	}

	// ✅ pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(400, gin.H{"error": "invalid page"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(400, gin.H{"error": "invalid limit"})
		return
	}

	if limit > 50 {
		limit = 50
	}

	offset := (page - 1) * limit

	// ✅ parse category
	var categoryUUID *uuid.UUID
	if categoryIDStr != "" {
		parsed, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid category_id"})
			return
		}
		categoryUUID = &parsed
	}

	products, err := h.service.GetAllProducts(
		categoryUUID,
		minPrice,
		maxPrice,
		sort,
		limit,
		offset,
		page,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, products)
}
