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
		log.Fatal("Error loading .env file:", err)
		if err != nil {
			log.Fatal("ENV not found. Failed to connect to database:", err)
		}

		dbUser := os.Getenv("DB_USER")
		dbPass := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

		db, err := database.ConnectDB(dsn)
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}
		defer db.Close()

		// auth route
		userRepo := repository.NewUserRepository(db)
		authService := service.NewAuthService(userRepo)
		authHandler := handler.NewAuthHandler(authService)

		r := gin.Default()

		// health check
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
				"status":  "Server and DB are up and running securely!",
			})
		})

		// API routes
		r = gin.Default()
		api := r.Group("/api")
		{
			api.POST("/auth/google", authHandler.GoogleLogin)
		}

		// start the server
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		fmt.Printf("Server running securely on port %s\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}
}
