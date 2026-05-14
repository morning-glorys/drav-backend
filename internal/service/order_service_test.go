package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
)

type mockOrderRepo struct {
	createOrderFn     func(ctx context.Context, tx *sql.Tx, order *model.Order) (int, error)
	createOrderItemFn func(ctx context.Context, tx *sql.Tx, item *model.OrderItem) error
}

func (m *mockOrderRepo) CreateOrder(ctx context.Context, tx *sql.Tx, order *model.Order) (int, error) {
	if m.createOrderFn != nil {
		return m.createOrderFn(ctx, tx, order)
	}
	return 0, nil
}

func (m *mockOrderRepo) CreateOrderItem(ctx context.Context, tx *sql.Tx, item *model.OrderItem) error {
	if m.createOrderItemFn != nil {
		return m.createOrderItemFn(ctx, tx, item)
	}
	return nil
}

func TestCheckout_InvalidInput(t *testing.T) {
	svc := NewOrderService(nil, &mockOrderRepo{}, &mockCartRepo{}, &mockProductRepoForCart{})

	_, err := svc.Checkout(context.Background(), 1, nil)
	if !errors.Is(err, apperror.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCheckout_CartEmpty(t *testing.T) {
	svc := NewOrderService(nil, &mockOrderRepo{}, &mockCartRepo{
		getCartByUserIDFn: func(ctx context.Context, userID int) (int, error) {
			return 0, repository.ErrCartNotFound
		},
	}, &mockProductRepoForCart{})

	_, err := svc.Checkout(context.Background(), 1, &model.CheckoutRequest{Address: "Jakarta Selatan 12345"})
	if !errors.Is(err, apperror.ErrOrderCartEmpty) {
		t.Fatalf("expected ErrOrderCartEmpty, got %v", err)
	}
}
