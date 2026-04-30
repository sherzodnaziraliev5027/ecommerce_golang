package products

import (
	"ecommerce/pkg/database"

	"github.com/google/uuid"
)

type Repository struct{}

// 🔹 Create product
func (r *Repository) CreateProduct(product *Product) error {
	return database.DB.Create(product).Error
}

// 🔹 Create variation
func (r *Repository) CreateVariation(variation *ProductVariation) error {
	return database.DB.Create(variation).Error
}

// 🔹 Get all products
func (r *Repository) FindAllProducts() ([]Product, error) {
	var products []Product
	err := database.DB.Find(&products).Error
	return products, err
}

// 🔹 Get variations by product
func (r *Repository) FindVariationsByProduct(productID uuid.UUID) ([]ProductVariation, error) {
	var variations []ProductVariation
	err := database.DB.Where("product_id = ?", productID).Find(&variations).Error
	return variations, err
}

func (r *Repository) FindAllVariations() ([]ProductVariation, error) {
	var variations []ProductVariation
	err := database.DB.Find(&variations).Error
	return variations, err
}

func (r *Repository) FindProductByID(id uuid.UUID) (*Product, error) {
	var product Product
	err := database.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) FindVariationsByProductID(id uuid.UUID) ([]ProductVariation, error) {
	var variations []ProductVariation
	err := database.DB.Where("product_id = ?", id).Find(&variations).Error
	return variations, err
}

func (r *Repository) FindVariationByID(id uuid.UUID) (*ProductVariation, error) {
	var v ProductVariation
	err := database.DB.Where("id = ?", id).First(&v).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}
