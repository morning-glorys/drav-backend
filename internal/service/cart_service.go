package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
)

type CartService interface {
	AddToCart(ctx context.Context, userID int, req *model.AddToCartRequest) error
	GetMyCart(ctx context.Context, userID int) ([]model.CartItem, error)
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

// add to cart
func (s *cartService) AddToCart(ctx context.Context, userID int, req *model.AddToCartRequest) error {
	if req.Quantity <= 0 {
		return apperror.ErrCartInvalidInput
	}

	product, err := s.productRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, repository.ErrProductNotFound) {
			return apperror.ErrCartProductNotFound
		}
		return err
	}

	if product.Stock < req.Quantity {
		return apperror.ErrCartInsufficientStock
	}

	cartID, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}

	return s.cartRepo.UpsertCartItem(ctx, cartID, req.ProductID, req.Quantity)
}

// get my cart items
func (s *cartService) GetMyCart(ctx context.Context, userID int) ([]model.CartItem, error) {
	cartID, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.cartRepo.GetCartItems(ctx, cartID)
}
