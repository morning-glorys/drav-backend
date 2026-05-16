package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/service"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// @Summary Add a product review
// @Description Add a review and rating to a product
// @Tags reviews
// @Accept json
// @Produce json
// @Param body body model.CreateReviewRequest true "Review request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	//TODO: Implementasikan handler untuk membuat review, termasuk validasi input dan memastikan user hanya bisa mereview produk yang pernah dibeli gunakan eror handling di pkg/apperror
}

// GetProductReviews godoc
// @Summary Get reviews by product ID
// @Description Retrieve all reviews for a specific product
// @Tags reviews
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{product_id}/reviews [get]
func (h *ReviewHandler) GetReviewsByProductID(c *gin.Context) {
	// TODO: Implementasikan handler untuk mengambil reviews berdasarkan product id, termasuk nama user yang mereview gunakan eror handling di pkg/apperror
	// Pastikan untuk mengembalikan response yang sesuai dengan format yang diharapkan oleh frontend
}
