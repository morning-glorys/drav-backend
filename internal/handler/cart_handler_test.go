package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
)

type mockCartService struct {
	addToCartFn func(ctx context.Context, userID int, req *model.AddToCartRequest) error
	getMyCartFn func(ctx context.Context, userID int) ([]model.CartItem, error)
}

func (m *mockCartService) AddToCart(ctx context.Context, userID int, req *model.AddToCartRequest) error {
	if m.addToCartFn != nil {
		return m.addToCartFn(ctx, userID, req)
	}
	return nil
}

func (m *mockCartService) GetMyCart(ctx context.Context, userID int) ([]model.CartItem, error) {
	if m.getMyCartFn != nil {
		return m.getMyCartFn(ctx, userID)
	}
	return nil, nil
}

func TestAddToCart_InvalidInput_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCartHandler(&mockCartService{
		addToCartFn: func(ctx context.Context, userID int, req *model.AddToCartRequest) error {
			return apperror.ErrCartInvalidInput
		},
	})

	r := gin.New()
	r.POST("/api/carts", func(c *gin.Context) {
		c.Set("userID", 1)
		h.AddToCart(c)
	})

	body := map[string]any{"product_id": 1, "quantity": 0}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/carts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAddToCart_ProductNotFound_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCartHandler(&mockCartService{
		addToCartFn: func(ctx context.Context, userID int, req *model.AddToCartRequest) error {
			return apperror.ErrCartProductNotFound
		},
	})

	r := gin.New()
	r.POST("/api/carts", func(c *gin.Context) {
		c.Set("userID", 1)
		h.AddToCart(c)
	})

	body := map[string]any{"product_id": 999, "quantity": 1}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/carts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetMyCart_Empty_ReturnsArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCartHandler(&mockCartService{
		getMyCartFn: func(ctx context.Context, userID int) ([]model.CartItem, error) {
			return []model.CartItem{}, nil
		},
	})

	r := gin.New()
	r.GET("/api/carts", func(c *gin.Context) {
		c.Set("userID", 1)
		h.GetMyCart(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/carts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() == "" {
		t.Fatal("expected response body")
	}
}

func TestGetMyCart_InternalError_Returns500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewCartHandler(&mockCartService{
		getMyCartFn: func(ctx context.Context, userID int) ([]model.CartItem, error) {
			return nil, errors.New("db timeout")
		},
	})

	r := gin.New()
	r.GET("/api/carts", func(c *gin.Context) {
		c.Set("userID", 1)
		h.GetMyCart(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/carts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
