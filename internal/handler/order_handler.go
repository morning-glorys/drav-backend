package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/service"
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
	// TODO: implemented handler checkout
}
