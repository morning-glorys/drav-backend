package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/service"
	"github.com/morning-glorys/drav-backend/pkg/apperror"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Checkout godoc
// @Summary Checkout cart
// @Description Process items in cart and create an order
// @Tags orders
// @Accept json
// @Produce json
// @Param body body model.CheckoutRequest true "Checkout request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /orders/checkout [post]
func (h *OrderHandler) Checkout(c *gin.Context) {
	userIDRaw, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, ok := userIDRaw.(int)
	if !ok || userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req model.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	order, err := h.orderService.Checkout(c.Request.Context(), userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperror.ErrOrderCartEmpty):
			c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		case errors.Is(err, apperror.ErrOrderInsufficientStock):
			c.JSON(http.StatusConflict, gin.H{"error": "insufficient stock"})
		case errors.Is(err, apperror.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "checkout successful, order created",
		"data":    order,
	})
}
