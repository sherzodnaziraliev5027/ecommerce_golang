package reviews

import (
	"ecommerce/internal/orders"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo      *Repository
	orderRepo *orders.Repository
}

func NewService(r *Repository, o *orders.Repository) *Service {
	return &Service{
		repo:      r,
		orderRepo: o,
	}
}

func (s *Service) CreateReview(userID, productID uuid.UUID, rating int, comment string) error {

	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// 🔥 NEW CHECK

	review := &Review{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: productID,
		Rating:    rating,
		Comment:   comment,
	}

	return s.repo.Create(review)
}

func (s *Service) GetProductReviews(productID uuid.UUID) ([]Review, error) {
	return s.repo.FindByProduct(productID.String())
}
func (s *Service) GetProductReviewsWithAverage(productID uuid.UUID) (float64, []Review, error) {

	reviews, err := s.repo.FindByProduct(productID.String())
	if err != nil {
		return 0, nil, err
	}

	avg, err := s.repo.GetAverageRating(productID.String())
	if err != nil {
		return 0, nil, err
	}

	return avg, reviews, nil
}
