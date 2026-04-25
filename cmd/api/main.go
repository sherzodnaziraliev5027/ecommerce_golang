package main

import (
	"ecommerce/internal/users"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"ecommerce/pkg/database"

	"ecommerce/pkg/middleware"
)

func main() {
	// 🔹 Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	// 🔹 Connect DB

	database.Connect()

	err = database.DB.AutoMigrate(&users.User{})

	if err != nil {
		log.Fatal("migration failed:", err)
	}

	// 🔹 Init Gin
	r := gin.Default()

	// 🔹 Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 🔹 Setup Users Module
	repo := &users.Repository{}
	service := users.NewService(repo)
	handler := users.NewHandler(service)

	// 🔹 Public Routes
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// 🔐 Protected routes (require JWT)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/me", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		c.JSON(200, gin.H{
			"user_id": userID,
			"role":    role,
		})
	})

	// 🔹 Start server
	r.Run(":8080")

}
