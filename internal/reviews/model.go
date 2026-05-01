package reviews

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_product"`
	ProductID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_product"`
	Rating    int
	Comment   string
	CreatedAt time.Time
}
