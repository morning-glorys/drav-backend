package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/service"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// AddToCart godoc
// @Summary Add item to cart
// @Description Add a product to user's shopping cart
// @Tags carts
// @Accept json
// @Produce json
// @Param body body model.AddToCartRequest true "Add to cart request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /carts [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req model.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	err := h.cartService.AddToCart(c.Request.Context(), userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrCartInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be at least 1"})
		case errors.Is(err, apperror.ErrCartProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		case errors.Is(err, apperror.ErrCartInsufficientStock):
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient product stock"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "item added to cart successfully",
	})
}

// GetMyCart godoc
// @Summary Get user's cart
// @Description Retrieve all items in the authenticated user's cart
// @Tags carts
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /carts [get]
func (h *CartHandler) GetMyCart(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	items, err := h.cartService.GetMyCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve cart"})
		return
	}

	if items == nil {
		items = []model.CartItem{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cart retrieved successfully",
		"data":    items,
	})
}
