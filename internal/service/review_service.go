package service

import (
	"context"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
)

type ReviewService interface {
	CreateReview(ctx context.Context, userID int, req *model.CreateReviewRequest) error
	GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error)
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo}
}

// create review
func (s *reviewService) CreateReview(ctx context.Context, userID int, req *model.CreateReviewRequest) error {
	return nil
	// TODO: Implementasikan logic untuk membuat review, termasuk validasi input dan memastikan user hanya bisa mereview produk yang pernah dibeli gunakan eror handling di pkg/apperror
}

// Get reviews by product id
func (s *reviewService) GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error) {
	return nil, nil
	//TODO: Implementasikan logic untuk mengambil reviews berdasarkan product id, termasuk nama user yang mereview gunakan eror handling di pkg/apperror
}
