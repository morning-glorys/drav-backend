package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
)

var (
	ErrSellerInvalidInput = errors.New("invalid seller input")
	ErrSellerConflict     = errors.New("seller already exists")
	ErrSellerNotFound     = errors.New("seller not found")
)

type SellerService interface {
	RegisterStore(ctx context.Context, userID int, req *model.Seller) error
	GetStoreProfile(ctx context.Context, userID int) (*model.Seller, error)
}

type sellerService struct {
	sellerRepo repository.SellerRepository
}

func NewSellerService(sellerRepo repository.SellerRepository) SellerService {
	return &sellerService{sellerRepo: sellerRepo}
}

// register store
func (s *sellerService) RegisterStore(ctx context.Context, userID int, req *model.Seller) error {
	req.StoreName = strings.TrimSpace(req.StoreName)
	nameLength := utf8.RuneCountInString(req.StoreName)

	if userID <= 0 {
		return ErrSellerInvalidInput
	}
	req.StoreName = strings.TrimSpace(req.StoreName)
	nameLength = utf8.RuneCountInString(req.StoreName)

	if nameLength < 3 || nameLength > 255 {
		return ErrSellerInvalidInput
	}

	req.UserID = userID
	err := s.sellerRepo.CreateSeller(ctx, req)
	if err != nil {
		if errors.Is(err, repository.ErrSellerAlreadyExists) {
			return ErrSellerConflict
		}
		return err
	}
	return nil
}

// get store profile
func (s *sellerService) GetStoreProfile(ctx context.Context, userID int) (*model.Seller, error) {
	seller, err := s.sellerRepo.GetSellerByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrSellerNotFound) {
			return nil, ErrSellerNotFound
		}
		return nil, err
	}
	return seller, nil
}
