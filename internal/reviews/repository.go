package reviews

import (
	"ecommerce/pkg/database"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(review *Review) error {
	return database.DB.Create(review).Error
}

func (r *Repository) FindByProduct(productID string) ([]Review, error) {
	var reviews []Review
	err := database.DB.Where("product_id = ?", productID).Find(&reviews).Error
	return reviews, err
}

func (r *Repository) GetAverageRating(productID string) (float64, error) {
	var avg float64

	err := database.DB.
		Model(&Review{}).
		Where("product_id = ?", productID).
		Select("AVG(rating)").
		Scan(&avg).Error

	return avg, err
}
