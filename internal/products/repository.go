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
func (r *Repository) FindProducts(limit, offset int) ([]Product, error) {
	var products []Product

	err := database.DB.
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, err
}

func (r *Repository) FindProductsByCategory(categoryID uuid.UUID, limit, offset int) ([]Product, error) {
	var products []Product

	err := database.DB.
		Where("category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error

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

func (r *Repository) CountProducts(categoryID *uuid.UUID) (int64, error) {
	var count int64

	query := database.DB.Model(&Product{})

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *Repository) FindProductsAdvanced(
	categoryID *uuid.UUID,
	minPrice, maxPrice float64,
	sort string,
	limit, offset int,
) ([]Product, error) {

	var products []Product

	query := database.DB.Model(&Product{}).
		Joins("JOIN product_variations ON product_variations.product_id = products.id")

	if categoryID != nil {
		query = query.Where("products.category_id = ?", *categoryID)
	}

	if minPrice > 0 {
		query = query.Where("product_variations.price >= ?", minPrice)
	}

	if maxPrice > 0 {
		query = query.Where("product_variations.price <= ?", maxPrice)
	}

	switch sort {
	case "price_asc":
		query = query.Order("product_variations.price ASC")
	case "price_desc":
		query = query.Order("product_variations.price DESC")
	default:
		query = query.Order("products.created_at DESC")
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, err
}

func (r *Repository) CountProductsAdvanced(
	categoryID *uuid.UUID,
	minPrice, maxPrice float64,
) (int64, error) {

	var count int64

	query := database.DB.Model(&Product{}).
		Joins("JOIN product_variations ON product_variations.product_id = products.id")

	if categoryID != nil {
		query = query.Where("products.category_id = ?", *categoryID)
	}

	if minPrice > 0 {
		query = query.Where("product_variations.price >= ?", minPrice)
	}

	if maxPrice > 0 {
		query = query.Where("product_variations.price <= ?", maxPrice)
	}

	err := query.Distinct("products.id").Count(&count).Error

	return count, err
}
