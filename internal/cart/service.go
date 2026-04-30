package cart

import (
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

type CartResponse struct {
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Total       float64 `json:"total"`
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) AddToCart(userID uuid.UUID, variationID uuid.UUID, quantity int) error {

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// 1. Check if item already exists
	existing, err := s.repo.FindByUserAndVariation(userID, variationID)

	if err == nil && existing != nil {
		// 🔥 Item exists → increase quantity
		existing.Quantity += quantity
		return s.repo.Update(existing)
	}

	// 2. Item does not exist → create new
	item := &CartItem{
		ID:                 uuid.New(),
		UserID:             userID,
		ProductVariationID: variationID,
		Quantity:           quantity,
	}

	return s.repo.Create(item)
}

//	func (s *Service) GetUserCart(userID uuid.UUID) ([]CartItem, error) {
//		return s.repo.FindByUser(userID)
//	}
func (s *Service) GetUserCart(userID uuid.UUID) ([]CartResponse, error) {

	items, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	var result []CartResponse

	for _, item := range items {

		// 1. Get variation
		variation, err := s.repo.FindVariationByID(item.ProductVariationID)
		if err != nil {
			return nil, err
		}

		// 2. Get product
		product, err := s.repo.FindProductByID(variation.ProductID)
		if err != nil {
			return nil, err
		}

		// 3. Build response
		result = append(result, CartResponse{
			ProductName: product.Name,
			Price:       variation.Price,
			Quantity:    item.Quantity,
			Total:       variation.Price * float64(item.Quantity),
		})
	}

	return result, nil
}

func (s *Service) UpdateQuantity(userID, variationID uuid.UUID, quantity int) error {

	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	item, err := s.repo.FindByUserAndVariation(userID, variationID)
	if err != nil {
		return err
	}

	item.Quantity = quantity

	return s.repo.Update(item)
}

func (s *Service) RemoveFromCart(userID, variationID uuid.UUID) error {

	item, err := s.repo.FindByUserAndVariation(userID, variationID)
	if err != nil {
		return err
	}

	return s.repo.Delete(item)
}
