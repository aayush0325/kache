package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aayush0325/consistent-hashing/internal/routes"
	"github.com/aayush0325/consistent-hashing/internal/services/coordinator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	go coordinator.Background(ctx)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		routes.SetupKeyRouter(api)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, coordinator.GlobalState)
	})

	router.Run()
}
