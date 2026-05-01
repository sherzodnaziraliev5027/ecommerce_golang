package orders

import (
	"errors"

	"ecommerce/internal/cart"
	"ecommerce/internal/products"
	"ecommerce/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"math/rand"
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

type OrderResponse struct {
	ID     string              `json:"id"`
	Total  float64             `json:"total"`
	Status string              `json:"status"`
	Items  []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
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
			Status: "pending",
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

// func (s *Service) GetUserOrders(userID uuid.UUID) ([]OrderResponse, error) {

// 	orders, err := s.orderRepo.FindOrdersByUser(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result []OrderResponse

// 	for _, order := range orders {

// 		items, err := s.orderRepo.FindItemsByOrder(order.ID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		var itemResponses []OrderItemResponse

// 		for _, item := range items {

// 			// 🔥 get variation
// 			variation, err := s.productRepo.FindVariationByID(item.ProductVariationID)
// 			if err != nil {
// 				return nil, err
// 			}

// 			// 🔥 get product
// 			product, err := s.productRepo.FindProductByID(variation.ProductID)
// 			if err != nil {
// 				return nil, err
// 			}

// 			itemResponses = append(itemResponses, OrderItemResponse{
// 				ProductName: product.Name,
// 				Price:       item.Price,
// 				Quantity:    item.Quantity,
// 			})
// 		}

// 		result = append(result, OrderResponse{
// 			ID:    order.ID.String(),
// 			Total: order.Total,
// 			Items: itemResponses,
// 		})
// 	}

//		return result, nil
//	}
func (s *Service) GetUserOrders(userID uuid.UUID) ([]OrderResponse, error) {

	orders, err := s.orderRepo.FindOrdersWithItems(userID)
	if err != nil {
		return nil, err
	}

	var result []OrderResponse

	for _, order := range orders {

		var items []OrderItemResponse

		for _, item := range order.Items {

			items = append(items, OrderItemResponse{
				ProductName: item.ProductVariation.Product.Name,
				Price:       item.Price,
				Quantity:    item.Quantity,
			})
		}

		result = append(result, OrderResponse{
			ID:     order.ID.String(),
			Total:  order.Total,
			Status: order.Status,
			Items:  items,
		})
	}

	return result, nil
}

func (s *Service) PayOrder(userID, orderID uuid.UUID) (string, error) {

	var order Order

	err := database.DB.
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error

	if err != nil {
		return "", err
	}

	if order.Status != "pending" {
		return "", errors.New("order already processed")
	}

	// 🔥 simulate payment
	if rand.Intn(2) == 0 {
		order.Status = "failed"
	} else {
		order.Status = "paid"
	}

	err = database.DB.Save(&order).Error
	if err != nil {
		return "", err
	}

	return order.Status, nil
}

func (s *Service) UpdateOrderStatus(orderID uuid.UUID, newStatus string) error {

	var order Order

	err := database.DB.
		Where("id = ?", orderID).
		First(&order).Error

	if err != nil {
		return err
	}

	// 🔥 allowed transitions
	switch order.Status {

	case "paid":
		if newStatus != "shipped" {
			return errors.New("invalid transition")
		}

	case "shipped":
		if newStatus != "delivered" {
			return errors.New("invalid transition")
		}

	default:
		return errors.New("cannot update status from current state")
	}

	order.Status = newStatus

	return database.DB.Save(&order).Error
}
