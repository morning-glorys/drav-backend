package service

import (
	"context"
	"errors"
	"testing"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
)

type mockReviewRepo struct {
	createFn        func(ctx context.Context, review *model.Review) error
	getByProductFn  func(ctx context.Context, productID int) ([]model.Review, error)
	productExistsFn func(ctx context.Context, productID int) (bool, error)
	hasPurchasedFn  func(ctx context.Context, userID int, productID int) (bool, error)
}

func (m *mockReviewRepo) CreateReviewAndUpdateSeller(ctx context.Context, review *model.Review) error {
	if m.createFn != nil {
		return m.createFn(ctx, review)
	}
	return nil
}

func (m *mockReviewRepo) GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error) {
	if m.getByProductFn != nil {
		return m.getByProductFn(ctx, productID)
	}
	return nil, nil
}

func (m *mockReviewRepo) ProductExists(ctx context.Context, productID int) (bool, error) {
	if m.productExistsFn != nil {
		return m.productExistsFn(ctx, productID)
	}
	return false, nil
}

func (m *mockReviewRepo) HasUserPurchasedProduct(ctx context.Context, userID int, productID int) (bool, error) {
	if m.hasPurchasedFn != nil {
		return m.hasPurchasedFn(ctx, userID, productID)
	}
	return false, nil
}

func TestCreateReview_NotPurchased(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{
		productExistsFn: func(ctx context.Context, productID int) (bool, error) { return true, nil },
		hasPurchasedFn:  func(ctx context.Context, userID int, productID int) (bool, error) { return false, nil },
	})

	err := svc.CreateReview(context.Background(), 1, &model.CreateReviewRequest{ProductID: 1, Rating: 5, Comment: "Great"})
	if !errors.Is(err, apperror.ErrReviewNotPurchased) {
		t.Fatalf("expected ErrReviewNotPurchased, got %v", err)
	}
}

func TestCreateReview_ProductNotFound(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{
		productExistsFn: func(ctx context.Context, productID int) (bool, error) { return false, nil },
	})

	err := svc.CreateReview(context.Background(), 1, &model.CreateReviewRequest{ProductID: 100, Rating: 5, Comment: "Great"})
	if !errors.Is(err, apperror.ErrReviewProductNotFound) {
		t.Fatalf("expected ErrReviewProductNotFound, got %v", err)
	}
}

func TestGetReviewsByProductID_ProductNotFound(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{
		productExistsFn: func(ctx context.Context, productID int) (bool, error) { return false, nil },
	})

	_, err := svc.GetReviewsByProductID(context.Background(), 123)
	if !errors.Is(err, apperror.ErrReviewProductNotFound) {
		t.Fatalf("expected ErrReviewProductNotFound, got %v", err)
	}
}
