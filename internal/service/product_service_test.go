package service

import (
	"context"
	"errors"
	"testing"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
)

type mockProductRepo struct {
	getAllFn     func(ctx context.Context, query model.ProductListQuery) ([]model.Product, error)
	getByIDFn    func(ctx context.Context, id int) (*model.Product, error)
	createProdFn func(ctx context.Context, product *model.Product) error
}

func (m *mockProductRepo) GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, query)
	}
	return nil, nil
}

func (m *mockProductRepo) GetProductByID(ctx context.Context, id int) (*model.Product, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockProductRepo) CreateProduct(ctx context.Context, product *model.Product) error {
	if m.createProdFn != nil {
		return m.createProdFn(ctx, product)
	}
	return nil
}

func TestGetProductByID_InvalidID(t *testing.T) {
	svc := NewProductService(&mockProductRepo{})

	_, err := svc.GetProductByID(context.Background(), 0)
	if !errors.Is(err, ErrInvalidProductID) {
		t.Fatalf("expected ErrInvalidProductID, got %v", err)
	}
}

func TestGetProductByID_ProductNotFound(t *testing.T) {
	svc := NewProductService(&mockProductRepo{
		getByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return nil, repository.ErrProductNotFound
		},
	})

	_, err := svc.GetProductByID(context.Background(), 10)
	if !errors.Is(err, repository.ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestGetProductByID_RepoErrorWrapped(t *testing.T) {
	rootErr := errors.New("db timeout details")
	svc := NewProductService(&mockProductRepo{
		getByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return nil, rootErr
		},
	})

	_, err := svc.GetProductByID(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, rootErr) {
		t.Fatalf("expected wrapped root error, got %v", err)
	}
}

func TestCreateProduct_InvalidData(t *testing.T) {
	svc := NewProductService(&mockProductRepo{})

	err := svc.CreateProduct(context.Background(), &model.Product{SellerID: 1, Name: "", Price: 1, Stock: 0})
	if !errors.Is(err, ErrInvalidProductData) {
		t.Fatalf("expected ErrInvalidProductData, got %v", err)
	}
}

func TestCreateProduct_NilRequest(t *testing.T) {
	svc := NewProductService(&mockProductRepo{})

	err := svc.CreateProduct(context.Background(), nil)
	if !errors.Is(err, ErrInvalidProductData) {
		t.Fatalf("expected ErrInvalidProductData, got %v", err)
	}
}

func TestCreateProduct_Success(t *testing.T) {
	called := false
	svc := NewProductService(&mockProductRepo{
		createProdFn: func(ctx context.Context, product *model.Product) error {
			called = true
			product.ID = 99
			return nil
		},
	})

	product := &model.Product{SellerID: 1, Name: "Laptop", Price: 1500, Stock: 3, Category: "electronics"}
	err := svc.CreateProduct(context.Background(), product)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected repository CreateProduct to be called")
	}
	if product.ID != 99 {
		t.Fatalf("expected product ID to be set by repo, got %d", product.ID)
	}
}

func TestGetAllProducts_InvalidQuery(t *testing.T) {
	svc := NewProductService(&mockProductRepo{})

	_, err := svc.GetAllProducts(context.Background(), model.ProductListQuery{Page: 1, Limit: 200})
	if !errors.Is(err, ErrInvalidQueryParam) {
		t.Fatalf("expected ErrInvalidQueryParam, got %v", err)
	}
}

func TestGetAllProducts_SetsDefaultPagination(t *testing.T) {
	received := model.ProductListQuery{}
	svc := NewProductService(&mockProductRepo{
		getAllFn: func(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
			received = query
			return []model.Product{}, nil
		},
	})

	_, err := svc.GetAllProducts(context.Background(), model.ProductListQuery{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Page != 1 {
		t.Fatalf("expected default page 1, got %d", received.Page)
	}

	if received.Limit != 10 {
		t.Fatalf("expected default limit 10, got %d", received.Limit)
	}
}
