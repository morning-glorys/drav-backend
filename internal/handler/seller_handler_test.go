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
	"github.com/morning-glorys/drav-backend/internal/service"
)

type mockSellerService struct {
	registerFn func(ctx context.Context, userID int, req *model.Seller) error
	getStoreFn func(ctx context.Context, userID int) (*model.Seller, error)
}

func (m *mockSellerService) RegisterStore(ctx context.Context, userID int, req *model.Seller) error {
	if m.registerFn != nil {
		return m.registerFn(ctx, userID, req)
	}
	return nil
}

func (m *mockSellerService) GetStoreProfile(ctx context.Context, userID int) (*model.Seller, error) {
	if m.getStoreFn != nil {
		return m.getStoreFn(ctx, userID)
	}
	return nil, nil
}

func TestRegisterStore_Conflict_Returns409(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSellerHandler(&mockSellerService{
		registerFn: func(ctx context.Context, userID int, req *model.Seller) error {
			return service.ErrSellerConflict
		},
	})

	r := gin.New()
	r.POST("/api/sellers/register", func(c *gin.Context) {
		c.Set("userID", 1)
		h.RegisterStore(c)
	})

	body := map[string]any{"store_name": "Store One"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/sellers/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestRegisterStore_InternalError_NoLeak(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSellerHandler(&mockSellerService{
		registerFn: func(ctx context.Context, userID int, req *model.Seller) error {
			return errors.New("db detail leak")
		},
	})

	r := gin.New()
	r.POST("/api/sellers/register", func(c *gin.Context) {
		c.Set("userID", 1)
		h.RegisterStore(c)
	})

	body := map[string]any{"store_name": "Store One"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/sellers/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
	if got := w.Body.String(); got != "{\"error\":\"internal server error\"}" {
		t.Fatalf("unexpected body: %s", got)
	}
}

func TestGetStoreProfile_NotFound_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSellerHandler(&mockSellerService{
		getStoreFn: func(ctx context.Context, userID int) (*model.Seller, error) {
			return nil, service.ErrSellerNotFound
		},
	})

	r := gin.New()
	r.GET("/api/sellers/me", func(c *gin.Context) {
		c.Set("userID", 2)
		h.GetStoreProfile(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/sellers/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}
