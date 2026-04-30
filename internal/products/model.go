package products

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	CategoryID  uuid.UUID
}

type ProductVariation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid"`
	Price     float64
	Stock     int
}
