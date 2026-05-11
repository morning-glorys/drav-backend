package service

import (
	"context"
	"errors"
	"testing"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
)

type mockCartRepo struct {
	getCartByUserIDFn func(ctx context.Context, userID int) (int, error)
	getOrCreateCartFn func(ctx context.Context, userID int) (int, error)
	getItemQuantityFn func(ctx context.Context, cartID, productID int) (int, error)
	upsertCartItemFn  func(ctx context.Context, cartID, productID, quantity int) error
	getCartItemsFn    func(ctx context.Context, cartID int) ([]model.CartItem, error)
}

func (m *mockCartRepo) GetCartByUserID(ctx context.Context, userID int) (int, error) {
	if m.getCartByUserIDFn != nil {
		return m.getCartByUserIDFn(ctx, userID)
	}
	return 0, nil
}

func (m *mockCartRepo) GetOrCreateCart(ctx context.Context, userID int) (int, error) {
	if m.getOrCreateCartFn != nil {
		return m.getOrCreateCartFn(ctx, userID)
	}
	return 0, nil
}

func (m *mockCartRepo) GetCartItemQuantity(ctx context.Context, cartID, productID int) (int, error) {
	if m.getItemQuantityFn != nil {
		return m.getItemQuantityFn(ctx, cartID, productID)
	}
	return 0, nil
}

func (m *mockCartRepo) UpsertCartItem(ctx context.Context, cartID, productID, quantity int) error {
	if m.upsertCartItemFn != nil {
		return m.upsertCartItemFn(ctx, cartID, productID, quantity)
	}
	return nil
}

func (m *mockCartRepo) GetCartItems(ctx context.Context, cartID int) ([]model.CartItem, error) {
	if m.getCartItemsFn != nil {
		return m.getCartItemsFn(ctx, cartID)
	}
	return nil, nil
}

type mockProductRepoForCart struct {
	getProductByIDFn func(ctx context.Context, id int) (*model.Product, error)
}

func (m *mockProductRepoForCart) GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
	return nil, nil
}

func (m *mockProductRepoForCart) GetProductByID(ctx context.Context, id int) (*model.Product, error) {
	if m.getProductByIDFn != nil {
		return m.getProductByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockProductRepoForCart) CreateProduct(ctx context.Context, product *model.Product) error {
	return nil
}

func (m *mockProductRepoForCart) ProductExists(ctx context.Context, id int) (bool, error) {
	return false, nil
}

func (m *mockProductRepoForCart) ProductOwnedByUser(ctx context.Context, productID int, userID int) (bool, error) {
	return false, nil
}

func (m *mockProductRepoForCart) CreateProductImage(ctx context.Context, productID int, imageURL string) (int, error) {
	return 0, nil
}

func TestAddToCart_InvalidInput(t *testing.T) {
	svc := NewCartService(&mockCartRepo{}, &mockProductRepoForCart{})

	err := svc.AddToCart(context.Background(), 1, &model.AddToCartRequest{ProductID: 1, Quantity: 0})
	if !errors.Is(err, apperror.ErrCartInvalidInput) {
		t.Fatalf("expected ErrCartInvalidInput, got %v", err)
	}
}

func TestAddToCart_ProductNotFound(t *testing.T) {
	svc := NewCartService(&mockCartRepo{}, &mockProductRepoForCart{
		getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return nil, repository.ErrProductNotFound
		},
	})

	err := svc.AddToCart(context.Background(), 1, &model.AddToCartRequest{ProductID: 1, Quantity: 1})
	if !errors.Is(err, apperror.ErrCartProductNotFound) {
		t.Fatalf("expected ErrCartProductNotFound, got %v", err)
	}
}

func TestAddToCart_InsufficientStockWithExistingQty(t *testing.T) {
	svc := NewCartService(&mockCartRepo{
		getOrCreateCartFn: func(ctx context.Context, userID int) (int, error) { return 10, nil },
		getItemQuantityFn: func(ctx context.Context, cartID, productID int) (int, error) { return 5, nil },
	}, &mockProductRepoForCart{
		getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return &model.Product{ID: id, Stock: 8}, nil
		},
	})

	err := svc.AddToCart(context.Background(), 1, &model.AddToCartRequest{ProductID: 1, Quantity: 4})
	if !errors.Is(err, apperror.ErrCartInsufficientStock) {
		t.Fatalf("expected ErrCartInsufficientStock, got %v", err)
	}
}

func TestAddToCart_Success(t *testing.T) {
	called := false
	svc := NewCartService(&mockCartRepo{
		getOrCreateCartFn: func(ctx context.Context, userID int) (int, error) { return 10, nil },
		getItemQuantityFn: func(ctx context.Context, cartID, productID int) (int, error) { return 1, nil },
		upsertCartItemFn: func(ctx context.Context, cartID, productID, quantity int) error {
			called = true
			return nil
		},
	}, &mockProductRepoForCart{
		getProductByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return &model.Product{ID: id, Stock: 10}, nil
		},
	})

	err := svc.AddToCart(context.Background(), 1, &model.AddToCartRequest{ProductID: 1, Quantity: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected UpsertCartItem to be called")
	}
}

func TestGetMyCart_NotFoundReturnsEmpty(t *testing.T) {
	svc := NewCartService(&mockCartRepo{
		getCartByUserIDFn: func(ctx context.Context, userID int) (int, error) {
			return 0, repository.ErrCartNotFound
		},
	}, &mockProductRepoForCart{})

	items, err := svc.GetMyCart(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty items, got %d", len(items))
	}
}
