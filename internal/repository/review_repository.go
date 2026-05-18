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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertQuery := `INSERT INTO reviews (user_id, product_id, rating, comment) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err = tx.QueryRowContext(ctx, insertQuery, review.UserID, review.ProductID, review.Rating, review.Comment).Scan(&review.ID, &review.CreatedAt)
	if err != nil {
		return err
	}

	updateSellerQuery := `UPDATE sellers
		SET rating = (
			SELECT COALESCE(AVG(r.rating), 0)
			FROM reviews r
			JOIN products p ON r.product_id = p.id
			WHERE p.seller_id = sellers.id
		)
		WHERE id = (SELECT seller_id FROM products WHERE id = $1)`
	_, err = tx.ExecContext(ctx, updateSellerQuery, review.ProductID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// get reviews by product id
func (r *reviewRepository) GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error) {
	query := `SELECT r.id, r.user_id, r.product_id, r.rating, r.comment, r.created_at, u.name
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC `
	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []model.Review
	for rows.Next() {
		var rev model.Review
		err := rows.Scan(&rev.ID, &rev.UserID, &rev.ProductID, &rev.Rating, &rev.Comment, &rev.CreatedAt, &rev.UserName)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}
