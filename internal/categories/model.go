package categories

import "github.com/google/uuid"

type Category struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string
	ParentID *uuid.UUID
}
