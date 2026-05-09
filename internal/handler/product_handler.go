package handler

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/morning-glorys/drav-backend/internal/model"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
	uploadImageFn  func(ctx context.Context, file multipart.File, params uploader.UploadParams) (string, error)
}

type AttachProductImageResponse struct {
	ImageID  int    `json:"image_id"`
	ImageURL string `json:"image_url"`
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		uploadImageFn:  uploadToCloudinary,
	}
}

const maxUploadImageSize = 5 * 1024 * 1024

var allowedImageMimeTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
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

	var minPrice *int
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		parsed, parseErr := strconv.Atoi(minPriceStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query param"})
			return
		}
		minPrice = &parsed
	}

	var maxPrice *int
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		parsed, parseErr := strconv.Atoi(maxPriceStr)
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

// upload product image
// @Summary Upload product image
// @Description Upload a product image to cloud storage
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param image formData file true "Image file (jpg, png, webp), max 5MB"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /products/upload-image [post]
func (h *ProductHandler) UploadProductImage(c *gin.Context) {
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

	productID, err := strconv.Atoi(c.PostForm("product_id"))
	if err != nil || productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	if fileHeader.Size <= 0 || fileHeader.Size > maxUploadImageSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image size"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process image"})
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, readErr := file.Read(buf)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image file"})
		return
	}

	contentType := http.DetectContentType(buf[:n])
	if _, ok := allowedImageMimeTypes[contentType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported image format"})
		return
	}

	if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process image"})
		return
	}

	imageURL, err := h.uploadImageFn(c.Request.Context(), file, uploader.UploadParams{Folder: "DRAV"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
		return
	}

	imageID, err := h.productService.AttachProductImage(c.Request.Context(), productID, userID, imageURL)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}

		if errors.Is(err, service.ErrProductForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		if errors.Is(err, service.ErrInvalidProductData) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "image uploaded successfully",
		"data": AttachProductImageResponse{
			ImageID:  imageID,
			ImageURL: imageURL,
		},
	})
}

// helper function to upload image to Cloudinary
func uploadToCloudinary(ctx context.Context, file multipart.File, params uploader.UploadParams) (string, error) {
	cldURL := os.Getenv("CLOUDINARY_URL")
	if cldURL == "" {
		return "", errors.New("cloudinary url is not configured")
	}

	cld, err := cloudinary.NewFromURL(cldURL)
	if err != nil {
		return "", err
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, params)
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
