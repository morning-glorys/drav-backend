package repository

import (
	"context"
	"database/sql"

	"github.com/morning-glorys/drav-backend/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, tx *sql.Tx, order *model.Order) (int, error)
	CreateOrderItem(ctx context.Context, tx *sql.Tx, item *model.OrderItem) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

// create order
func (r *orderRepository) CreateOrder(ctx context.Context, tx *sql.Tx, order *model.Order) (int, error) {
	query := `INSERT INTO orders (user_id, total_price, status) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := tx.QueryRowContext(ctx, query, order.UserID, order.TotalPrice, order.Status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// create order item
func (r *orderRepository) CreateOrderItem(ctx context.Context, tx *sql.Tx, item *model.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)`
	_, err := tx.ExecContext(ctx, query, item.OrderID, item.ProductID, item.Quantity, item.Price)
	if err != nil {
		return err
	}
	return nil
}
