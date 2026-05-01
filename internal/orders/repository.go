package orders

import (
	"ecommerce/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct{}

func (r *Repository) CreateOrder(tx *gorm.DB, order *Order) error {
	return tx.Create(order).Error
}

func (r *Repository) CreateOrderItem(tx *gorm.DB, item *OrderItem) error {
	return tx.Create(item).Error
}

func (r *Repository) FindOrdersWithItems(userID uuid.UUID) ([]Order, error) {
	var orders []Order

	err := database.DB.
		Preload("Items").
		Preload("Items.ProductVariation").
		Preload("Items.ProductVariation.Product").
		Where("user_id = ?", userID).
		Find(&orders).Error

	return orders, err
}

func NewRepository() *Repository {
	return &Repository{}
}
