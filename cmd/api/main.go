package main

import (
	"ecommerce/internal/categories"
	"ecommerce/internal/reviews"
	"ecommerce/internal/users"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"ecommerce/pkg/database"

	"ecommerce/pkg/middleware"

	"ecommerce/internal/cart"

	"ecommerce/internal/orders"

	"ecommerce/internal/products"
)

func main() {
	// 🔹 Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	// 🔹 Connect DB

	database.Connect()

	err = database.DB.AutoMigrate(&users.User{}, &categories.Category{},
		&products.Product{}, &products.ProductVariation{}, &cart.CartItem{}, &orders.Order{}, &orders.OrderItem{}, &reviews.Review{})

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
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("super_admin"))

	admin.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome SuperAdmin",
		})
	})

	catRepo := &categories.Repository{}
	catService := categories.NewService(catRepo)
	catHandler := categories.NewHandler(catService)

	prodRepo := &products.Repository{}
	prodService := products.NewService(prodRepo)
	prodHandler := products.NewHandler(prodService)

	cartRepo := &cart.Repository{}
	cartService := cart.NewService(cartRepo)
	cartHandler := cart.NewHandler(cartService)

	orderRepo := &orders.Repository{}
	orderService := orders.NewService(orderRepo, cartRepo, prodRepo)
	orderHandler := orders.NewHandler(orderService)

	reviewRepo := reviews.NewRepository()

	reviewService := reviews.NewService(reviewRepo, orderRepo)
	reviewHandler := reviews.NewHandler(reviewService)

	protected.POST("/categories", catHandler.Create)
	protected.POST("/products", prodHandler.CreateProduct)
	protected.POST("/cart", cartHandler.AddToCart)
	protected.POST("/checkout", orderHandler.Checkout)
	protected.POST("/pay/:order_id", orderHandler.PayOrder)

	protected.GET("/cart", cartHandler.GetCart)

	protected.GET("/categories", catHandler.GetAll)
	protected.GET("/categories/tree", catHandler.GetTree)

	protected.GET("/products", prodHandler.GetAllProducts)

	protected.GET("/products/:id", prodHandler.GetProductByID)
	protected.GET("/orders", orderHandler.GetOrders)

	protected.PUT("/cart", cartHandler.UpdateCart)
	protected.DELETE("/cart/:variation_id", cartHandler.RemoveFromCart)

	protected.PUT("/orders/:order_id/status", orderHandler.UpdateOrderStatus)

	protected.POST("/reviews", reviewHandler.CreateReview)
	protected.GET("/products/:id/reviews", reviewHandler.GetProductReviews)

	// 🔹 Start server
	r.Run(":8080")

}
