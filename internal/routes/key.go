package routes

import "github.com/gin-gonic/gin"

func SetupKeyRouter(rg *gin.RouterGroup) {
	keys := rg.Group("/keys")
	keys.GET("/:key", handleGet)
	keys.PUT("/:key", handlePut)
}
