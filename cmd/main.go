package main

import (
	"context"
	"net/http"

	"github.com/aayush0325/consistent-hashing/internal/routes"
	"github.com/aayush0325/consistent-hashing/internal/services/coordinator"
	"github.com/aayush0325/consistent-hashing/internal/services/metrics"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load .env.local if it exists, but don't fatal if it's missing (env vars might be set directly)
	_ = godotenv.Load(".env.local")

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

	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metrics.Reg, promhttp.HandlerOpts{Registry: metrics.Reg})))

	router.Run()
}
