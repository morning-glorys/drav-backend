package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/service"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
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
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}
	var req model.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format ulasan tidak valid"})
		return
	}
	err := h.reviewService.CreateReview(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, apperror.ErrReviewInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rating harus 1-5 dan komentar tidak boleh kosong"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "terjadi kesalahan internal server"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "ulasan berhasil ditambahkan"})

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
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id produk tidak valid"})
		return
	}
	reviews, err := h.reviewService.GetReviewsByProductID(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil ulasan produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "berhasil mengambil ulasan",
		"data":    reviews,
	})
}
