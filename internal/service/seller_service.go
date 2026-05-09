package service

import (
	"context"
	"errors"
	"fmt"
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

const (
	storeNameMinLength = 3
	storeNameMaxLength = 255
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
	if userID <= 0 || req == nil {
		return ErrSellerInvalidInput
	}

	req.StoreName = strings.TrimSpace(req.StoreName)
	if req.StoreName == "" {
		return ErrSellerInvalidInput
	}

	storeNameLen := utf8.RuneCountInString(req.StoreName)
	if storeNameLen < storeNameMinLength || storeNameLen > storeNameMaxLength {
		return ErrSellerInvalidInput
	}

	_, err := s.sellerRepo.GetSellerByUserID(ctx, userID)
	if err == nil {
		return ErrSellerConflict
	}

	if !errors.Is(err, repository.ErrSellerNotFound) {
		return fmt.Errorf("failed to check existing seller: %w", err)
	}

	req.UserID = userID
	err = s.sellerRepo.CreateSeller(ctx, req)
	if err != nil {
		if errors.Is(err, repository.ErrSellerAlreadyExists) {
			return ErrSellerConflict
		}
		return fmt.Errorf("failed to create seller: %w", err)
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
		return nil, fmt.Errorf("failed to get seller profile: %w", err)
	}
	return seller, nil
}
