package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
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
	if userID <= 0 || req == nil || strings.TrimSpace(req.Address) == "" {
		return nil, apperror.ErrInvalidInput
	}

	cartID, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrCartNotFound) {
			return nil, apperror.ErrOrderCartEmpty
		}
		return nil, err
	}

	cartItems, err := s.cartRepo.GetCartItems(ctx, cartID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, apperror.ErrOrderCartEmpty
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalPrice := 0
	var orderItems []model.OrderItem

	for _, item := range cartItems {
		product, err := s.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if item.Quantity <= 0 {
			return nil, apperror.ErrInvalidInput
		}
		totalPrice += product.Price * item.Quantity
		orderItems = append(orderItems, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	newOrder := &model.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Address:    strings.TrimSpace(req.Address),
	}
	orderID, err := s.orderRepo.CreateOrder(ctx, tx, newOrder)
	if err != nil {
		return nil, err
	}
	newOrder.ID = orderID
	for _, oi := range orderItems {
		oi.OrderID = orderID
		err = s.orderRepo.CreateOrderItem(ctx, tx, &oi)
		if err != nil {
			return nil, err
		}

		// query update stock table products
		updateStockQuery := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`
		res, err := tx.ExecContext(ctx, updateStockQuery, oi.Quantity, oi.ProductID)
		if err != nil {
			return nil, err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if affected != 1 {
			return nil, apperror.ErrOrderInsufficientStock
		}
	}
	// clear cart query
	clearCartQuery := `DELETE FROM cart_items WHERE cart_id = $1`
	if _, err := tx.ExecContext(ctx, clearCartQuery, cartID); err != nil {
		return nil, err
	}
	// commit and save transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return newOrder, nil
}
