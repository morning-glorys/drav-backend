package service

import (
	"context"
	"database/sql"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
)

type OrderService interface {
	Checkout(ctx context.Context, userID int, req *model.CheckoutRequest) (*model.Order, error)
}

type orderService struct {
	db          *sql.DB
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewOrderService(db *sql.DB, or repository.OrderRepository, cr repository.CartRepository, pr repository.ProductRepository) OrderService {
	return &orderService{db: db, orderRepo: or, cartRepo: cr, productRepo: pr}
}

// checkout
func (s *orderService) Checkout(ctx context.Context, userID int, req *model.CheckoutRequest) (*model.Order, error) {
	return nil, nil
	// TODO: implemented logic checkout with transaction
}
