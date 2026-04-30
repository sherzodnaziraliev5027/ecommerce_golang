package categories

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

// 🔹 Create category
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Name     string  `json:"name"`
		ParentID *string `json:"parent_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid input"})
		return
	}

	var parentUUID *uuid.UUID

	if req.ParentID != nil {
		parsed, err := uuid.Parse(*req.ParentID)
		if err == nil {
			parentUUID = &parsed
		}
	}

	err := h.service.Create(req.Name, parentUUID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "category created",
	})
}

// 🔹 Get all categories
func (h *Handler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, categories)

}

func (h *Handler) GetTree(c *gin.Context) {
	tree, err := h.service.GetTree()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tree)
}
