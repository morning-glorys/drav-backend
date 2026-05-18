package service

import (
	"context"
	"strings"

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
	if userID <= 0 || req == nil || req.ProductID <= 0 || req.Rating < 1 || req.Rating > 5 || strings.TrimSpace(req.Comment) == "" {
		return apperror.ErrReviewInvalid
	}

	exists, err := s.reviewRepo.ProductExists(ctx, req.ProductID)
	if err != nil {
		return apperror.ErrReviewFailed
	}
	if !exists {
		return apperror.ErrReviewProductNotFound
	}

	purchased, err := s.reviewRepo.HasUserPurchasedProduct(ctx, userID, req.ProductID)
	if err != nil {
		return apperror.ErrReviewFailed
	}
	if !purchased {
		return apperror.ErrReviewNotPurchased
	}

	review := &model.Review{
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   strings.TrimSpace(req.Comment),
	}
	err = s.reviewRepo.CreateReviewAndUpdateSeller(ctx, review)
	if err != nil {
		return apperror.ErrReviewFailed
	}
	return nil
}

// Get reviews by product id
func (s *reviewService) GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error) {
	if productID <= 0 {
		return nil, apperror.ErrReviewInvalid
	}

	exists, err := s.reviewRepo.ProductExists(ctx, productID)
	if err != nil {
		return nil, apperror.ErrReviewFailed
	}
	if !exists {
		return nil, apperror.ErrReviewProductNotFound
	}

	reviews, err := s.reviewRepo.GetReviewsByProductID(ctx, productID)
	if err != nil {
		return nil, apperror.ErrReviewFailed
	}
	if reviews == nil {
		reviews = []model.Review{}
	}
	return reviews, nil
}
