package orders

import (
	"ecommerce/internal/products"
	"github.com/google/uuid"
)

type Order struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;index"`
	Total  float64
	Status string

	Items []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID            uuid.UUID `gorm:"type:uuid"`
	ProductVariationID uuid.UUID `gorm:"type:uuid"`
	Quantity           int
	Price              float64

	ProductVariation products.ProductVariation `gorm:"foreignKey:ProductVariationID"`
}
