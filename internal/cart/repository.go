package cart

import (
	"ecommerce/pkg/database"

	"ecommerce/internal/products"
	"github.com/google/uuid"
)

type Repository struct{}

// 🔹 Find existing cart item (IMPORTANT)
func (r *Repository) FindByUserAndVariation(userID, variationID uuid.UUID) (*CartItem, error) {
	var item CartItem

	err := database.DB.
		Where("user_id = ? AND product_variation_id = ?", userID, variationID).
		First(&item).Error

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// 🔹 Create new cart item
func (r *Repository) Create(item *CartItem) error {
	return database.DB.Create(item).Error
}

// 🔹 Update existing cart item (e.g. quantity)
func (r *Repository) Update(item *CartItem) error {
	return database.DB.Save(item).Error
}

// 🔹 Get all cart items for a user
func (r *Repository) FindByUser(userID uuid.UUID) ([]CartItem, error) {
	var items []CartItem

	err := database.DB.
		Where("user_id = ?", userID).
		Find(&items).Error

	return items, err
}

func (r *Repository) FindVariationByID(id uuid.UUID) (*products.ProductVariation, error) {
	var v products.ProductVariation
	err := database.DB.Where("id = ?", id).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *Repository) FindProductByID(id uuid.UUID) (*products.Product, error) {
	var p products.Product
	err := database.DB.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) Delete(item *CartItem) error {
	return database.DB.Delete(item).Error
}
