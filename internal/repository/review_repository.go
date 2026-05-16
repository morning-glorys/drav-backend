package repository

import (
	"context"
	"database/sql"

	"github.com/morning-glorys/drav-backend/internal/model"
)

type ReviewRepository interface {
	CreateReviewAndUpdateSeller(ctx context.Context, review *model.Review) error
	GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error)
}

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

// create review dan update rating seller
func (r *reviewRepository) CreateReviewAndUpdateSeller(ctx context.Context, review *model.Review) error {
	return nil
	// TODO: Implementasikan logic untuk membuat review dan update rating seller secara atomik menggunakan transaction
}

// get reviews by product id
func (r *reviewRepository) GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error) {
	return nil, nil
	//TODO: Implementasikan logic untuk mengambil reviews berdasarkan product id, termasuk nama user yang mereview
}
