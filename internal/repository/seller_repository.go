package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/morning-glorys/drav-backend/internal/model"
)

var ErrSellerNotFound = errors.New("seller not found")

type SellerRepository interface {
	GetSellerByUserID(ctx context.Context, userID int) (*model.Seller, error)
	CreateSeller(ctx context.Context, seller *model.Seller) error
}

type sellerRepository struct {
	db *sql.DB
}

func NewSellerRepository(db *sql.DB) SellerRepository {
	return &sellerRepository{db: db}
}

// get seller by user id
func (r *sellerRepository) GetSellerByUserID(ctx context.Context, userID int) (*model.Seller, error) {
	query := `SELECT id, user_id, store_name, is_verified, rating, created_at FROM sellers WHERE user_id = $1`
	var s model.Seller
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&s.ID, &s.UserID, &s.StoreName, &s.IsVerified, &s.Rating, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSellerNotFound
		}
		return nil, err
	}

	return &s, nil
}

// create seller
func (r *sellerRepository) CreateSeller(ctx context.Context, seller *model.Seller) error {
	query := `
		INSERT INTO sellers (user_id, store_name)
		VALUES ($1, $2)
		RETURNING id, is_verified, rating, created_at
	`
	err := r.db.QueryRowContext(ctx, query, seller.UserID, seller.StoreName).Scan(&seller.ID, &seller.IsVerified, &seller.Rating, &seller.CreatedAt)
	return err
}
