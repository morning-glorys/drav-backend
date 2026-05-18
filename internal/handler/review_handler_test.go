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

type mockReviewService struct {
	createFn       func(ctx context.Context, userID int, req *model.CreateReviewRequest) error
	getByProductFn func(ctx context.Context, productID int) ([]model.Review, error)
}

func (m *mockReviewService) CreateReview(ctx context.Context, userID int, req *model.CreateReviewRequest) error {
	if m.createFn != nil {
		return m.createFn(ctx, userID, req)
	}
	return nil
}

func (m *mockReviewService) GetReviewsByProductID(ctx context.Context, productID int) ([]model.Review, error) {
	if m.getByProductFn != nil {
		return m.getByProductFn(ctx, productID)
	}
	return nil, nil
}

func TestCreateReview_NotPurchased_ReturnsForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewReviewHandler(&mockReviewService{
		createFn: func(ctx context.Context, userID int, req *model.CreateReviewRequest) error {
			return apperror.ErrReviewNotPurchased
		},
	})

	r := gin.New()
	r.POST("/api/reviews", func(c *gin.Context) {
		c.Set("userID", 1)
		h.CreateReview(c)
	})

	body := map[string]any{"product_id": 1, "rating": 5, "comment": "Great"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/reviews", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetReviewsByProductID_InvalidParam_ReturnsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewReviewHandler(&mockReviewService{})

	r := gin.New()
	r.GET("/api/products/:product_id/reviews", h.GetReviewsByProductID)

	req := httptest.NewRequest(http.MethodGet, "/api/products/abc/reviews", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetReviewsByProductID_ProductNotFound_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewReviewHandler(&mockReviewService{
		getByProductFn: func(ctx context.Context, productID int) ([]model.Review, error) {
			return nil, apperror.ErrReviewProductNotFound
		},
	})

	r := gin.New()
	r.GET("/api/products/:product_id/reviews", h.GetReviewsByProductID)

	req := httptest.NewRequest(http.MethodGet, "/api/products/1/reviews", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCreateReview_InternalError_Returns500(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewReviewHandler(&mockReviewService{
		createFn: func(ctx context.Context, userID int, req *model.CreateReviewRequest) error {
			return errors.New("db failure")
		},
	})

	r := gin.New()
	r.POST("/api/reviews", func(c *gin.Context) {
		c.Set("userID", 1)
		h.CreateReview(c)
	})

	body := map[string]any{"product_id": 1, "rating": 5, "comment": "Great"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/reviews", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
