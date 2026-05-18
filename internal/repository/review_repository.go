package repository

import (
	"context"
	"database/sql"

	"github.com/morning-glorys/drav-backend/internal/model"
)

type ReviewRepository interface {
	CreateReviewAndUpdateSeller(ctx context.Context, review *model.Review) error
	GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error)
	ProductExists(ctx context.Context, productID int) (bool, error)
	HasUserPurchasedProduct(ctx context.Context, userID int, productID int) (bool, error)
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
	query := `SELECT r.id, r.user_id, r.product_id, r.rating, r.comment, r.created_at, COALESCE(u.name, 'Deleted User')
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
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

func (r *reviewRepository) HasUserPurchasedProduct(ctx context.Context, userID int, productID int) (bool, error) {
	query := `
		SELECT 1
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		WHERE o.user_id = $1 AND oi.product_id = $2
		LIMIT 1
	`

	var exists int
	err := r.db.QueryRowContext(ctx, query, userID, productID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *reviewRepository) ProductExists(ctx context.Context, productID int) (bool, error) {
	query := `SELECT 1 FROM products WHERE id = $1 LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
