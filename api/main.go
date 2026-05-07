package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/morning-glorys/drav-backend/internal/handler"
	"github.com/morning-glorys/drav-backend/internal/repository"
	"github.com/morning-glorys/drav-backend/internal/service"
	"github.com/morning-glorys/drav-backend/pkg/database"
)

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
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "Server and DB are up and running securely!",
		})
	})

	// api grouping
	api := r.Group("/api")
	{
		api.POST("/auth/google", authHandler.GoogleLogin)
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
