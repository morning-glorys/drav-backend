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
	return 0, nil
	// TODO: implemented create order with transaction
}

// create order item
func (r *orderRepository) CreateOrderItem(ctx context.Context, tx *sql.Tx, item *model.OrderItem) error {
	// TODO: implemented create order item with transaction
	return nil
}
