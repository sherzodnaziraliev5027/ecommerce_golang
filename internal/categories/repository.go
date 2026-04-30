package categories

import "ecommerce/pkg/database"

type Repository struct{}

// 🔹 Create category
func (r *Repository) Create(category *Category) error {
	return database.DB.Create(category).Error
}

// 🔹 Get all categories
func (r *Repository) FindAll() ([]Category, error) {
	var categories []Category
	err := database.DB.Find(&categories).Error
	return categories, err
}
