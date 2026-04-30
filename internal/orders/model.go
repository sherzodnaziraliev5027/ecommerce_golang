package orders

import "github.com/google/uuid"

type Order struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid"`
	Total  float64
}

type OrderItem struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID            uuid.UUID `gorm:"type:uuid"`
	ProductVariationID uuid.UUID `gorm:"type:uuid"`
	Quantity           int
	Price              float64
}
