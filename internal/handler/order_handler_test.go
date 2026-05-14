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

type mockOrderService struct {
	checkoutFn func(ctx context.Context, userID int, req *model.CheckoutRequest) (*model.Order, error)
}

func (m *mockOrderService) Checkout(ctx context.Context, userID int, req *model.CheckoutRequest) (*model.Order, error) {
	if m.checkoutFn != nil {
		return m.checkoutFn(ctx, userID, req)
	}
	return nil, nil
}

func TestCheckout_EmptyCart_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(&mockOrderService{
		checkoutFn: func(ctx context.Context, userID int, req *model.CheckoutRequest) (*model.Order, error) {
			return nil, apperror.ErrOrderCartEmpty
		},
	})

	r := gin.New()
	r.POST("/api/orders/checkout", func(c *gin.Context) {
		c.Set("userID", 1)
		h.Checkout(c)
	})

	body := map[string]any{"address": "Jl Sudirman No 1 Jakarta"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/orders/checkout", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCheckout_InsufficientStock_ReturnsConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(&mockOrderService{
		checkoutFn: func(ctx context.Context, userID int, req *model.CheckoutRequest) (*model.Order, error) {
			return nil, apperror.ErrOrderInsufficientStock
		},
	})

	r := gin.New()
	r.POST("/api/orders/checkout", func(c *gin.Context) {
		c.Set("userID", 1)
		h.Checkout(c)
	})

	body := map[string]any{"address": "Jl Sudirman No 1 Jakarta"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/orders/checkout", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestCheckout_InternalError_Returns500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewOrderHandler(&mockOrderService{
		checkoutFn: func(ctx context.Context, userID int, req *model.CheckoutRequest) (*model.Order, error) {
			return nil, errors.New("db issue")
		},
	})

	r := gin.New()
	r.POST("/api/orders/checkout", func(c *gin.Context) {
		c.Set("userID", 1)
		h.Checkout(c)
	})

	body := map[string]any{"address": "Jl Sudirman No 1 Jakarta"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/orders/checkout", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
