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
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/internal/service"
)

type mockProductService struct {
	getAllFn  func(ctx context.Context, query model.ProductListQuery) ([]model.Product, error)
	getByIDFn func(ctx context.Context, id int) (*model.Product, error)
	createFn  func(ctx context.Context, req *model.Product) error
}

func (m *mockProductService) GetAllProducts(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx, query)
	}
	return nil, nil
}

func (m *mockProductService) GetProductByID(ctx context.Context, id int) (*model.Product, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockProductService) CreateProduct(ctx context.Context, req *model.Product) error {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil
}

func TestGetProductByID_NotFound_ReturnsGenericNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewProductHandler(&mockProductService{
		getByIDFn: func(ctx context.Context, id int) (*model.Product, error) {
			return nil, repository.ErrProductNotFound
		},
	})

	r := gin.New()
	r.GET("/api/products/:id", h.GetProductByID)

	req := httptest.NewRequest(http.MethodGet, "/api/products/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
	if got := w.Body.String(); got != "{\"error\":\"product not found\"}" {
		t.Fatalf("unexpected body: %s", got)
	}
}

func TestGetProductByID_InvalidParam_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewProductHandler(&mockProductService{})
	r := gin.New()
	r.GET("/api/products/:id", h.GetProductByID)

	req := httptest.NewRequest(http.MethodGet, "/api/products/not-a-number", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateProduct_InvalidData_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewProductHandler(&mockProductService{
		createFn: func(ctx context.Context, req *model.Product) error {
			return service.ErrInvalidProductData
		},
	})
	r := gin.New()
	r.POST("/api/products", h.CreateProduct)

	body := map[string]any{"name": "", "price": 100, "stock": 1}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	if got := w.Body.String(); got != "{\"error\":\"invalid product data\"}" {
		t.Fatalf("unexpected body: %s", got)
	}
}

func TestCreateProduct_InternalError_DoesNotLeakDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewProductHandler(&mockProductService{
		createFn: func(ctx context.Context, req *model.Product) error {
			return errors.New("pq: detailed db failure")
		},
	})
	r := gin.New()
	r.POST("/api/products", h.CreateProduct)

	body := map[string]any{"name": "Phone", "price": 100, "stock": 1}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewReader(bodyBytes))
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

func TestGetAllProducts_InvalidQueryParam_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewProductHandler(&mockProductService{})
	r := gin.New()
	r.GET("/api/products", h.GetAllProducts)

	req := httptest.NewRequest(http.MethodGet, "/api/products?page=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetAllProducts_ServiceInvalidQuery_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewProductHandler(&mockProductService{
		getAllFn: func(ctx context.Context, query model.ProductListQuery) ([]model.Product, error) {
			return nil, service.ErrInvalidQueryParam
		},
	})
	r := gin.New()
	r.GET("/api/products", h.GetAllProducts)

	req := httptest.NewRequest(http.MethodGet, "/api/products?limit=101", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
