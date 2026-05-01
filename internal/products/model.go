package products

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	CategoryID  uuid.UUID `gorm:"type:uuid;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductVariation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;index"`
	Price     float64
	Stock     int

	Product Product `gorm:"foreignKey:ProductID"`
}
