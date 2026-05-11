package repository

import (
	"context"
	"database/sql"

	"github.com/morning-glorys/drav-backend/internal/model"
)

type CartRepository interface {
	GetOrCreateCart(ctx context.Context, userID int) (int, error)
	UpsertCartItem(ctx context.Context, cartID, productID, quantity int) error
	GetCartItems(ctx context.Context, cartID int) ([]model.CartItem, error)
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

// get or create cart for user
func (r *cartRepository) GetOrCreateCart(ctx context.Context, userID int) (int, error) {
	var cartID int
	query := `
		INSERT INTO carts (user_id) 
		VALUES ($1) 
		ON CONFLICT (user_id) DO UPDATE SET user_id = EXCLUDED.user_id 
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&cartID)
	return cartID, err
}

// add to cart with upsert logic: if product already in cart, update quantity; otherwise insert new item
func (r *cartRepository) UpsertCartItem(ctx context.Context, cartID, productID, quantity int) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id) 
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`
	_, err := r.db.ExecContext(ctx, query, cartID, productID, quantity)
	return err
}

// get cart items by cart id
func (r *cartRepository) GetCartItems(ctx context.Context, cartID int) ([]model.CartItem, error) {
	query := `
		SELECT ci.id, ci.cart_id, ci.product_id, ci.quantity, p.name, p.price 
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id = $1
		ORDER BY ci.id DESC
	`
	rows, err := r.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItem
	for rows.Next() {
		var i model.CartItem
		if err := rows.Scan(&i.ID, &i.CartID, &i.ProductID, &i.Quantity, &i.ProductName, &i.ProductPrice); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
