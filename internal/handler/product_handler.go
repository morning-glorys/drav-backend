package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// get all products
// @Summary Get all products
// @Description Get product list with pagination and optional filters
// @Tags products
// @Produce json
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Items per page" minimum(1) maximum(100) default(10)
// @Param search query string false "Search by product name"
// @Param min_price query number false "Minimum product price" minimum(0)
// @Param max_price query number false "Maximum product price" minimum(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query param"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query param"})
		return
	}

	var minPrice *float64
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		parsed, parseErr := strconv.ParseFloat(minPriceStr, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query param"})
			return
		}
		minPrice = &parsed
	}

	var maxPrice *float64
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		parsed, parseErr := strconv.ParseFloat(maxPriceStr, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query param"})
			return
		}
		maxPrice = &parsed
	}

	products, err := h.productService.GetAllProducts(c.Request.Context(), model.ProductListQuery{
		Page:     page,
		Limit:    limit,
		Search:   c.Query("search"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidQueryParam) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query param"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "products retrieved successfully",
		"data":    products})
}

// get product by id
// @Summary Get product by ID
// @Description Get single product detail by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	product, err := h.productService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProductID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}

		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product retrieved successfully",
		"data":    product})
}

// create new product
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param body body model.Product true "Create product request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req model.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.productService.CreateProduct(c.Request.Context(), &req); err != nil {
		if errors.Is(err, service.ErrInvalidProductData) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product data"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "product created successfully", "data": req})
}
