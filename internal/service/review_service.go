package service

import (
	"context"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
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
	if req.Rating < 1 || req.Rating > 5 || req.Comment == "" {
		return apperror.ErrReviewInvalid
	}
	review := &model.Review{
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	err := s.reviewRepo.CreateReviewAndUpdateSeller(ctx, review)
	if err != nil {
		return apperror.ErrReviewFailed
	}
	return nil
}

// Get reviews by product id
func (s *reviewService) GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error) {
	reviews, err := s.reviewRepo.GetReviewsByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if reviews == nil {
		reviews = []model.Review{}
	}
	return reviews, err
}
