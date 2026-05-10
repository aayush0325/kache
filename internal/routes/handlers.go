package routes

import (
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

	val, err := hash.Get([]byte(key), c.Request.Context())
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": string(val)})
}

func handlePut(c *gin.Context) {
	key := c.Param("key")

	var req setRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := hash.Set(key, req.Value, req.TTL, c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
