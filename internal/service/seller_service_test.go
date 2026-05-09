package service

import (
	"context"
	"errors"
	"testing"

	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
)

type mockSellerRepo struct {
	getByUserIDFn func(ctx context.Context, userID int) (*model.Seller, error)
	createFn      func(ctx context.Context, seller *model.Seller) error
}

func (m *mockSellerRepo) GetSellerByUserID(ctx context.Context, userID int) (*model.Seller, error) {
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(ctx, userID)
	}
	return nil, nil
}

func (m *mockSellerRepo) CreateSeller(ctx context.Context, seller *model.Seller) error {
	if m.createFn != nil {
		return m.createFn(ctx, seller)
	}
	return nil
}

func TestRegisterStore_InvalidInput(t *testing.T) {
	svc := NewSellerService(&mockSellerRepo{})

	err := svc.RegisterStore(context.Background(), 0, &model.Seller{StoreName: "Store"})
	if !errors.Is(err, ErrSellerInvalidInput) {
		t.Fatalf("expected ErrSellerInvalidInput, got %v", err)
	}
}

func TestRegisterStore_Conflict(t *testing.T) {
	svc := NewSellerService(&mockSellerRepo{
		getByUserIDFn: func(ctx context.Context, userID int) (*model.Seller, error) {
			return &model.Seller{ID: 1, UserID: userID}, nil
		},
	})

	err := svc.RegisterStore(context.Background(), 1, &model.Seller{StoreName: "Store"})
	if !errors.Is(err, ErrSellerConflict) {
		t.Fatalf("expected ErrSellerConflict, got %v", err)
	}
}

func TestRegisterStore_Success(t *testing.T) {
	called := false
	svc := NewSellerService(&mockSellerRepo{
		getByUserIDFn: func(ctx context.Context, userID int) (*model.Seller, error) {
			return nil, repository.ErrSellerNotFound
		},
		createFn: func(ctx context.Context, seller *model.Seller) error {
			called = true
			seller.ID = 10
			return nil
		},
	})

	req := &model.Seller{StoreName: "Store One"}
	err := svc.RegisterStore(context.Background(), 3, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected CreateSeller to be called")
	}
	if req.UserID != 3 {
		t.Fatalf("expected UserID to be set from token context, got %d", req.UserID)
	}
}

func TestGetStoreProfile_NotFound(t *testing.T) {
	svc := NewSellerService(&mockSellerRepo{
		getByUserIDFn: func(ctx context.Context, userID int) (*model.Seller, error) {
			return nil, repository.ErrSellerNotFound
		},
	})

	_, err := svc.GetStoreProfile(context.Background(), 99)
	if !errors.Is(err, ErrSellerNotFound) {
		t.Fatalf("expected ErrSellerNotFound, got %v", err)
	}
}

func TestRegisterStore_TrimStoreName(t *testing.T) {
	var receivedName string
	svc := NewSellerService(&mockSellerRepo{
		getByUserIDFn: func(ctx context.Context, userID int) (*model.Seller, error) {
			return nil, repository.ErrSellerNotFound
		},
		createFn: func(ctx context.Context, seller *model.Seller) error {
			receivedName = seller.StoreName
			return nil
		},
	})

	err := svc.RegisterStore(context.Background(), 1, &model.Seller{StoreName: "   Toko Maju   "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedName != "Toko Maju" {
		t.Fatalf("expected trimmed store name, got %q", receivedName)
	}
}

func TestRegisterStore_StoreNameTooShort(t *testing.T) {
	svc := NewSellerService(&mockSellerRepo{})

	err := svc.RegisterStore(context.Background(), 1, &model.Seller{StoreName: "ab"})
	if !errors.Is(err, ErrSellerInvalidInput) {
		t.Fatalf("expected ErrSellerInvalidInput, got %v", err)
	}
}

func TestRegisterStore_StoreNameTooLong(t *testing.T) {
	longName := ""
	for i := 0; i < 256; i++ {
		longName += "a"
	}

	svc := NewSellerService(&mockSellerRepo{})
	err := svc.RegisterStore(context.Background(), 1, &model.Seller{StoreName: longName})
	if !errors.Is(err, ErrSellerInvalidInput) {
		t.Fatalf("expected ErrSellerInvalidInput, got %v", err)
	}
}

func TestRegisterStore_CreateReturnsAlreadyExists_MapsConflict(t *testing.T) {
	svc := NewSellerService(&mockSellerRepo{
		getByUserIDFn: func(ctx context.Context, userID int) (*model.Seller, error) {
			return nil, repository.ErrSellerNotFound
		},
		createFn: func(ctx context.Context, seller *model.Seller) error {
			return repository.ErrSellerAlreadyExists
		},
	})

	err := svc.RegisterStore(context.Background(), 10, &model.Seller{StoreName: "Store One"})
	if !errors.Is(err, ErrSellerConflict) {
		t.Fatalf("expected ErrSellerConflict, got %v", err)
	}
}

func TestRegisterStore_UnicodeRuneLengthValidation(t *testing.T) {
	svc := NewSellerService(&mockSellerRepo{})

	err := svc.RegisterStore(context.Background(), 1, &model.Seller{StoreName: "あ"})
	if !errors.Is(err, ErrSellerInvalidInput) {
		t.Fatalf("expected ErrSellerInvalidInput for 1-rune name, got %v", err)
	}
}
