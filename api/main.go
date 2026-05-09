package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/morning-glorys/drav-backend/docs"
	"github.com/morning-glorys/drav-backend/internal/handler"
	"github.com/morning-glorys/drav-backend/internal/middleware"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/internal/service"
	"github.com/morning-glorys/drav-backend/pkg/database"
)

// @title Drav Backend API
// @version 1.0
// @description REST API for Drav e-commerce backend.
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: Bearer {token}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: file .env tidak ditemukan, menggunakan variabel bawaan sistem.")
	}

	// db setup
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := database.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer db.Close()

	// Inject User
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Inject Product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Inject Seller
	sellerRepo := repository.NewSellerRepository(db)
	sellerService := service.NewSellerService(sellerRepo)
	sellerHandler := handler.NewSellerHandler(sellerService)

	r := gin.Default()
	r.Use(middleware.CORS())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- PUBLIC ROUTES ---
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "Server and DB are up and running securely!",
		})
	})

	api := r.Group("/api")
	{
		api.POST("/auth/google", middleware.RateLimitByIP(2, 5), authHandler.GoogleLogin)
		api.GET("/products", productHandler.GetAllProducts)
		api.GET("/products/:id", productHandler.GetProductByID)
	}

	// --- PROTECTED ROUTES ---
	protectedAPI := r.Group("/api")
	protectedAPI.Use(middleware.RequireAuth())
	{
		protectedAPI.GET("/profile", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user dari token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Akses Diberikan! Selamat datang di area rahasia.",
				"user_id": userID,
			})
		})
		protectedAPI.POST("/products", productHandler.CreateProduct)
		protectedAPI.POST("/sellers/register", sellerHandler.RegisterStore)
		protectedAPI.GET("/sellers/me", sellerHandler.GetStoreProfile)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server siap menerima traffic secara aman di port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}
