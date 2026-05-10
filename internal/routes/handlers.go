package routes

import (
	"log"
	"net/http"

	"github.com/aayush0325/consistent-hashing/internal/services/hash"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type setRequest struct {
	Value string `json:"value"`
	TTL   uint64 `json:"ttl"`
}

func handleGet(c *gin.Context) {
	key := c.Param("key")
	log.Printf("GET request for key: %s", key)

	val, err := hash.Get([]byte(key), c.Request.Context())
	if err != nil {
		if err == redis.Nil {
			log.Printf("Key not found: %s", key)
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		log.Printf("Error getting key %s: %v", key, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully retrieved key: %s", key)
	c.JSON(http.StatusOK, gin.H{"value": string(val)})
}

func handlePut(c *gin.Context) {
	key := c.Param("key")
	log.Printf("PUT request for key: %s", key)

	var req setRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON for key %s: %v", key, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := hash.Set(key, req.Value, req.TTL, c.Request.Context()); err != nil {
		log.Printf("Error setting key %s: %v", key, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully set key: %s", key)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
