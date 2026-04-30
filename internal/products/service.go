package products

import (
	"ecommerce/pkg/database"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"errors"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

type ProductResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	CategoryID  string              `json:"category_id"`
	Variations  []VariationResponse `json:"variations"`
}

type VariationResponse struct {
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// 🔥 Create product with variations
func (s *Service) CreateProduct(
	name string,
	description string,
	categoryID uuid.UUID,
	variations []ProductVariation,
) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		// 1. Create product
		product := &Product{
			ID:          uuid.New(),
			Name:        name,
			Description: description,
			CategoryID:  categoryID,
		}

		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// 2. Create variations
		for i := range variations {
			variations[i].ID = uuid.New()
			variations[i].ProductID = product.ID

			if err := tx.Create(&variations[i]).Error; err != nil {
				return err
			}
		}

		// 2. Create variations
		for i := range variations {
			variations[i].ID = uuid.New()
			variations[i].ProductID = product.ID

			if err := tx.Create(&variations[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) GetAllProducts(categoryID *uuid.UUID) ([]ProductResponse, error) {

	products, err := s.repo.FindAllProducts()
	if err != nil {
		return nil, err
	}

	if categoryID != nil {
		filtered := []Product{}

		for _, p := range products {
			if p.CategoryID == *categoryID {
				filtered = append(filtered, p)
			}
		}

		products = filtered
	}

	variations, err := s.repo.FindAllVariations()
	if err != nil {
		return nil, err
	}

	// map productID → variations
	variationMap := make(map[string][]VariationResponse)

	for _, v := range variations {
		pid := v.ProductID.String()

		variationMap[pid] = append(variationMap[pid], VariationResponse{
			Price: v.Price,
			Stock: v.Stock,
		})
	}

	// build final response
	var result []ProductResponse

	for _, p := range products {
		id := p.ID.String()

		result = append(result, ProductResponse{
			ID:          id,
			Name:        p.Name,
			Description: p.Description,
			CategoryID:  p.CategoryID.String(),
			Variations:  variationMap[id],
		})
	}

	return result, nil
}

func (s *Service) GetProductByID(id uuid.UUID) (*ProductResponse, error) {

	product, err := s.repo.FindProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // not found
		}
		return nil, err
	}

	variations, err := s.repo.FindVariationsByProductID(id)
	if err != nil {
		return nil, err
	}

	var vr []VariationResponse
	for _, v := range variations {
		vr = append(vr, VariationResponse{
			Price: v.Price,
			Stock: v.Stock,
		})
	}

	res := &ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		CategoryID:  product.CategoryID.String(),
		Variations:  vr,
	}

	return res, nil
}
