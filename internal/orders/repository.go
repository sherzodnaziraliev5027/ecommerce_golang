package orders

import "gorm.io/gorm"

type Repository struct{}

func (r *Repository) CreateOrder(tx *gorm.DB, order *Order) error {
	return tx.Create(order).Error
}

func (r *Repository) CreateOrderItem(tx *gorm.DB, item *OrderItem) error {
	return tx.Create(item).Error
}
