package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API running "})
	})

	log.Println("Server running on :8080")

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
