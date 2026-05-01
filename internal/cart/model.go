package cart

import "github.com/google/uuid"

type CartItem struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID             uuid.UUID `gorm:"type:uuid;index"`
	ProductVariationID uuid.UUID `gorm:"type:uuid"`
	Quantity           int
}
