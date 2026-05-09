package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/service"
)

type SellerHandler struct {
	sellerService service.SellerService
}

func NewSellerHandler(sellerService service.SellerService) *SellerHandler {
	return &SellerHandler{sellerService: sellerService}
}

// register store
// @Summary Register seller store
// @Description Register a seller store for authenticated user
// @Tags sellers
// @Accept json
// @Produce json
// @Param body body model.Seller true "Register seller request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /sellers/register [post]
func (h *SellerHandler) RegisterStore(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req model.Seller
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.sellerService.RegisterStore(c.Request.Context(), uid, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSellerInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller data"})
		case errors.Is(err, service.ErrSellerConflict):
			c.JSON(http.StatusConflict, gin.H{"error": "seller already registered"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "seller registered successfully",
		"data":    req,
	})
}

// @Summary Get seller profile
// @Description Get authenticated user's seller profile
// @Tags sellers
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /sellers/me [get]
func (h *SellerHandler) GetStoreProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	store, err := h.sellerService.GetStoreProfile(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, service.ErrSellerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "seller profile retrieved successfully",
		"data":    store,
	})
}
