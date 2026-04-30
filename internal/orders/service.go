package orders

import (
	"errors"

	"ecommerce/internal/cart"
	"ecommerce/internal/products"
	"ecommerce/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	orderRepo   *Repository
	cartRepo    *cart.Repository
	productRepo *products.Repository
}

func NewService(o *Repository, c *cart.Repository, p *products.Repository) *Service {
	return &Service{
		orderRepo:   o,
		cartRepo:    c,
		productRepo: p,
	}
}

func (s *Service) Checkout(userID uuid.UUID) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		// 1. Get cart items
		items, err := s.cartRepo.FindByUser(userID)
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return errors.New("cart is empty")
		}

		// 2. Validate stock
		for _, item := range items {
			variation, err := s.productRepo.FindVariationByID(item.ProductVariationID)
			if err != nil {
				return err
			}

			if variation.Stock < item.Quantity {
				return errors.New("insufficient stock")
			}
		}

		// 3. Create order
		order := &Order{
			ID:     uuid.New(),
			UserID: userID,
			Total:  0,
		}

		if err := s.orderRepo.CreateOrder(tx, order); err != nil {
			return err
		}

		var total float64

		// 4. Create order items + deduct stock
		for _, item := range items {

			variation, err := s.productRepo.FindVariationByID(item.ProductVariationID)
			if err != nil {
				return err
			}

			orderItem := &OrderItem{
				ID:                 uuid.New(),
				OrderID:            order.ID,
				ProductVariationID: variation.ID,
				Quantity:           item.Quantity,
				Price:              variation.Price,
			}

			if err := s.orderRepo.CreateOrderItem(tx, orderItem); err != nil {
				return err
			}

			total += variation.Price * float64(item.Quantity)

			// deduct stock
			variation.Stock -= item.Quantity
			if err := tx.Save(variation).Error; err != nil {
				return err
			}
		}

		// 5. Update total
		order.Total = total
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		// 6. Clear cart
		if err := tx.Where("user_id = ?", userID).Delete(&cart.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})
}
