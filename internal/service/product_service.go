package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
)

var (
	ErrInvalidProductID   = errors.New("invalid product id")
	ErrInvalidProductData = errors.New("invalid product data")
	ErrInvalidQueryParam  = errors.New("invalid query param")
)

type ProductService interface {
	GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error)
	GetProductByID(ctx context.Context, id int) (*model.Product, error)
	CreateProduct(ctx context.Context, req *model.Product) error
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

// get all products
func (s *productService) GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
	if query.Page <= 0 {
		query.Page = 1
	}

	if query.Limit <= 0 {
		query.Limit = 10
	}

	if query.Limit > 100 {
		return nil, ErrInvalidQueryParam
	}

	if query.MinPrice != nil && *query.MinPrice < 0 {
		return nil, ErrInvalidQueryParam
	}

	if query.MaxPrice != nil && *query.MaxPrice < 0 {
		return nil, ErrInvalidQueryParam
	}

	if query.MinPrice != nil && query.MaxPrice != nil && *query.MinPrice > *query.MaxPrice {
		return nil, ErrInvalidQueryParam
	}

	products, err := s.productRepo.GetAllProducts(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}

	return products, nil
}

// get product by id
func (s *productService) GetProductByID(ctx context.Context, id int) (*model.Product, error) {
	if id <= 0 {
		return nil, ErrInvalidProductID
	}
	product, err := s.productRepo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, repository.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return product, nil
}

// create new product
func (s *productService) CreateProduct(ctx context.Context, req *model.Product) error {
	if req == nil {
		return ErrInvalidProductData
	}

	if req.Name == "" || req.Price <= 0 || req.Stock < 0 {
		return ErrInvalidProductData
	}

	if req.SellerID <= 0 {
		return ErrInvalidProductData
	}

	if err := s.productRepo.CreateProduct(ctx, req); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}
